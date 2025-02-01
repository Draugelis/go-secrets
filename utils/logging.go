package utils

import (
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
