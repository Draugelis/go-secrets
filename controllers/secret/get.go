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

func GetSecret(ctx *gin.Context) {
	// Parse the secret key path
	fullPath := ctx.Param("key")
	secretKeyPath := strings.TrimPrefix(fullPath, "/")
	if secretKeyPath == "" {
		slog.Warn("missing secret key path")
		errors.ErrAPIMissingPath.WithRequestID(ctx).JSON(ctx)
		return
	}

	// Generate token HMAC
	token := utils.GetHeaderToken(ctx)
	tokenHMAC, err := utils.AuthTokenHMAC(ctx)
	if err != nil {
		slog.Error("failed to get token hmac", slog.String("error", err.Error()))
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
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
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	encryptedValue, err := getCmd.Result()
	if err != nil {
		slog.Warn("secret not found", slog.String("error", err.Error()))
		errors.ErrNotFound.WithRequestID(ctx).JSON(ctx)
		return
	}

	// Get TTL
	ttl, err := ttlCmd.Result()
	if err != nil || ttl <= 0 {
		slog.Error("failed to get secret TTL", slog.String("error", err.Error()))
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	// Decrypt secret
	decryptedValue, err := utils.Decrypt(encryptedValue, token)
	if err != nil {
		slog.Error("failed to decrypt secret", slog.String("error", err.Error()))
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	response := models.GetSecretResponse{
		Value: decryptedValue,
		TTL:   int(ttl.Seconds()),
	}

	ctx.JSON(http.StatusOK, response)
}
