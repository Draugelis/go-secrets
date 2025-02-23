package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs details about incoming requests and their processing time.
func LoggingMiddleware() gin.HandlerFunc {
	// TODO: Log request_id
	return func(ctx *gin.Context) {
		startTime := time.Now()

		slog.Debug("incoming request",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.String("remote_ip", ctx.ClientIP()),
			slog.String("request_id", ctx.GetHeader("X-Request-ID")),
		)

		ctx.Next()

		duration := time.Since(startTime).Milliseconds()
		slog.Info("request processed",
			slog.Int64("duration", duration),
			slog.Int("status", ctx.Writer.Status()),
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.String("remote_ip", ctx.ClientIP()),
			slog.String("request_id", ctx.GetHeader("X-Request-ID")),
		)
	}
}
