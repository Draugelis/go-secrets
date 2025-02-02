package controllers

import (
	"context"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "secret key path is required"})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete secret"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
