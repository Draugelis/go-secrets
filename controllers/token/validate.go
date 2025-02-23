package controllers

import (
	"go-secrets/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Validate a token
// @Description Validates if a token is still active
// @Tags token
// @Security BearerAuth
// @Success 200 {object} models.TokenValidationResponse "Token is valid"
// @Failure 401 {object} models.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /token/valid [get]
func (tc *TokenControllerImpl) Validate(ctx *gin.Context) {
	response := models.TokenValidationResponse{
		Valid: true,
	}
	ctx.JSON(http.StatusOK, response)
}
