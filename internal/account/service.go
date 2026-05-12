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
	repository *repository
	cfg        *config.Config
}

func newService(repository *repository, cfg *config.Config) *service {
	return &service{
		repository: repository,
		cfg:        cfg,
	}
}

func (service *service) findByID(ctx context.Context, id string) (*account, error) {
	return service.repository.findByID(ctx, id)
}

func (service *service) deleteByID(ctx context.Context, id string) (*account, error) {
	return service.repository.deleteByID(ctx, id)
}

func (service *service) deactivateByID(ctx context.Context, id string) (*account, error) {
	return service.repository.deactivateByID(ctx, id)
}

func (service *service) activateByID(ctx context.Context, id string) (*account, error) {
	return service.repository.activateByID(ctx, id)
}

func (service *service) create(ctx context.Context, username, password string) (*account, error) {
	if err := auth.ValidatePassword(password, username); err != nil {
		return nil, err
	}

	passwordHash, err := auth.HashPassword(password, service.cfg.BcryptCost)
	if err != nil {
		return nil, err
	}

	return service.repository.create(ctx, username, passwordHash)
}

func (service *service) login(ctx context.Context, username, password, userAgent, ip string) (*account, string, string, error) {
	defaultErr := errors.New("could not login")

	// load account
	account, err := service.repository.findByUsername(ctx, username)
	if err != nil {
		return nil, "", "", defaultErr
	}

	if !account.Active {
		return nil, "", "", defaultErr
	}

	// check if account has >= 3 failed logins
	if account.FailedLoginAttempts >= 3 {
		return nil, "", "", defaultErr
	}

	// verify password
	if err := auth.ComparePassword(password, account.PasswordHash); err != nil {
		// increment failed login attempts
		if _, err := service.repository.incrementFailedLoginAttemptsByID(ctx, account.ID); err != nil {
			log.Printf("failed to increment login attempts for account %s: %v", account.ID, err)
		}
		return nil, "", "", defaultErr
	}

	// reset failed login attempts if > 0
	if account.FailedLoginAttempts > 0 {
		if _, err := service.repository.resetFailedLoginAttemptsByID(ctx, account.ID); err != nil {
			log.Printf("failed to reset login attempts for account %s: %v", account.ID, err)
			return nil, "", "", defaultErr
		}
	}

	// generate jwt
	accessToken, err := auth.GenerateJWT(
		account.ID,
		service.cfg.JSONWebTokenIssuer,
		service.cfg.JSONWebTokenSecret,
		service.cfg.JSONWebTokenExpireMinutes,
	)
	if err != nil {
		return nil, "", "", defaultErr
	}

	// generate refresh token
	refreshToken, refreshTokenHash, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", defaultErr
	}

	// parse ip
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, "", "", defaultErr
	}

	// save refreshToken to db
	_, err = service.repository.createRefreshToken(
		ctx,
		account.ID,
		refreshTokenHash,
		time.Now().AddDate(0, 0, service.cfg.RefreshTokenExpireDays),
		userAgent,
		ipAddr,
	)
	if err != nil {
		return nil, "", "", defaultErr
	}

	return account, accessToken, refreshToken, nil
}

func (service *service) refresh(ctx context.Context, token, userAgent, ip string) (string, string, error) {
	defaultErr := errors.New("could not refresh")

	// hash sent token to find in db
	refreshTokenHash := auth.HashRefreshToken(token)

	// find token in db
	refreshToken, err := service.repository.findRefreshTokenByHash(ctx, refreshTokenHash)
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
	if _, err := service.repository.revokeRefreshTokenByID(ctx, refreshToken.ID); err != nil {
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
	_, err = service.repository.createRefreshToken(
		ctx,
		refreshToken.AccountID,
		newRefreshTokenHash,
		time.Now().AddDate(0, 0, service.cfg.RefreshTokenExpireDays),
		userAgent,
		ipAddr,
	)
	if err != nil {
		return "", "", defaultErr
	}

	// generate new jwt
	jwt, err := auth.GenerateJWT(
		refreshToken.AccountID,
		service.cfg.JSONWebTokenIssuer,
		service.cfg.JSONWebTokenSecret,
		service.cfg.JSONWebTokenExpireMinutes,
	)
	if err != nil {
		return "", "", defaultErr
	}

	return newRefreshToken, jwt, nil
}

func (service *service) me(ctx context.Context) (*account, error) {
	id, ok := auth.AccountIDFromContext(ctx)
	if !ok {
		return nil, errors.New("unable to load account data")
	}
	return service.repository.findByID(ctx, id)
}

func (service *service) revokeRefreshToken(ctx context.Context, token string) {
	tokenHash := auth.HashRefreshToken(token)
	service.repository.revokeRefreshTokenByHash(ctx, tokenHash)
}
