package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	URL      string `json:"url"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

var RedisClient *redis.Client

func loadRedisConfig(config *RedisConfig) error {
	config.URL = os.Getenv("REDIS_URL")
	if config.URL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}

	config.Password = getEnvWithDefault("REDIS_PASSWORD", "")

	var err error
	config.DB, err = parseIntWithDefault("REDIS_DB", 0)
	if err != nil {
		return fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	return nil
}

func InitRedis(config RedisConfig) {
	if config.URL == "" {
		panic("REDIS_URL environment variable is not set")
	}

	opt, err := redis.ParseURL(config.URL)
	if err != nil {
		panic(err)
	}
	RedisClient = redis.NewClient(opt)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	fmt.Println("Successfully connected to Redis!")
}

// CloseRedis closes the Redis connection gracefully
func CloseRedis() error {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			return fmt.Errorf("failed to close Redis connection: %w", err)
		}
		RedisClient = nil
		fmt.Println("Redis connection closed successfully")
	}
	return nil
}
