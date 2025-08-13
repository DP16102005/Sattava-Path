package database

import (
	"context"
	"fmt"

	"digital-wellbeing-backend/internal/config"

	"github.com/go-redis/redis/v8"
)

func InitializeRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	// Test connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Warning: Redis connection failed: %v\n", err)
		return nil
	}

	fmt.Println("Redis connected successfully")
	return rdb
}
