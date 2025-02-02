package controllers

import (
	"context"
	"go-secrets/errors"
	"go-secrets/models"
	"go-secrets/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const DefaultTTL = 900 // Default TTL of 15 minutes
const MaxTTL = 3600    // Max TTL of 60 minutes

// IssueToken handles the generation and storage of a new token with a specified TTL.
func IssueToken(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	ttlStr := ctx.Query("ttl")
	ttl := DefaultTTL

	if ttlStr != "" {
		parsedTTL, err := strconv.Atoi(ttlStr)
		if err != nil || parsedTTL <= 0 || parsedTTL > MaxTTL {
			utils.LogError(context.Background(), "invalid TTL value", requestID, err)
			errors.ErrInvalidRequest.WithRequestID(ctx).JSON(ctx)
			return
		}
		ttl = parsedTTL
	}

	token := utils.RandomToken()
	tokenHMAC, err := utils.HMAC(token)
	if err != nil {
		utils.LogError(context.Background(), "failed to get token hmac", requestID, err)
		return
	}

	redisClient := utils.GetRedisClient()
	err = redisClient.Set(context.Background(), tokenHMAC, "1", time.Duration(ttl)*time.Second).Err()
	if err != nil {
		utils.LogError(context.Background(), "failed to store token", requestID, err)
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	response := models.IssueTokenResponse{
		Token: token,
		TTL:   ttl,
	}

	ctx.JSON(http.StatusOK, response)
}
