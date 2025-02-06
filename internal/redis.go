package internal

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisService defines the interface for Redis operations.
type RedisService interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	TTL(ctx context.Context, key string) (time.Duration, error)
}

// Redis represents a struct that holds the Redis client for performing Redis operations.
type Redis struct {
	Client *redis.Client
}

// SetupRedis initializes a Redis client connection and returns it.
func SetupRedis(address string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Redis{Client: client}, nil
}

// Set stores a value in Redis with a specified TTL.
func (r *Redis) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("could not set key: %w", err)
	}
	slog.Info("key set in redis", "key", key)
	return nil
}

// Get retrieves a value from Redis.
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("could not get key: %w", err)
	}
	return value, nil
}

// Delete removes a key from Redis.
func (r *Redis) Delete(ctx context.Context, key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("could not delete key: %w", err)
	}
	return nil
}

// TTL retrieves the TTL (time to live) for a key in Redis.
func (r *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.Client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("could not get TTL for key: %w", err)
	}
	return ttl, nil
}
