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

type Redis struct {
	Client *redis.Client
}

var once sync.Once
var instance *Redis

func SetupRedis(address string) (*Redis, error) {
	var err error
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr: address,
		})

		// Test connection to Redis
		if err = client.Ping(context.Background()).Err(); err != nil {
			return
		}

		// Assign the singleton instance
		instance = &Redis{Client: client}
		slog.Info("redis connection established")
	})

	return instance, err
}

func GetRedisClient() *redis.Client {
	if instance == nil {
		LogError(context.Background(), "redis client is not initialized, call SetupRedis first", "", nil)
	}
	return instance.Client
}

func (r *Redis) Close() error {
	if err := r.Client.Close(); err != nil {
		LogWarn(context.Background(), "failed to close redis connection", "", err)
		return err
	}
	slog.Info("redis connection closed")
	return nil
}
