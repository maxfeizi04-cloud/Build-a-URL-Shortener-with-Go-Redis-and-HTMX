package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func NewRedisClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost:6379"
	}
	log.Println("Connecting to Redis server on:", host)

	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return rdb
}

func SetKey(ctx context.Context, rdb *redis.Client, key string, value string, ttl time.Duration) error {
	log.Printf("Setting key: %s to %s in Redis", key, value)

	if err := rdb.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	log.Printf("Key %s has been set successfully", key)
	return nil
}

func GetLongURL(ctx context.Context, rdb *redis.Client, shortURL string) (string, error) {
	longURL, err := rdb.Get(ctx, shortURL).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("short URL not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve from Redis: %v", err)
	}
	return longURL, nil
}
