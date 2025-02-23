package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"go-secrets/config"
	"go-secrets/helpers"
	"io"
)

const NonceSize = 12

// CryptoService defines the methods for encrypting and decrypting data.
type CryptoService interface {
	Encrypt(plaintext string, token string) (string, error)
	Decrypt(encrypted string, token string) (string, error)
	GenerateHMAC(data string) (string, error)
	ValidateHMAC(data string, hmac string) bool
}

type CryptoServiceImpl struct{}

func NewCryptoService() CryptoService {
	return &CryptoServiceImpl{}
}

// deriveKey generates a derived key by applying HMAC with SHA256 using the server token and the provided token.
func (c *CryptoServiceImpl) deriveKey(token string) ([]byte, error) {
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
func (c *CryptoServiceImpl) Encrypt(value string, token string) (string, error) {
	key, err := c.deriveKey(token)
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
func (c *CryptoServiceImpl) Decrypt(encrypted string, token string) (string, error) {
	key, err := c.deriveKey(token)
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

func (c *CryptoServiceImpl) GenerateHMAC(data string) (string, error) {
	serverToken, err := config.GetServerToken()
	if err != nil {
		return "", err
	}

	hmac, err := helpers.GenerateHMAC(data, serverToken)
	if err != nil {
		return "", err
	}

	return hmac, nil
}

func (c *CryptoServiceImpl) ValidateHMAC(data string, receivedHMAC string) bool {
	expectedHMAC, _ := c.GenerateHMAC(data)
	return hmac.Equal([]byte(expectedHMAC), []byte(receivedHMAC))
}
