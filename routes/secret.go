package routes

import (
	controllers "go-secrets/controllers/secret"
	"go-secrets/middlewares"
	"go-secrets/utils"

	"github.com/gin-gonic/gin"
)

// SecretRoutes defines the routes for managing secrets under the `/secret` endpoint.
func SecretRoutes(router *gin.Engine) {
	crypto := &utils.AesGcmCrypto{}
	secretGroup := router.Group("/secret").Use(middlewares.AuthMiddleware())
	{
		secretGroup.POST("/*key", func(ctx *gin.Context) {
			controllers.StoreSecret(ctx, crypto)
		})
		secretGroup.GET("/*key", func(ctx *gin.Context) {
			controllers.GetSecret(ctx, crypto)
		})
		secretGroup.DELETE("/*key", controllers.DeleteSecret)
	}
}
