package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/krotesq/vivery/internal/config"
)

func ConnectRedis(ctx context.Context, cfg *config.Config) (rdb *redis.Client, err error) {
	rdb = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			Protocol: 2,
		},
	)
	err = rdb.Ping(ctx).Err()
	return
}
