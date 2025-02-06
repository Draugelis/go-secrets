package internal

import (
	"context"
	"io"
	"log/slog"
)

type LoggerService interface {
	Log(ctx context.Context, level slog.Level, message string, requestID string, err error)
	LogError(ctx context.Context, message string, requestID string, err error)
	LogWarn(ctx context.Context, message string, requestID string, err error)
}

type LoggerServiceImpl struct {
	logger *slog.Logger
	writer io.Writer
}

func NewLogger(level slog.Level, writer io.Writer) *LoggerServiceImpl {
	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewJSONHandler(writer, handlerOpts)
	logger := slog.New(handler)

	return &LoggerServiceImpl{
		logger: logger,
		writer: writer,
	}
}

// Log logs a message at the specified level, including optional request ID and error details.
func (l *LoggerServiceImpl) Log(ctx context.Context, level slog.Level, message string, requestID string, err error) {
	var logAttrs []any

	if requestID != "" {
		logAttrs = append(logAttrs, slog.String("request_id", requestID))
	}

	if err != nil {
		logAttrs = append(logAttrs, slog.String("error", err.Error()))
	}

	l.logger.Log(ctx, level, message, logAttrs...)
}

// LogInfo logs a message at the info level, including optional request ID and error details.
func (l *LoggerServiceImpl) LogInfo(ctx context.Context, message string, requestID string) {
	l.Log(ctx, slog.LevelInfo, message, requestID, nil)
}

// LogError logs an error message at the error level, including optional request ID and error details.
func (l *LoggerServiceImpl) LogError(ctx context.Context, message string, requestID string, err error) {
	l.Log(ctx, slog.LevelError, message, requestID, err)
}

// LogWarn logs a warning message at the warning level, including optional request ID and error details.
func (l *LoggerServiceImpl) LogWarn(ctx context.Context, message string, requestID string, err error) {
	l.Log(ctx, slog.LevelWarn, message, requestID, err)
}
