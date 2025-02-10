package controllers

import (
	"fmt"
	"go-secrets/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Delete all secrets for a token
// @Description Deletes all stored secrets for the authenticated token
// @Tags token
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /token [delete]
func (tc *TokenControllerImpl) Delete(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	requestID := ctx.GetString("request_id")
	tokenHMAC, err := tc.Token.AuthTokenHMAC(ctx)
	if err != nil {
		tc.Logger.LogError(requestCtx, "failed to get token hmac", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}
	keyPattern := fmt.Sprintf("%s*", tokenHMAC)

	iter, err := tc.Redis.NewScanner(requestCtx, keyPattern)
	if err != nil {
		tc.Logger.LogError(requestCtx, "failed to create scanner", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	var keysToDelete []string
	for iter.Next(requestCtx) {
		keysToDelete = append(keysToDelete, iter.Val())
	}

	if err := iter.Err(); err != nil {
		tc.Logger.LogError(requestCtx, "failed to scan secrets", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	// No keys to delete, exit early
	if len(keysToDelete) == 0 {
		ctx.Status(http.StatusNoContent)
		return
	}

	pipeline, err := tc.Redis.NewPipeline(requestCtx)
	if err != nil {
		tc.Logger.LogError(requestCtx, "failed to create pipeline", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	for _, key := range keysToDelete {
		if err := pipeline.Del(requestCtx, key); err != nil {
			tc.Logger.LogError(requestCtx, "failed to delete key", requestID, err)
		}
	}

	_, err = pipeline.Exec(requestCtx)
	if err != nil {
		tc.Logger.LogError(requestCtx, "failed to execute pipeline", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)

}
