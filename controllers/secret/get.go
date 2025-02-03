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

// GetSecret handles the process of retrieving a secret associated with a specified key path.
func GetSecret(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	fullPath := ctx.Param("key")
	secretKeyPath := strings.TrimPrefix(fullPath, "/")
	if secretKeyPath == "" {
		utils.LogWarn(context.Background(), "missing secret key path", requestID, nil)
		errors.ErrAPIMissingPath.WithRequestID(ctx).JSON(ctx)
		return
	}

	token := utils.GetHeaderToken(ctx)
	tokenHMAC, err := utils.AuthTokenHMAC(ctx)
	if err != nil {
		utils.LogError(context.Background(), "failed to get token hmac", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	secretPath, err := utils.FormatSecretPath(tokenHMAC, secretKeyPath)
	if err != nil {
		utils.LogError(context.Background(), "failed to generate secret key", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	redisClient := utils.GetRedisClient()
	pipe := redisClient.Pipeline()
	getCmd := pipe.Get(context.Background(), secretPath)
	ttlCmd := pipe.TTL(context.Background(), secretPath)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		utils.LogError(context.Background(), "failed to execute redis pipeline", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	encryptedValue, err := getCmd.Result()
	if err != nil {
		utils.LogWarn(context.Background(), "secret not found", requestID, err)
		errors.ErrNotFound.WithRequestID(ctx).JSON(ctx)
		return
	}

	ttl, err := ttlCmd.Result()
	if err != nil || ttl <= 0 {
		utils.LogError(context.Background(), "failed to get secret TTL", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	decryptedValue, err := utils.Decrypt(encryptedValue, token)
	if err != nil {
		utils.LogError(context.Background(), "failed to decrypt secret", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	response := models.GetSecretResponse{
		Value: decryptedValue,
		TTL:   int(ttl.Seconds()),
	}

	ctx.JSON(http.StatusOK, response)
}
