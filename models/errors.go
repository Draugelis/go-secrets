package models

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
	RequestID  string `json:"request_id,omitempty"`
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

func (e ErrorResponse) WithRequestID(ctx *gin.Context) ErrorResponse {
	requestID, _ := ctx.Get("request_id")
	if id, ok := requestID.(string); ok {
		e.RequestID = id
	}
	return e
}
