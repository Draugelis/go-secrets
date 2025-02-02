package controllers

import (
	"context"
	"go-secrets/errors"
	"go-secrets/models"
	"go-secrets/utils"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func StoreSecret(ctx *gin.Context) {
	// Parse the secret key path
	fullPath := ctx.Param("key")
	secretKeyPath := strings.TrimPrefix(fullPath, "/")
	if secretKeyPath == "" {
		slog.Warn("missing secret key path")
		errors.ErrAPIMissingPath.JSON(ctx)
		return
	}

	// Generate token HMAC
	token := utils.GetHeaderToken(ctx)
	tokenHMAC, err := utils.AuthTokenHMAC(ctx)
	if err != nil {
		return
	}

	var req models.StoreSecretRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Warn("invalid request format", slog.String("error", err.Error()))
		errors.ErrInvalidRequest.JSON(ctx)
		return
	}

	redisClient := utils.GetRedisClient()

	// Get token TTL
	ttl, err := redisClient.TTL(context.Background(), tokenHMAC).Result()
	if err != nil || ttl <= 0 {
		slog.Warn("invalid or expired token", slog.String("error", err.Error()))
		errors.ErrUnauthorized.JSON(ctx)
		return
	}

	// Store secret in Redis using {HMAC}:secret:{key} pattern
	encryptedValue, err := utils.Encrypt(req.Value, token)
	if err != nil {
		slog.Error("encryption failed", slog.String("error", err.Error()))
		errors.ErrInternalServer.JSON(ctx)
		return
	}

	secretPath := utils.FormatSecretPath(tokenHMAC, secretKeyPath)
	err = redisClient.Set(context.Background(), secretPath, encryptedValue, ttl).Err()
	if err != nil {
		slog.Error("failed to store secret", slog.String("error", err.Error()))
		errors.ErrInternalServer.JSON(ctx)
		return
	}

	response := models.StoreSecretResponse{
		Key: secretKeyPath,
		TTL: int(ttl.Seconds()),
	}

	ctx.JSON(http.StatusOK, response)
}
