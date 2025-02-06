package internal

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type TokenService interface {
	GenerateToken(byteLength ...int) (string, error)
	GetHeaderToken(ctx *gin.Context) (string, error)
	AuthTokenHMAC(ctx *gin.Context) (string, error)
	ValidateToken(token string, r RedisService) (bool, error)
}

type TokenServiceImpl struct{}

// GenerateToken generates a random token of the specified byte length (defaulting to DefaultByteLength if not provided).
func (t *TokenServiceImpl) GenerateToken(byteLength ...int) (string, error) {
	length := 32 // Default value
	if len(byteLength) > 0 && byteLength[0] > 0 {
		length = byteLength[0]
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ValidateToken checks if the given HMAC token is valid by verifying its presence and value in Redis.
func (t *TokenServiceImpl) ValidateToken(token string, r RedisService) (bool, error) {
	val, err := r.Get(context.Background(), token)
	if err == redis.Nil {
		return false, nil // Token does not exist, but it's not an error
	} else if err != nil {
		return false, err
	}

	return val != "", nil
}

// GetHeaderToken retrieves the token from the Authorization header in the request.
func (t *TokenServiceImpl) GetHeaderToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}

// AuthTokenHMAC computes the HMAC of the token from the request header.
func (t *TokenServiceImpl) AuthTokenHMAC(ctx *gin.Context, c CryptoService) (string, error) {
	token, err := t.GetHeaderToken(ctx)
	if err != nil {
		return "", err
	}

	hmacToken, err := c.GenerateHMAC(token)
	if err != nil {
		return "", err
	}

	return hmacToken, nil
}
