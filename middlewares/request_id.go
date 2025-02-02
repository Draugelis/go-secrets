package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates or retrieves a request ID for each request and sets it in the context and response header.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx.Set("request_id", requestID)
		ctx.Writer.Header().Set("X-Request-ID", requestID)
		ctx.Next()
	}
}
