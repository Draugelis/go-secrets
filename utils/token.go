package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
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
		log.Fatalf("Failed to generate random secret: %v", err)
	}
	return hex.EncodeToString(bytes)
}
