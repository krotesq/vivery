package db

import (
	"context"
	"fmt"

	"github.com/krotesq/vivery/internal/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context, cfg *config.Config) (rdb *redis.Client, err error) {
	rdb = redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
			Password: cfg.RedisPassword,
			DB:       0,
			Protocol: cfg.RedisProtocolVersion,
		},
	)
	err = rdb.Ping(ctx).Err()
	return
}
