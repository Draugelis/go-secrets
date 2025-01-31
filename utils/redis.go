package utils

import (
	"context"
	"log"
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
		log.Println("Redis connection established")
	})

	return instance, err
}

func GetRedisClient() *redis.Client {
	if instance == nil {
		log.Fatal("Redis client is not initialized. Call SetupRedis first.")
	}
	return instance.Client
}

func (r *Redis) Close() error {
	if err := r.Client.Close(); err != nil {
		log.Printf("Warning: Failed to close Redis connection: %v", err)
		return err
	}
	log.Println("Redis connection closed")
	return nil
}
