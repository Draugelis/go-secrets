package routes

import (
	"go-secrets/controllers"

	"github.com/gin-gonic/gin"
)

func TokenRoute(router *gin.Engine) {
	tokenGroup := router.Group("/token")
	{
		tokenGroup.GET("", controllers.IssueToken)
	}
}
