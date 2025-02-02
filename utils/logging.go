package utils

import (
	"context"
	"log/slog"
	"os"
)

func InitializeLogger(level slog.Level) {
	handlerOps := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewJSONHandler(os.Stdout, handlerOps)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

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

func LogError(ctx context.Context, message string, requestID string, err error) {
	Log(ctx, slog.LevelError, message, requestID, err)
}

func LogWarn(ctx context.Context, message string, requestID string, err error) {
	Log(ctx, slog.LevelWarn, message, requestID, err)
}
