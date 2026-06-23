package mediamtx

import "github.com/jackc/pgx/v5/pgxpool"

type repository struct {
	pool *pgxpool.Pool
}

func newRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
