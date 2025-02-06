package internal

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the CryptoService interface with only the necessary methods (GenerateHMAC, Decrypt)
type MockCryptoService struct {
	mock.Mock
}

func (c *MockCryptoService) GenerateHMAC(data string) (string, error) {
	args := c.Called(data)
	return args.String(0), args.Error(1)
}

func (c *MockCryptoService) Decrypt(encrypted, token string) (string, error) {
	args := c.Called(encrypted, token)
	return args.String(0), args.Error(1)
}

func (c *MockCryptoService) Encrypt(plaintext, token string) (string, error) {
	// No-op for testing purposes
	return "", nil
}

func (c *MockCryptoService) ValidateHMAC(data, hmac string) bool {
	// No-op for testing purposes
	return false
}

// Mocking the RedisService interface with only the necessary methods (Get, Set)
type MockRedisService struct {
	mock.Mock
}

func (r *MockRedisService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	args := r.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (r *MockRedisService) Get(ctx context.Context, key string) (string, error) {
	args := r.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (r *MockRedisService) Delete(ctx context.Context, key string) error {
	args := r.Called(ctx, key)
	return args.Error(0)
}

func (r *MockRedisService) TTL(ctx context.Context, key string) (time.Duration, error) {
	// No-op for testing purposes
	return 0, nil
}

// Test function for GenerateToken
func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name        string
		byteLength  int
		expectedLen int
	}{
		{"default length", 0, 64},    // Default length (32 bytes, represented by 64 hex chars)
		{"specified length", 16, 32}, // Specified length (16 bytes, represented by 32 hex chars)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenService := &TokenServiceImpl{}

			// Act
			mockToken, err := tokenService.GenerateToken(tt.byteLength)

			// Assert
			assert.NoError(t, err)
			assert.Len(t, mockToken, tt.expectedLen)
		})
	}
}

// Test function for ValidateToken
func TestValidateToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		redisVal string
		expected bool
	}{
		{"token exists", "valid_token", "some_value", true},
		{"token does not exist", "invalid_token", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(MockRedisService)
			mockRedis.On("Get", context.Background(), tt.token).Return(tt.redisVal, nil)

			// Act
			tokenService := &TokenServiceImpl{}
			valid, err := tokenService.ValidateToken(tt.token, mockRedis)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, valid)
			mockRedis.AssertExpectations(t)
		})
	}
}

// Test function for GetHeaderToken
func TestGetHeaderToken(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expected    string
		expectedErr error
	}{
		{"valid header", "Bearer valid_token", "valid_token", nil},
		{"invalid header", "InvalidHeader", "", errors.New("invalid authorization header format")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &gin.Context{}
			ctx.Request = &http.Request{
				Header: map[string][]string{
					"Authorization": {tt.authHeader},
				},
			}

			// Act
			tokenService := &TokenServiceImpl{}
			token, err := tokenService.GetHeaderToken(ctx)

			// Assert
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, token)
			}
		})
	}
}

// Test function for AuthTokenHMAC
func TestAuthTokenHMAC(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		mockHMAC      string
		mockCryptoErr error
		expectedHMAC  string
		expectedErr   error
	}{
		{
			name:          "valid_token",
			token:         "validToken",
			mockHMAC:      "mockedHMAC",
			mockCryptoErr: nil,
			expectedHMAC:  "mockedHMAC",
			expectedErr:   nil,
		},
		{
			name:          "crypto_error",
			token:         "validToken",
			mockHMAC:      "",
			mockCryptoErr: errors.New("crypto error"),
			expectedHMAC:  "",
			expectedErr:   errors.New("crypto error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock CryptoService and set expectations
			mockCrypto := new(MockCryptoService)
			mockCrypto.On("GenerateHMAC", tt.token).Return(tt.mockHMAC, tt.mockCryptoErr)

			// Create a new gin context for the test
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = &http.Request{
				Header: map[string][]string{
					"Authorization": {"Bearer " + tt.token}, // Set the Authorization header for the test
				},
			}

			// Create the token service
			tokenService := &TokenServiceImpl{}

			// Act
			hmac, err := tokenService.AuthTokenHMAC(ctx, mockCrypto)

			// Assert
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedHMAC, hmac)
			}

			// Assert expectations on mock
			mockCrypto.AssertExpectations(t)
		})
	}
}
