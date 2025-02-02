package controllers

import (
	"context"
	"fmt"
	"go-secrets/errors"
	"go-secrets/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Delete token and all secrets that belong to it
func DeleteToken(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	tokenHMAC, err := utils.AuthTokenHMAC(ctx)
	if err != nil {
		utils.LogError(context.Background(), "failed to get token hmac", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
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
		utils.LogError(context.Background(), "failed to scan secrets", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
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
			utils.LogError(context.Background(), "failed to delete secrets", requestID, err)
			errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}
