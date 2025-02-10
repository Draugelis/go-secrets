package controllers

import (
	"go-secrets/errors"
	"go-secrets/helpers"
	"go-secrets/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// StoreSecret handles the process of storing a secret with a specified key path.
func (sc *SecretsControllerImpl) Set(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	requestID := ctx.GetString("request_id")
	fullPath := ctx.Param("key")

	secretKeyPath := strings.TrimPrefix(fullPath, "/")
	if secretKeyPath == "" {
		sc.Logger.LogWarn(requestCtx, "missing secret key path", requestID, nil)
		errors.ErrAPIMissingPath.WithRequestID(ctx).JSON(ctx)
		return
	}

	token, err := sc.Token.GetHeaderToken(ctx)
	if err != nil {
		sc.Logger.LogError(requestCtx, "failed to get token from header", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	tokenHMAC, err := sc.Token.AuthTokenHMAC(ctx)
	if err != nil {
		sc.Logger.LogError(requestCtx, "failed to get token hmac", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	var req models.StoreSecretRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		sc.Logger.LogWarn(requestCtx, "invalid request format", requestID, err)
		errors.ErrInvalidRequest.WithRequestID(ctx).JSON(ctx)
		return
	}

	ttl, err := sc.Redis.TTL(requestCtx, tokenHMAC)
	if err != nil || ttl <= 0 {
		sc.Logger.LogWarn(requestCtx, "invalid or expired token", requestID, err)
		errors.ErrUnauthorized.WithRequestID(ctx).JSON(ctx)
		return
	}

	encryptedValue, err := sc.Crypto.Encrypt(req.Value, token)
	if err != nil {
		sc.Logger.LogError(requestCtx, "encryption failed", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	secretPath, err := helpers.FormatSecretPath(tokenHMAC, secretKeyPath)
	if err != nil {
		sc.Logger.LogError(requestCtx, "failed to generate secret key", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	err = sc.Redis.Set(requestCtx, secretPath, encryptedValue, ttl)
	if err != nil {
		sc.Logger.LogError(requestCtx, "failed to store secret", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	response := models.StoreSecretResponse{
		Key: secretKeyPath,
		TTL: int(ttl.Seconds()),
	}

	ctx.JSON(http.StatusOK, response)
}
