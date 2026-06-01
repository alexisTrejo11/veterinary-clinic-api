package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"clinic-vet-api/internal/shared/rediskeys"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	URL       string `json:"url"`
	KeyPrefix string `json:"key_prefix"`
}

var RedisClient *redis.Client

func loadRedisConfig(config *RedisConfig) error {
	config.URL = os.Getenv("REDIS_URL")
	if config.URL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	config.KeyPrefix = getEnvWithDefault("REDIS_KEY_PREFIX", "vet-clinic-api")
	return nil
}

func InitRedis(config RedisConfig) {
	if config.URL == "" {
		panic("REDIS_URL environment variable is not set")
	}

	rediskeys.Init(config.KeyPrefix)

	opt, err := redis.ParseURL(config.URL)
	if err != nil {
		panic(err)
	}
	RedisClient = redis.NewClient(opt)
	if rediskeys.Prefix() != "" {
		RedisClient.AddHook(keyPrefixHook{})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	fmt.Printf("Successfully connected to Redis (key prefix: %q)!\n", rediskeys.Prefix())
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
