package utils

import "fmt"

// FormatSecretPath formats the secret path by combining HMAC and key into a single string.
func FormatSecretPath(HMAC string, key string) (string, error) {
	if HMAC == "" || key == "" {
		return "", fmt.Errorf("HMAC and key cannot be empty")
	}
	return fmt.Sprintf("%s:secret:%s", HMAC, key), nil
}
