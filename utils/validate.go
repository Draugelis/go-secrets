package utils

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func IsValidToken(tokenHMAC string) bool {
	redisClient := GetRedisClient()
	val, err := redisClient.Get(context.Background(), tokenHMAC).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Printf("Error checking token: %v", err)
		return false
	}

	return val == "1"
}
