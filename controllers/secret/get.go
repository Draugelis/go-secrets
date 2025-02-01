package controllers

import (
	"context"
	"go-secrets/utils"
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
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "secret key path is required",
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
	// Secret path
	secretPath := utils.FormatSecretPath(tokenHMAC, secretKeyPath)

	// Fetch secret
	redisClient := utils.GetRedisClient()
	encryptedValue, err := redisClient.Get(context.Background(), secretKeyPath).Result()
	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "secret not found",
			},
		)
		return
	}

	// Get TTL
	ttl, err := redisClient.TTL(context.Background(), secretPath).Result()
	if err != nil || ttl <= 0 {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to get secret TTL",
			},
		)
		return
	}

	// Decrypt secret
	decryptedValue, err := utils.Decrypt(encryptedValue, token)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to decrypt secret",
			},
		)
		return
	}

	response := SecretResponse{
		Value: decryptedValue,
		TTL:   int(ttl.Seconds()),
	}

	ctx.JSON(http.StatusOK, response)
}
