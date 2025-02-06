package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvService_Get(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		expected     string
		setEnvValue  string
		defaultValue string
	}{
		{
			name:         "Get existing environment variable",
			key:          "TEST_ENV",
			expected:     "test_value",
			setEnvValue:  "test_value",
			defaultValue: "default_value",
		},
		{
			name:         "Get non-existing environment variable with default value",
			key:          "NON_EXISTING_ENV",
			expected:     "default_value",
			setEnvValue:  "",
			defaultValue: "default_value",
		},
		{
			name:         "Get non-existing environment variable without default value",
			key:          "NON_EXISTING_ENV",
			expected:     "",
			setEnvValue:  "",
			defaultValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnvValue != "" {
				os.Setenv(tt.key, tt.setEnvValue)
				defer os.Unsetenv(tt.key) // Clean up after test
			}

			envService := &EnvServiceImpl{}
			result := envService.Get(tt.key, tt.defaultValue)

			assert.Equal(t, tt.expected, result)
		})
	}
}
