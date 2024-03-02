package redis

import (
	"context"
	"fmt"
	"snappfood/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, cfg config.Config) (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, err
}
