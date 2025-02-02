package controllers

import (
	"context"
	"go-secrets/errors"
	"go-secrets/utils"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Delete secret
func DeleteSecret(ctx *gin.Context) {
	// Parse the secret key path
	fullPath := ctx.Param("key")
	secretKeyPath := strings.TrimPrefix(fullPath, "/")
	if secretKeyPath == "" {
		slog.Warn("missing secret key path")
		errors.ErrAPIMissingPath.WithRequestID(ctx).JSON(ctx)
		return
	}

	tokenHMAC, err := utils.AuthTokenHMAC(ctx)
	if err != nil {
		return
	}

	secretKey := utils.FormatSecretPath(tokenHMAC, secretKeyPath)

	redisClient := utils.GetRedisClient()
	if err := redisClient.Del(context.Background(), secretKey).Err(); err != nil {
		slog.Error("failed to delete secret", slog.String("error", err.Error()))
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}
