package middlewares

import (
	"go-secrets/errors"
	"go-secrets/internal"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddlewareImpl struct {
	Crypto internal.CryptoService
	Token  internal.TokenService
	Redis  internal.RedisService
}

// AuthMiddleware handles the authorization of incoming requests by validating the Authorization header.
func (a *AuthMiddlewareImpl) AuthMiddleware() gin.HandlerFunc {
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
		tokenHMAC, err := a.Crypto.GenerateHMAC(token)
		if err != nil {
			errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
			ctx.Abort()
			return
		}

		valid, err := a.Token.ValidateToken(tokenHMAC, a.Redis)
		if !valid || err != nil {
			errors.ErrUnauthorized.WithRequestID(ctx).JSON(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
