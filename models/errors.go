package models

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func NewErrorResponse(status int, message string) ErrorResponse {
	return ErrorResponse{
		StatusCode: status,
		Message:    message,
	}
}

func (e ErrorResponse) JSON(ctx *gin.Context) {
	ctx.JSON(e.StatusCode, e)
}
