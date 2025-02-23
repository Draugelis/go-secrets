package routes

import (
	controllers "go-secrets/controllers/token"
	"go-secrets/internal"
	"go-secrets/middlewares"

	"github.com/gin-gonic/gin"
)

// TokenRoute defines the routes for managing tokens under the `/token` endpoint.
func TokenRoute(router *gin.Engine, logger internal.LoggerService, crypto internal.CryptoService, redis internal.RedisService, token internal.TokenService) {
	// Initialize the TokenController
	controller := &controllers.TokenControllerImpl{
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

	tokenGroup := router.Group("/token")
	{
		tokenGroup.GET("", controller.Generate)
		tokenGroup.GET("/valid", authMiddleware.AuthMiddleware(), controller.Validate)
		tokenGroup.DELETE("", authMiddleware.AuthMiddleware(), controller.Delete)
	}
}
