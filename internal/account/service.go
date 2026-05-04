package account

import (
	"context"
	"errors"
	"net/netip"
	"time"

	"github.com/krotesq/strowger/internal/auth"
)

type service struct {
	repository     *repository
	jwtSecret      string
	jwtIssuer      string
	jwtExpMin      int
	refreshExpDays int
	bcryptCost     int
}

func newService(repository *repository, jwtSecret string, jwtIssuer string, jwtExpMin int, refreshExpDays int, bcryptCost int) *service {
	return &service{
		repository:     repository,
		jwtSecret:      jwtSecret,
		jwtIssuer:      jwtIssuer,
		jwtExpMin:      jwtExpMin,
		refreshExpDays: refreshExpDays,
		bcryptCost:     bcryptCost,
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
	passwordHash, err := auth.HashPassword(password, service.bcryptCost)
	if err != nil {
		return nil, err
	}

	return service.repository.create(ctx, username, passwordHash)
}

func (service *service) login(ctx context.Context, username, password, userAgent, ip string) (*account, string, string, error) {
	// load account
	account, err := service.repository.findByUsername(ctx, username)
	if err != nil {
		return nil, "", "", err
	}

	if !account.Active {
		return nil, "", "", errors.New("account is deactivated")
	}

	// check if account has >= 3 failed logins
	if account.FailedLoginAttempts >= 3 {
		return nil, "", "", errors.New("too many failed login attempts")
	}

	// verify password
	if err := auth.ComparePassword(password, account.PasswordHash); err != nil {
		// increment failed login attempts
		account, _ = service.repository.incrementFailedLoginAttemptsByID(ctx, account.ID)
		return nil, "", "", err
	}

	// reset failed login attempts if > 0
	if account.FailedLoginAttempts > 0 {
		account, _ = service.repository.resetFailedLoginAttemptsByID(ctx, account.ID)
	}

	// generate jwt
	accessToken, err := auth.GenerateJWT(account.ID, service.jwtIssuer, service.jwtSecret, service.jwtExpMin)
	if err != nil {
		return nil, "", "", err
	}

	// generate refresh token
	refreshToken, refreshTokenHash, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", err
	}

	// parse ip
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil, "", "", err
	}

	// save refreshToken to db
	_, err = service.repository.createRefreshToken(ctx, account.ID, refreshTokenHash, time.Now().AddDate(0, 0, service.refreshExpDays), userAgent, ipAddr)
	if err != nil {
		return nil, "", "", err
	}

	return account, accessToken, refreshToken, nil
}

func (service *service) refresh(ctx context.Context, token, userAgent, ip string) (string, string, error) {
	// hash token
	tokenHash := auth.HashRefreshToken(token)

	// find token in db
	refreshToken, err := service.repository.findRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return "", "", err
	}

	// check if token is not expired and not revoked
	if refreshToken.RevokedAt != nil {
		return "", "", errors.New("refresh token was revoked")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("refresh token is expired")
	}

	// revoke old refresh token
	if _, err := service.repository.revokeRefreshTokenByHash(ctx, tokenHash); err != nil {
		return "", "", err
	}

	// generate new refresh token
	newToken, newTokenHash, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// parse ip
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return "", "", err
	}

	// save new token in db
	refreshToken, err = service.repository.createRefreshToken(ctx, refreshToken.AccountID, newTokenHash, time.Now().AddDate(0, 0, service.refreshExpDays), userAgent, ipAddr)
	if err != nil {
		return "", "", err
	}

	// generate new jwt
	jwt, err := auth.GenerateJWT(refreshToken.AccountID, service.jwtIssuer, service.jwtSecret, service.jwtExpMin)
	if err != nil {
		return "", "", err
	}

	return newToken, jwt, nil
}

func (service *service) me(ctx context.Context) (*account, error) {
	id, ok := auth.AccountIDFromContext(ctx)
	if !ok {
		return nil, errors.New("unable to load id from context")
	}
	return service.repository.findByID(ctx, id)
}

func (service *service) revokeRefreshToken(ctx context.Context, token string) {
	tokenHash := auth.HashRefreshToken(token)
	service.repository.revokeRefreshTokenByHash(ctx, tokenHash)
}
