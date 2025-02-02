package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
