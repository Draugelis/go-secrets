package controllers

import (
	"context"
	"go-secrets/utils"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SecretRequest struct {
	Value string `json:"value" binding:"required"`
}

type StoreResponse struct {
	Key string `json:"key"`
	TTL int    `json:"ttl"`
}

func StoreSecret(ctx *gin.Context) {
	// Parse the secret key path
	fullPath := ctx.Param("key")
	secretKeyPath := strings.TrimPrefix(fullPath, "/")
	if secretKeyPath == "" {
		slog.Warn("missing secret key path")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "secret key path is required"})
		return
	}

	// Generate token HMAC
	token := utils.GetHeaderToken(ctx)
	tokenHMAC, err := utils.AuthTokenHMAC(ctx)
	if err != nil {
		return
	}

	var req SecretRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		slog.Warn("invalid request format", slog.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	redisClient := utils.GetRedisClient()

	// Get token TTL
	ttl, err := redisClient.TTL(context.Background(), tokenHMAC).Result()
	if err != nil || ttl <= 0 {
		slog.Warn("invalid or expired token", slog.String("error", err.Error()))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	// Store secret in Redis using {HMAC}:secret:{key} pattern
	encryptedValue, err := utils.Encrypt(req.Value, token)
	if err != nil {
		slog.Error("encryption failed", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "encryption failed"})
		return
	}

	secretPath := utils.FormatSecretPath(tokenHMAC, secretKeyPath)
	err = redisClient.Set(context.Background(), secretPath, encryptedValue, ttl).Err()
	if err != nil {
		slog.Error("failed to store secret", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store secret"})
		return
	}

	response := StoreResponse{
		Key: secretKeyPath,
		TTL: int(ttl.Seconds()),
	}

	ctx.JSON(http.StatusOK, response)
}
