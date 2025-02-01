package routes

import (
	"go-secrets/controllers"
	"go-secrets/middleware"

	"github.com/gin-gonic/gin"
)

func TokenRoute(router *gin.Engine) {
	tokenGroup := router.Group("/token")
	{
		tokenGroup.GET("", controllers.IssueToken)
		tokenGroup.GET("/valid", middleware.AuthMiddleware(), controllers.ValidateToken)
	}
}
