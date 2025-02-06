package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretPathService_FormatSecretPath(t *testing.T) {
	tests := []struct {
		name        string
		hmac        string
		key         string
		expected    string
		expectedErr error
	}{
		{
			name:        "Valid HMAC and key",
			hmac:        "validHmac",
			key:         "validKey",
			expected:    "validHmac:secret:validKey",
			expectedErr: nil,
		},
		{
			name:        "Empty HMAC",
			hmac:        "",
			key:         "validKey",
			expected:    "",
			expectedErr: fmt.Errorf("hmac and key cannot be empty"),
		},
		{
			name:        "Empty key",
			hmac:        "validHmac",
			key:         "",
			expected:    "",
			expectedErr: fmt.Errorf("hmac and key cannot be empty"),
		},
		{
			name:        "Empty HMAC and key",
			hmac:        "",
			key:         "",
			expected:    "",
			expectedErr: fmt.Errorf("hmac and key cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secretPathService := &SecretPathServiceImpl{}

			result, err := secretPathService.FormatSecretPath(tt.hmac, tt.key)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "Expected error message does not match")
			} else {
				assert.NoError(t, err, "Unexpected error occurred")
				assert.Equal(t, tt.expected, result, "Unexpected result from FormatSecretPath")
			}
		})
	}
}
