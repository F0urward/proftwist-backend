package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/redis/go-redis/v9"
)

func NewClient(cfg *config.Config) *redis.Client {
	const op = "redis.NewClient"
	logger := ctxutil.GetLogger(context.Background()).WithField("op", op)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,

		DialTimeout:  cfg.Redis.DialTimeout * time.Second,
		ReadTimeout:  cfg.Redis.ReadTimeout * time.Second,
		WriteTimeout: cfg.Redis.WriteTimeout * time.Second,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  cfg.Redis.PoolTimeout * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.WithError(err).Error("cannot ping redis instance")
	}

	return client
}
