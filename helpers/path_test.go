package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSecretPath(t *testing.T) {
	t.Run("formats secret path correctly", func(t *testing.T) {
		namespace := "my_namespace"
		key := "my_key"
		expectedPath := "my_namespace:secret:my_key"

		result, err := FormatSecretPath(namespace, key)

		assert.NoError(t, err)
		assert.Equal(t, expectedPath, result)
	})

	t.Run("returns error when namespace is empty", func(t *testing.T) {
		key := "my_key"

		result, err := FormatSecretPath("", key)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.EqualError(t, err, "namespace and key cannot be empty")
	})

	t.Run("returns error when key is empty", func(t *testing.T) {
		namespace := "my_namespace"

		result, err := FormatSecretPath(namespace, "")

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.EqualError(t, err, "namespace and key cannot be empty")
	})

	t.Run("returns error when both namespace and key are empty", func(t *testing.T) {
		result, err := FormatSecretPath("", "")

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.EqualError(t, err, "namespace and key cannot be empty")
	})
}
