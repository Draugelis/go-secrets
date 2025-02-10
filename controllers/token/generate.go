package controllers

import (
	"go-secrets/errors"
	"go-secrets/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const DefaultTTL = 900 // Default TTL of 15 minutes
const MaxTTL = 3600    // Max TTL of 60 minutes

// IssueToken handles the generation and storage of a new token with a specified TTL.
func (tc *TokenControllerImpl) Generate(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	requestID := ctx.GetString("request_id")
	ttlStr := ctx.Query("ttl")
	ttl := DefaultTTL

	if ttlStr != "" {
		parsedTTL, err := strconv.Atoi(ttlStr)
		if err != nil || parsedTTL <= 0 || parsedTTL > MaxTTL {
			tc.Logger.LogError(requestCtx, "invalid ttl value", requestID, err)
			errors.ErrInvalidRequest.WithRequestID(ctx).JSON(ctx)
			return
		}
		ttl = parsedTTL
	}

	token, err := tc.Token.GenerateToken()
	if err != nil {
		tc.Logger.LogError(requestCtx, "failed to generate token", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	tokenHMAC, err := tc.Crypto.GenerateHMAC(token)
	if err != nil {
		tc.Logger.LogError(requestCtx, "failed to generate hmac", requestID, err)
		return
	}

	err = tc.Redis.Set(requestCtx, tokenHMAC, "1", time.Duration(ttl)*time.Second)
	if err != nil {
		tc.Logger.LogError(requestCtx, "failed to store token", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	response := models.IssueTokenResponse{
		Token: token,
		TTL:   ttl,
	}

	ctx.JSON(http.StatusOK, response)
}
