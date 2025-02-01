package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"go-secrets/config"
)

func HMAC(message string) (string, error) {
	serverToken, err := config.GetServerToken()
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(serverToken))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil)), nil
}
