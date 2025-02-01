package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"go-secrets/config"
	"io"
)

const NonceSize = 12

// Derives an AES key from a user and server tokens
func deriveKey(token string) ([]byte, error) {
	serverToken, err := config.GetServerToken()
	if err != nil {
		return nil, err
	}

	h := hmac.New(sha256.New, []byte(serverToken))
	h.Write([]byte(token))
	return h.Sum(nil)[:32], nil
}

// Encrypts a value using AES-GCM with a derived key
func Encrypt(value string, token string) (string, error) {
	key, err := deriveKey(token)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, NonceSize) // GCM standard nonce size
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(value), nil)
	encrypted := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypts a value using AES-GCM with a derived key
func Decrypt(encrypted string, token string) (string, error) {
	key, err := deriveKey(token)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < NonceSize {
		return "", errors.New("invalid encrypted data")
	}

	nonce := data[:NonceSize]
	ciphertext := data[NonceSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
