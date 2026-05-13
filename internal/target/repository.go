package target

import (
	"context"
	"database/sql"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

// handles all database interactions and returns models

type repository struct {
	pool *pgxpool.Pool
}

func newRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (repository *repository) findAll(ctx context.Context, accountID string) ([]target, error) {
	var targets []target
	const query = `SELECT id, name, description, target_type_id, account_id, created_at FROM target WHERE account_id = $1`
	if err := pgxscan.Select(ctx, repository.pool, &targets, query, accountID); err != nil {
		return []target{}, err
	}
	return targets, nil
}

func (repository *repository) findWithRtmpByID(ctx context.Context, id string, accountID string) (*target, *rtmp, error) {
	var t target
	const targetQuery = `SELECT id, name, description, target_type_id, account_id, created_at FROM target WHERE id = $1 AND account_id = $2`
	if err := pgxscan.Get(ctx, repository.pool, &t, targetQuery, id, accountID); err != nil {
		return &target{}, &rtmp{}, err
	}

	var r rtmp
	const rtmpQuery = `SELECT target_id, url, stream_key, created_at FROM target_rtmp WHERE target_id = $1`
	if err := pgxscan.Get(ctx, repository.pool, &r, rtmpQuery, t.ID); err != nil {
		return &target{}, &rtmp{}, err
	}

	return &t, &r, nil
}

func (repository *repository) createWithRtmp(ctx context.Context, name string, description string, url string, stream_key string, accountID string) (*target, *rtmp, error) {
	tx, err := repository.pool.Begin(ctx)
	if err != nil {
		return &target{}, &rtmp{}, err
	}
	defer tx.Rollback(ctx)

	var targetType _type
	const targetTypeQuery = `SELECT id, name FROM target_type WHERE name = 'rtmp'`
	if err := pgxscan.Get(ctx, tx, &targetType, targetTypeQuery); err != nil {
		return &target{}, &rtmp{}, err
	}

	const targetQuery = `
			INSERT INTO target (name, description, target_type_id, account_id)
			VALUES (
				$1,
				$2,
				$3,
				$4
			)
			RETURNING id, name, description, target_type_id, account_id, created_at
		`
	var t target
	if err := pgxscan.Get(ctx, tx, &t, targetQuery, name, description, targetType.ID, accountID); err != nil {
		return &target{}, &rtmp{}, err
	}

	const rtmpQuery = `
			INSERT INTO target_rtmp (target_id, url, stream_key)
			VALUES ($1, $2, $3)
			RETURNING target_id, url, stream_key, created_at
		`
	var r rtmp
	if err := pgxscan.Get(ctx, tx, &r, rtmpQuery, t.ID, url, stream_key); err != nil {
		return &target{}, &rtmp{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return &target{}, &rtmp{}, err
	}

	return &t, &r, err
}

func (repository *repository) deleteByID(ctx context.Context, id string, accountID string) error {
	const query = `
		DELETE FROM target
		WHERE id = $1 AND account_id = $2
	`
	cmdTag, err := repository.pool.Exec(ctx, query, id, accountID)

	// always check error first before accesing cmdTag
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return err
}
