package controllers

import (
	"context"
	"fmt"
	"go-secrets/utils"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Delete token and all secrets that belong to it
func DeleteToken(ctx *gin.Context) {
	tokenHMAC, err := utils.AuthTokenHMAC(ctx)
	if err != nil {
		return
	}
	keyPattern := fmt.Sprintf("%s*", tokenHMAC)

	// Fetch all keys related to token
	redisClient := utils.GetRedisClient()
	iter := redisClient.Scan(context.Background(), 0, keyPattern, 100).Iterator()
	var keysToDelete []string

	for iter.Next(context.Background()) {
		keysToDelete = append(keysToDelete, iter.Val())
	}

	if err := iter.Err(); err != nil {
		slog.Error("failed to scan secrets", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan secrets"})
		return
	}

	// Delete keys
	if len(keysToDelete) > 0 {
		_, err := redisClient.Pipelined(context.Background(), func(p redis.Pipeliner) error {
			for _, key := range keysToDelete {
				p.Del(context.Background(), key)
			}
			return nil
		})

		if err != nil {
			slog.Error("failed to delete secrets", slog.String("error", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete secrets"})
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}
