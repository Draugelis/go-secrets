package routes

import (
	controllers "go-secrets/controllers/token"
	"go-secrets/middlewares"

	"github.com/gin-gonic/gin"
)

func TokenRoute(router *gin.Engine) {
	tokenGroup := router.Group("/token")
	{
		tokenGroup.GET("", controllers.IssueToken)
		tokenGroup.GET("/valid", middlewares.AuthMiddleware(), controllers.ValidateToken)
		tokenGroup.DELETE("", middlewares.AuthMiddleware(), controllers.DeleteToken)
	}
}
