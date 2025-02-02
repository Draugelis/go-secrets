package utils

import (
	"context"
	"log/slog"
	"sync"

	"github.com/redis/go-redis/v9"
)

// RedisClient defines an interface for Redis operations
type RedisClient interface {
	CreateUser(token string) error
	CleanupUser() error
	Close() error
}

// Redis represents a struct that holds the Redis client for performing Redis operations.
type Redis struct {
	Client *redis.Client
}

var once sync.Once
var instance *Redis

// SetupRedis initializes a Redis client connection and ensures only one connection instance is created.
func SetupRedis(address string) (*Redis, error) {
	var err error
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr: address,
		})

		if err = client.Ping(context.Background()).Err(); err != nil {
			return
		}

		instance = &Redis{Client: client}
		slog.Info("redis connection established")
	})

	return instance, err
}

// GetRedisClient returns the Redis client instance, ensuring it has been initialized via SetupRedis.
func GetRedisClient() *redis.Client {
	if instance == nil {
		LogError(context.Background(), "redis client is not initialized, call SetupRedis first", "", nil)
	}
	return instance.Client
}

// Close closes the Redis connection.
func (r *Redis) Close() error {
	if err := r.Client.Close(); err != nil {
		LogWarn(context.Background(), "failed to close redis connection", "", err)
		return err
	}
	slog.Info("redis connection closed")
	return nil
}
