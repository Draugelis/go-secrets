package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
		LogError(context.Background(), "failed to generate a random token", "", err)
	}
	return hex.EncodeToString(bytes)
}
