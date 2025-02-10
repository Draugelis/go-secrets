package models

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response structure with status, message, and an optional request ID.
// @Description API error response format
// @Example { "status": 400, "message": "invalid request", "request_id": "abc-123" }
type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
	RequestID  string `json:"request_id,omitempty"`
}

// NewErrorResponse is a constructor function for ErrorResponse
func NewErrorResponse(status int, message string) ErrorResponse {
	return ErrorResponse{
		StatusCode: status,
		Message:    message,
	}
}

// JSON sends the ErrorResponse as a JSON response with the appropriate status code.
func (e ErrorResponse) JSON(ctx *gin.Context) {
	ctx.JSON(e.StatusCode, e)
}

// WithRequestID attaches the request ID to the ErrorResponse if available and returns the updated response.
func (e ErrorResponse) WithRequestID(ctx *gin.Context) ErrorResponse {
	requestID, _ := ctx.Get("request_id")
	if id, ok := requestID.(string); ok {
		e.RequestID = id
	}
	return e
}
