package routes

import (
	controllers "go-secrets/controllers/secret"
	"go-secrets/middlewares"

	"github.com/gin-gonic/gin"
)

func SecretRoutes(router *gin.Engine) {
	secretGroup := router.Group("/secret").Use(middlewares.AuthMiddleware())
	{
		secretGroup.POST("/*key", controllers.StoreSecret)
		secretGroup.GET("/*key", controllers.GetSecret)
	}
}
