package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Run("returns existing environment variable", func(t *testing.T) {
		os.Setenv("TEST_KEY", "exists")
		value, err := GetEnv("TEST_KEY")

		assert.NoError(t, err)
		assert.Equal(t, "exists", value)

		os.Unsetenv("TEST_KEY") // Cleanup after test
	})

	t.Run("returns default value when env var is missing", func(t *testing.T) {
		os.Unsetenv("TEST_KEY")
		value, err := GetEnv("TEST_KEY", "default")

		assert.NoError(t, err)
		assert.Equal(t, "default", value)
	})

	t.Run("returns error when env var is missing and no default is provided", func(t *testing.T) {
		os.Unsetenv("TEST_KEY")
		value, err := GetEnv("TEST_KEY")

		assert.Error(t, err)
		assert.Empty(t, value)
		assert.EqualError(t, err, "environment variable TEST_KEY is not set and no default value provided")
	})
}
