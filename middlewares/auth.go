package middlewares

import (
	"go-secrets/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			ctx.Abort()
			return
		}

		token := parts[1]
		tokenHMAC, err := utils.HMAC(token)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get token hmac"})
			ctx.Abort()
			return
		}

		if !utils.IsValidToken(tokenHMAC) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
