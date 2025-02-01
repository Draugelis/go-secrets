package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
)

const DefaultByteLength = 32

func RandomToken(byteLength ...int) string {
	length := DefaultByteLength
	if len(byteLength) > 0 && byteLength[0] > 0 {
		length = byteLength[0]
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		slog.Error("failed to generate a random token", slog.String("error", err.Error()))
	}
	return hex.EncodeToString(bytes)
}
