package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptoService_Encrypt(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		token   string
		wantErr bool
	}{
		{"valid encryption", "secret data", "token123", false},
		{"empty string", "", "token123", false},
		{"special characters", "data@#$%^&*()", "token123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cryptoSvc := &CryptoServiceImpl{}
			got, err := cryptoSvc.Encrypt(tt.input, tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, got)
			assert.True(t, len(got)%4 == 0) // Base64 strings length must be multiple of 4
		})
	}
}

func TestCryptoService_Decrypt(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		token     string
		wantValue string
		wantErr   bool
	}{
		{"valid decryption", "secret data", "token123", "secret data", false},
		{"empty string", "", "token123", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cryptoSvc := &CryptoServiceImpl{}
			encrypted, err := cryptoSvc.Encrypt(tt.input, tt.token)
			assert.NoError(t, err)

			decrypted, err := cryptoSvc.Decrypt(encrypted, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantValue, decrypted)
		})
	}
}

func TestCryptoService_GenerateAndValidateHMAC(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{"valid HMAC generation", "test data", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cryptoSvc := &CryptoServiceImpl{}
			hmacValue, err := cryptoSvc.GenerateHMAC(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, hmacValue)

			isValid := cryptoSvc.ValidateHMAC(tt.data, hmacValue)
			assert.True(t, isValid)
		})
	}
}

func TestCryptoService_InvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		token   string
		method  func(string, string) (string, error)
		wantErr bool
	}{
		{"invalid base64", "invalid-base64", "token123",
			func(s string, token string) (string, error) {
				return (&CryptoServiceImpl{}).Decrypt(s, token)
			}, true},
		{"too short encrypted string", "SGVsbG8gd29ybGQ=", "token123",
			func(s string, token string) (string, error) {
				return (&CryptoServiceImpl{}).Decrypt(s, token)
			}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.method(tt.input, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestCryptoService_KeyDerivation(t *testing.T) {
	cryptoSvc := &CryptoServiceImpl{}
	key1, err1 := cryptoSvc.deriveKey("token123")
	assert.NoError(t, err1)
	assert.Len(t, key1, 32) // AES-256 key size

	key2, err2 := cryptoSvc.deriveKey("token123")
	assert.NoError(t, err2)
	assert.Equal(t, key1, key2) // Same token should produce same key

	key3, err3 := cryptoSvc.deriveKey("different-token")
	assert.NoError(t, err3)
	assert.NotEqual(t, key1, key3) // Different token should produce different key
}
