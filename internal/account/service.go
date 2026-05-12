package account

import (
	"context"
	"errors"
	"log"
	"net/netip"
	"time"

	"github.com/krotesq/vivery/internal/auth"
	"github.com/krotesq/vivery/internal/config"
)

type service struct {
	r *repository
	c *config.Config
}

func newService(r *repository, c *config.Config) *service {
	return &service{
		r: r,
		c: c,
	}
}

func (s *service) findByID(ctx context.Context, id string) (*account, error) {
	return s.r.findByID(ctx, id)
}

func (s *service) deleteByID(ctx context.Context, id string) (*account, error) {
	return s.r.deleteByID(ctx, id)
}

func (s *service) deactivateByID(ctx context.Context, id string) (*account, error) {
	return s.r.deactivateByID(ctx, id)
}

func (s *service) activateByID(ctx context.Context, id string) (*account, error) {
	return s.r.activateByID(ctx, id)
}

func (s *service) create(ctx context.Context, username, password string) (*account, error) {
	if err := auth.ValidatePassword(password, username); err != nil {
		return nil, err
	}

	passwordHash, err := auth.HashPassword(password, s.c.BcryptCost)
	if err != nil {
		return nil, err
	}

	return s.r.create(ctx, username, passwordHash)
}

func (s *service) login(ctx context.Context, username, password, userAgent, ip string) (*account, string, string, error) {
	defaultErr := errors.New("could not login")

	// load account
	acc, err := s.r.findByUsername(ctx, username)
	if err != nil {
		return nil, "", "", defaultErr
	}
	accID := acc.ID

	if !acc.Active {
		return nil, "", "", defaultErr
	}

	if !acc.LockedUntil.Before(time.Now()) {
		return nil, "", "", defaultErr
	}

	// verify password
	if err := auth.ComparePassword(password, acc.PasswordHash); err != nil {

		// increment failed login attempts and override account
		acc, err = s.r.incrementFailedLoginAttemptsByID(ctx, acc.ID)
		if err != nil {
			log.Printf("Failed to increment login attempts for account %s: %v", accID, err)
			return nil, "", "", defaultErr
		}

		// check if failed_login_attempts is now > 2, if yes lock account
		if acc.FailedLoginAttempts > 2 {
			lockedMinutes := 5 * acc.FailedLoginAttempts
			acc, err = s.r.updateLockedUntil(ctx, acc.ID, time.Now().Add(time.Minute*time.Duration(lockedMinutes)))
			if err != nil {
				log.Printf("Failed to update locked_until for account %s: %v", accID, err)
				return nil, "", "", defaultErr
			}
			log.Printf("Account %s locked for %d minutes", acc.Username, lockedMinutes)
		}

		return nil, "", "", defaultErr
	}

	// reset failed login attempts if > 0
	if acc.FailedLoginAttempts > 0 {
		if _, err := s.r.resetFailedLoginAttemptsByID(ctx, acc.ID); err != nil {
			log.Printf("failed to reset login attempts for account %s: %v", acc.ID, err)
			return nil, "", "", defaultErr
		}
	}

	// generate jwt
	at, err := auth.GenerateJWT(
		acc.ID,
		s.c.JSONWebTokenIssuer,
		s.c.JSONWebTokenSecret,
		s.c.JSONWebTokenExpireMinutes,
	)
	if err != nil {
		return nil, "", "", defaultErr
	}

	// generate refresh token
	rt, rth, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", defaultErr
	}

	// parse ip
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, "", "", defaultErr
	}

	// save refreshToken to db
	_, err = s.r.createRefreshToken(
		ctx,
		acc.ID,
		rth,
		time.Now().AddDate(0, 0, s.c.RefreshTokenExpireDays),
		userAgent,
		ipAddr,
	)
	if err != nil {
		return nil, "", "", defaultErr
	}

	return acc, at, rt, nil
}

func (s *service) refresh(ctx context.Context, token, userAgent, ip string) (string, string, error) {
	defaultErr := errors.New("could not refresh")

	// hash sent token to find in db
	refreshTokenHash := auth.HashRefreshToken(token)

	// find token in db
	refreshToken, err := s.r.findRefreshTokenByHash(ctx, refreshTokenHash)
	if err != nil {
		return "", "", err
	}

	// check if token is not expired and not revoked
	if refreshToken.RevokedAt != nil {
		return "", "", defaultErr
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return "", "", defaultErr
	}

	// revoke old refresh token
	if _, err := s.r.revokeRefreshTokenByID(ctx, refreshToken.ID); err != nil {
		return "", "", defaultErr
	}

	// generate new refresh token
	newRefreshToken, newRefreshTokenHash, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", defaultErr
	}

	// parse ip
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return "", "", defaultErr
	}

	// save new token in db
	_, err = s.r.createRefreshToken(
		ctx,
		refreshToken.AccountID,
		newRefreshTokenHash,
		time.Now().AddDate(0, 0, s.c.RefreshTokenExpireDays),
		userAgent,
		ipAddr,
	)
	if err != nil {
		return "", "", defaultErr
	}

	// generate new jwt
	jwt, err := auth.GenerateJWT(
		refreshToken.AccountID,
		s.c.JSONWebTokenIssuer,
		s.c.JSONWebTokenSecret,
		s.c.JSONWebTokenExpireMinutes,
	)
	if err != nil {
		return "", "", defaultErr
	}

	return newRefreshToken, jwt, nil
}

func (s *service) me(ctx context.Context) (*account, error) {
	id, ok := auth.AccountIDFromContext(ctx)
	if !ok {
		return nil, errors.New("unable to load account data")
	}
	return s.r.findByID(ctx, id)
}

func (s *service) revokeRefreshToken(ctx context.Context, token string) {
	tokenHash := auth.HashRefreshToken(token)
	s.r.revokeRefreshTokenByHash(ctx, tokenHash)
}
