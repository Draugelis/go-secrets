package helpers

import (
	"fmt"
)

// FormatSecretPath formats the secret path by combining HMAC and key into a single string.
func FormatSecretPath(namespace string, key string) (string, error) {
	if namespace == "" || key == "" {
		return "", fmt.Errorf("namespace and key cannot be empty")
	}
	return fmt.Sprintf("%s:secret:%s", namespace, key), nil
}
