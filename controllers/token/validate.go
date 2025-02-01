package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateToken(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"valid": true,
		},
	)
}
