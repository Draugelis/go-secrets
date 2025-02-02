package controllers

import (
	"context"
	"go-secrets/errors"
	"go-secrets/models"
	"go-secrets/utils"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const DefaultTTL = 900 // Default TTL of 15 minutes
const MaxTTL = 3600    // Max TTL of 60 minutes

func IssueToken(ctx *gin.Context) {
	ttlStr := ctx.Query("ttl")
	ttl := DefaultTTL

	if ttlStr != "" {
		parsedTTL, err := strconv.Atoi(ttlStr)
		if err != nil || parsedTTL <= 0 || parsedTTL > MaxTTL {
			slog.Error("invalid TTL value", slog.String("error", err.Error()))
			errors.ErrInvalidRequest.WithRequestID(ctx).JSON(ctx)
			return
		}
		ttl = parsedTTL
	}

	token := utils.RandomToken()
	tokenHMAC, err := utils.HMAC(token)
	if err != nil {
		slog.Error("failed to get token hmac", slog.String("error", err.Error()))
		return
	}

	redisClient := utils.GetRedisClient()
	err = redisClient.Set(context.Background(), tokenHMAC, "1", time.Duration(ttl)*time.Second).Err()
	if err != nil {
		slog.Error("failed to store token", slog.String("error", err.Error()))
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return
	}

	response := models.IssueTokenResponse{
		Token: token,
		TTL:   ttl,
	}

	ctx.JSON(http.StatusOK, response)
}
