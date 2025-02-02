package middlewares

import (
	"go-secrets/errors"
	"go-secrets/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			errors.ErrUnauthorized.JSON(ctx)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errors.ErrUnauthorized.JSON(ctx)
			ctx.Abort()
			return
		}

		token := parts[1]
		tokenHMAC, err := utils.HMAC(token)
		if err != nil {
			errors.ErrInternalServer.JSON(ctx)
			ctx.Abort()
			return
		}

		if !utils.IsValidToken(tokenHMAC) {
			errors.ErrUnauthorized.JSON(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
