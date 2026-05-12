package account

import (
	"context"
	"net/netip"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

// handles all database interactions and returns models

type repository struct {
	p *pgxpool.Pool
}

func newRepository(p *pgxpool.Pool) *repository {
	return &repository{p: p}
}

func (r *repository) findByID(ctx context.Context, id string) (*account, error) {
	var acc account
	const q = `SELECT * FROM account WHERE id = $1`
	if err := pgxscan.Get(ctx, r.p, &acc, q, id); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) incrementFailedLoginAttemptsByID(ctx context.Context, id string) (*account, error) {
	var acc account
	const q = `UPDATE account SET failed_login_attempts = failed_login_attempts + 1 WHERE id = $1 RETURNING *`
	if err := pgxscan.Get(ctx, r.p, &acc, q, id); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) resetFailedLoginAttemptsByID(ctx context.Context, id string) (*account, error) {
	var acc account
	const q = `UPDATE account SET failed_login_attempts = 0 WHERE id = $1 RETURNING *`
	if err := pgxscan.Get(ctx, r.p, &acc, q, id); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) updateLockedUntil(ctx context.Context, id string, lockedUntil time.Time) (*account, error) {
	var acc account
	const q = `UPDATE account SET locked_until = $1 WHERE id = $2 RETURNING *`
	if err := pgxscan.Get(ctx, r.p, &acc, q, lockedUntil, id); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) findByUsername(ctx context.Context, username string) (*account, error) {
	var acc account
	const q = `SELECT * FROM account WHERE username = $1`
	if err := pgxscan.Get(ctx, r.p, &acc, q, username); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) create(ctx context.Context, username, passwordHash string) (*account, error) {
	var acc account
	const q = `INSERT INTO account (username, password_hash) VALUES ($1, $2) RETURNING *`
	if err := pgxscan.Get(ctx, r.p, &acc, q, username, passwordHash); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) deleteByID(ctx context.Context, id string) (*account, error) {
	var acc account
	const q = `DELETE FROM account WHERE id = $1 RETURNING *`
	if err := pgxscan.Get(ctx, r.p, &acc, q, id); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) deactivateByID(ctx context.Context, id string) (*account, error) {
	var acc account
	const q = `UPDATE account SET active = false WHERE id = $1 RETURNING *`
	if err := pgxscan.Get(ctx, r.p, &acc, q, id); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) activateByID(ctx context.Context, id string) (*account, error) {
	var acc account
	const q = `UPDATE account SET active = true WHERE id = $1 RETURNING *`
	if err := pgxscan.Get(ctx, r.p, &acc, q, id); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) createRefreshToken(ctx context.Context, accountID string, tokenHash []byte, expiresAt time.Time, userAgent string, ipAddress netip.Addr) (*refreshToken, error) {
	var rt refreshToken
	const q = `INSERT INTO refresh_token (account_id, token_hash, expires_at, user_agent, ip_address) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	if err := pgxscan.Get(ctx, r.p, &rt, q, accountID, tokenHash, expiresAt, userAgent, ipAddress); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *repository) findRefreshTokenByHash(ctx context.Context, hash []byte) (*refreshToken, error) {
	var rt refreshToken
	const q = `SELECT * FROM refresh_token WHERE token_hash = $1`
	if err := pgxscan.Get(ctx, r.p, &rt, q, hash); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *repository) revokeRefreshTokenByHash(ctx context.Context, hash []byte) (*refreshToken, error) {
	var rt refreshToken
	const q = `UPDATE refresh_token SET revoked_at = $2 WHERE token_hash = $1 RETURNING *`
	if err := pgxscan.Get(ctx, r.p, &rt, q, hash, time.Now()); err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *repository) revokeRefreshTokenByID(ctx context.Context, id string) (*refreshToken, error) {
	var rt refreshToken
	const q = `UPDATE refresh_token SET revoked_at = $2 WHERE id = $1 RETURNING *`
	if err := pgxscan.Get(ctx, r.p, &rt, q, id, time.Now()); err != nil {
		return nil, err
	}
	return &rt, nil
}
