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

// EncrypterDecrypter defines the methods for encrypting and decrypting data.
type EncrypterDecrypter interface {
	// Encrypt encrypts the provided plaintext string using a token for key derivation.
	// It returns the encrypted data as a base64 encoded string, or an error if the encryption fails.
	Encrypt(plaintext string, token string) (string, error)

	// Decrypt decrypts the provided encrypted base64 string using a token for key derivation.
	// It returns the decrypted plaintext string, or an error if the decryption fails.
	Decrypt(encrypted string, token string) (string, error)
}

// AesGcmCrypto implements the EncrypterDecrypter interface using AES-GCM encryption.
type AesGcmCrypto struct{}

// deriveKey generates a derived key by applying HMAC with SHA256 using the server token and the provided token.
func (e *AesGcmCrypto) deriveKey(token string) ([]byte, error) {
	serverToken, err := config.GetServerToken()
	if err != nil {
		return nil, err
	}

	h := hmac.New(sha256.New, []byte(serverToken))
	h.Write([]byte(token))
	return h.Sum(nil)[:32], nil
}

// Encrypt encrypts the provided value using AES-GCM with a derived key from the given token.
// The encrypted result is returned as a base64 encoded string.
func (e *AesGcmCrypto) Encrypt(value string, token string) (string, error) {
	key, err := e.deriveKey(token)
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

// Decrypt decrypts the provided base64 encoded encrypted string using AES-GCM with a derived key from the given token.
func (e *AesGcmCrypto) Decrypt(encrypted string, token string) (string, error) {
	key, err := e.deriveKey(token)
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
