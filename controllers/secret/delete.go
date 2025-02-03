package controllers

import (
	"context"
	"go-secrets/errors"
	"go-secrets/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DeleteSecret handles the process of deleting a secret associated with a specified key path.
func DeleteSecret(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	fullPath := ctx.Param("key")
	secretKeyPath := strings.TrimPrefix(fullPath, "/")
	if secretKeyPath == "" {
		utils.LogWarn(context.Background(), "missing secret key path", requestID, nil)
		errors.ErrAPIMissingPath.WithRequestID(ctx).JSON(ctx)
		return
	}

	tokenHMAC, err := utils.AuthTokenHMAC(ctx)
	if err != nil {
		utils.LogError(context.Background(), "failed to get token hmac", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	secretPath, err := utils.FormatSecretPath(tokenHMAC, secretKeyPath)
	if err != nil {
		utils.LogError(context.Background(), "failed to generate secret key", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	redisClient := utils.GetRedisClient()
	if err := redisClient.Del(context.Background(), secretPath).Err(); err != nil {
		utils.LogError(context.Background(), "failed to delete secret", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}
