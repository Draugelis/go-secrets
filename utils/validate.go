package utils

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func IsValidToken(tokenHMAC string) bool {
	redisClient := GetRedisClient()
	val, err := redisClient.Get(context.Background(), tokenHMAC).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		slog.Error("error validating token", slog.String("error", err.Error()))
		return false
	}

	return val == "1"
}
