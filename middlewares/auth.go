package middlewares

import (
	"go-secrets/errors"
	"go-secrets/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles the authorization of incoming requests by validating the Authorization header.
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			errors.ErrUnauthorized.WithRequestID(ctx).JSON(ctx)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errors.ErrUnauthorized.WithRequestID(ctx).JSON(ctx)
			ctx.Abort()
			return
		}

		token := parts[1]
		tokenHMAC, err := utils.HMAC(token)
		if err != nil {
			errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
			ctx.Abort()
			return
		}

		if !utils.IsValidToken(tokenHMAC) {
			errors.ErrUnauthorized.WithRequestID(ctx).JSON(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
