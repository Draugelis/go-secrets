package controllers

import (
	"context"
	"go-secrets/utils"
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
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Secret key path is required",
			},
		)
		return
	}

	// Generate token HMAC
	authHeader := ctx.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	token := parts[1]
	tokenHMAC, err := utils.HMAC(token)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to get token hmac",
			},
		)
		ctx.Abort()
		return
	}

	var req SecretRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	redisClient := utils.GetRedisClient()

	// Get token TTL
	ttl, err := redisClient.TTL(context.Background(), tokenHMAC).Result()
	if err != nil || ttl <= 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Store secret in Redis using {HMAC}:secret:{key} pattern
	encryptedValue, err := utils.Encrypt(req.Value, token)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Encryption failed",
			},
		)
		return
	}

	secretPath := utils.FormatSecretPath(tokenHMAC, secretKeyPath)
	err = redisClient.Set(context.Background(), secretPath, encryptedValue, ttl).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store secret"})
		return
	}

	response := StoreResponse{
		Key: secretKeyPath,
		TTL: int(ttl.Seconds()),
	}

	ctx.JSON(http.StatusOK, response)
}
