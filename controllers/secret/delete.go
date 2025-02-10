package controllers

import (
	"go-secrets/errors"
	"go-secrets/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Delete a secret
// @Description Deletes a secret by key path
// @Tags secret
// @Param key path string true "Secret key"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse "Missing key path"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /secret/{key} [delete]
func (sc *SecretsControllerImpl) Delete(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	requestID := ctx.GetString("request_id")
	fullPath := ctx.Param("key")

	secretKeyPath := strings.TrimPrefix(fullPath, "/")
	if secretKeyPath == "" {
		sc.Logger.LogWarn(requestCtx, "missing secret key path", requestID, nil)
		errors.ErrAPIMissingPath.WithRequestID(ctx).JSON(ctx)
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

	if err := sc.Redis.Del(requestCtx, secretPath); err != nil {
		sc.Logger.LogError(requestCtx, "failed to delete secret", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}
