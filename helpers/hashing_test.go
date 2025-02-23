package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHMAC(t *testing.T) {
	t.Run("generates correct HMAC", func(t *testing.T) {
		data := "test_data"
		secret := "test_secret"
		expectedHMAC := "6a4848557f6dff818d54e7563310e2d11a06768ef39c7ae003b95737d4f1d2cb"

		result, err := GenerateHMAC(data, secret)

		assert.NoError(t, err)
		assert.Equal(t, expectedHMAC, result)
	})
}
