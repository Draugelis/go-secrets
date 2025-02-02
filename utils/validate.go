package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func IsValidToken(tokenHMAC string) bool {
	redisClient := GetRedisClient()
	val, err := redisClient.Get(context.Background(), tokenHMAC).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		LogError(context.Background(), "error validating token", "", err)
		return false
	}

	return val == "1"
}
