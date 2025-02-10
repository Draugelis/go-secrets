package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidateToken handles the validation of a token and responds with its validity status.
func (tc *TokenControllerImpl) Validate(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"valid": true})
}
