package controllers

import (
	"context"
	"go-secrets/utils"
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
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid TTL value"},
			)
			return
		}
		ttl = parsedTTL
	}

	token := utils.RandomToken()
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

	redisClient := utils.GetRedisClient()
	err = redisClient.Set(context.Background(), tokenHMAC, "1", time.Duration(ttl)*time.Second).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "ttl": ttl})
}
