package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(cfg *config.Config) *redis.Client {
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
		log.Fatalf("cannot ping redis instance: %v", err)
	}

	return client
}
