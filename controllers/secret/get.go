package controllers

import (
	"go-secrets/errors"
	"go-secrets/helpers"
	"go-secrets/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Retrieve a secret
// @Description Gets a secret by key path
// @Tags secret
// @Param key path string true "Secret key"
// @Security BearerAuth
// @Success 200 {object} models.GetSecretResponse "Secret retrieved"
// @Failure 400 {object} models.ErrorResponse "Missing key path"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /secret/{key} [get]
func (sc *SecretsControllerImpl) Get(ctx *gin.Context) {
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

	secretPath, err := helpers.FormatSecretPath(tokenHMAC, secretKeyPath)
	if err != nil {
		sc.Logger.LogError(requestCtx, "failed to format secret path", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	encryptedValue, err := sc.Redis.Get(requestCtx, secretPath)
	if err != nil {
		sc.Logger.LogError(requestCtx, "failed to get secret", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	ttl, err := sc.Redis.TTL(requestCtx, secretPath)
	if err != nil || ttl <= 0 {
		sc.Logger.LogError(requestCtx, "failed to get secret TTL", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	decryptedValue, err := sc.Crypto.Decrypt(encryptedValue, token)
	if err != nil {
		sc.Logger.LogError(requestCtx, "failed to decrypt secret", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	response := models.GetSecretResponse{
		Value: decryptedValue,
		TTL:   int(ttl.Seconds()),
	}

	ctx.JSON(http.StatusOK, response)
}
