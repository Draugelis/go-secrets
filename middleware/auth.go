package middleware

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
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Authorization header is missing",
				},
			)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Invalid Authorization header format",
				},
			)
			ctx.Abort()
			return
		}

		token := parts[1]
		tokenHMAC := utils.HMAC(token)

		if !utils.IsValidToken(tokenHMAC) {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Invalid or expired token",
				},
			)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
