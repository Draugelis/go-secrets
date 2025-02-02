package controllers

import (
	"context"
	"go-secrets/errors"
	"go-secrets/models"
	"go-secrets/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func StoreSecret(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	// Parse the secret key path
	fullPath := ctx.Param("key")
	secretKeyPath := strings.TrimPrefix(fullPath, "/")
	if secretKeyPath == "" {
		utils.LogWarn(context.Background(), "missing secret key path", requestID, nil)
		errors.ErrAPIMissingPath.WithRequestID(ctx).JSON(ctx)
		return
	}

	// Generate token HMAC
	token := utils.GetHeaderToken(ctx)
	tokenHMAC, err := utils.AuthTokenHMAC(ctx)
	if err != nil {
		utils.LogError(context.Background(), "failed to get token hmac", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	var req models.StoreSecretRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.LogWarn(context.Background(), "invalid request format", requestID, err)
		errors.ErrInvalidRequest.WithRequestID(ctx).JSON(ctx)
		return
	}

	redisClient := utils.GetRedisClient()

	// Get token TTL
	ttl, err := redisClient.TTL(context.Background(), tokenHMAC).Result()
	if err != nil || ttl <= 0 {
		utils.LogWarn(context.Background(), "invalid or expired token", requestID, err)
		errors.ErrUnauthorized.WithRequestID(ctx).JSON(ctx)
		return
	}

	// Store secret in Redis using {HMAC}:secret:{key} pattern
	encryptedValue, err := utils.Encrypt(req.Value, token)
	if err != nil {
		utils.LogError(context.Background(), "encryption failed", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	secretPath := utils.FormatSecretPath(tokenHMAC, secretKeyPath)
	err = redisClient.Set(context.Background(), secretPath, encryptedValue, ttl).Err()
	if err != nil {
		utils.LogError(context.Background(), "failed to store secret", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	response := models.StoreSecretResponse{
		Key: secretKeyPath,
		TTL: int(ttl.Seconds()),
	}

	ctx.JSON(http.StatusOK, response)
}
