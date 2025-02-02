package routes

import (
	controllers "go-secrets/controllers/secret"
	"go-secrets/middlewares"

	"github.com/gin-gonic/gin"
)

// SecretRoutes defines the routes for managing secrets under the `/secret` endpoint.
func SecretRoutes(router *gin.Engine) {
	secretGroup := router.Group("/secret").Use(middlewares.AuthMiddleware())
	{
		secretGroup.POST("/*key", controllers.StoreSecret)
		secretGroup.GET("/*key", controllers.GetSecret)
		secretGroup.DELETE("/*key", controllers.DeleteSecret)
	}
}
