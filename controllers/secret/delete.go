package controllers

import (
	"context"
	"go-secrets/utils"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "secret key path is required"})
		return
	}

	token := utils.GetHeaderToken(ctx)
	tokenHMAC, err := utils.HMAC(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get token hmac"})
		ctx.Abort()
		return
	}

	secretKey := utils.FormatSecretPath(tokenHMAC, secretKeyPath)

	redisClient := utils.GetRedisClient()
	if err := redisClient.Del(context.Background(), secretKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete secret"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
