package controllers

import (
	"context"
	"go-secrets/utils"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SecretResponse struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

func GetSecret(ctx *gin.Context) {
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
	tokenHMAC, err := utils.HMAC(token)
	if err != nil {
		slog.Error("failed to get token hmac", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get token hmac"})
		return
	}
	// Secret path
	secretPath := utils.FormatSecretPath(tokenHMAC, secretKeyPath)

	// Fetch values from redis
	redisClient := utils.GetRedisClient()

	pipe := redisClient.Pipeline()
	getCmd := pipe.Get(context.Background(), secretPath)
	ttlCmd := pipe.TTL(context.Background(), secretPath)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		slog.Error("failed to execute redis pipeline", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to execute redis pipeline"})
		return
	}

	encryptedValue, err := getCmd.Result()
	if err != nil {
		slog.Warn("secret not found", slog.String("error", err.Error()))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "secret not found"})
		return
	}

	// Get TTL
	ttl, err := ttlCmd.Result()
	if err != nil || ttl <= 0 {
		slog.Error("failed to get secret TTL", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get secret TTL"})
		return
	}

	// Decrypt secret
	decryptedValue, err := utils.Decrypt(encryptedValue, token)
	if err != nil {
		slog.Error("failed to decrypt secret", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decrypt secret"})
		return
	}

	response := SecretResponse{
		Value: decryptedValue,
		TTL:   int(ttl.Seconds()),
	}

	ctx.JSON(http.StatusOK, response)
}
