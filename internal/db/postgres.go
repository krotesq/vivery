package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krotesq/vivery/internal/config"
)

func ConnectPostgres(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.DatabaseUser, cfg.DatabasePassword),
		Host:     fmt.Sprintf("%s:%s", cfg.DatabaseHost, cfg.DatabasePort),
		Path:     cfg.DatabaseName,
		RawQuery: fmt.Sprintf("sslmode=%s", cfg.DatabaseSSLMode),
	}).String()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
