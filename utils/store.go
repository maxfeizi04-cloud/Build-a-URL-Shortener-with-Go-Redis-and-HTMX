package utils

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	fmt.Println("connecting to redis server on:", os.Getenv("REDIS_HOST"))

	// Create a new Redis client
	// 创建一个新的 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return rdb
}

func SetKey(ctx *context.Context, rdb *redis.Client, key string, value string, ttl int) {
	// We set the key value pair in Redis,
	// we use the context defined in main by reference and a TTL of 0 (no expiration)
	// 我们在Redis中设置键值对， 我们使用在main中定义的上下文引用，并且TTL为O（无过期）
	fmt.Println("setting key:", key, "to", value, "in Redis")

	rdb.Set(*ctx, key, value, 0)
	fmt.Println("the key:", key, "has been set to", value, "successfully")
}

func GetLongURL(ctx *context.Context, rdb *redis.Client, shortURL string) (string, error) {
	longURL, err := rdb.Get(*ctx, shortURL).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("short URL not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve from Redis: %v", err)
	}
	return longURL, nil
}
