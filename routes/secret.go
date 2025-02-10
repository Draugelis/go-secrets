package routes

import (
	controllers "go-secrets/controllers/secret"
	"go-secrets/internal"
	"go-secrets/middlewares"

	"github.com/gin-gonic/gin"
)

// SecretRoutes defines the routes for managing secrets under the `/secret` endpoint.
func SecretRoutes(router *gin.Engine, logger internal.LoggerService, crypto internal.CryptoService, redis internal.RedisService, token internal.TokenService) {
	// Initialize the SecretsController
	controller := &controllers.SecretsControllerImpl{
		Logger: logger,
		Crypto: crypto,
		Redis:  redis,
		Token:  token,
	}

	// Initialize AuthMiddlewareImpl
	authMiddleware := &middlewares.AuthMiddlewareImpl{
		Crypto: crypto,
		Token:  token,
		Redis:  redis,
	}

	secretGroup := router.Group("/secret").Use(authMiddleware.AuthMiddleware())
	{
		secretGroup.POST("/*key", controller.Set)
		secretGroup.GET("/*key", controller.Get)
		secretGroup.DELETE("/*key", controller.Delete)
	}
}
