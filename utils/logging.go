package utils

import (
	"context"
	"log/slog"
	"os"
)

// InitializeLogger sets up the default logger with the specified log level and outputs logs in JSON format.
func InitializeLogger(level slog.Level) {
	handlerOps := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewJSONHandler(os.Stdout, handlerOps)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// Log logs a message at the specified level, including optional request ID and error details.
func Log(ctx context.Context, level slog.Level, message string, requestID string, err error) {
	var logAttrs []any

	if requestID != "" {
		logAttrs = append(logAttrs, slog.String("request_id", requestID))
	}

	if err != nil {
		logAttrs = append(logAttrs, slog.String("error", err.Error()))
	}

	slog.Log(ctx, level, message, logAttrs...)
}

// LogError logs an error message at the error level, including optional request ID and error details.
func LogError(ctx context.Context, message string, requestID string, err error) {
	Log(ctx, slog.LevelError, message, requestID, err)
}

// LogWarn logs a warning message at the warning level, including optional request ID and error details.
func LogWarn(ctx context.Context, message string, requestID string, err error) {
	Log(ctx, slog.LevelWarn, message, requestID, err)
}
