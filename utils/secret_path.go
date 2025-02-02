package utils

import "fmt"

// FormatSecretPath formats the secret path by combining HMAC and key into a single string.
func FormatSecretPath(HMAC string, key string) string {
	return fmt.Sprintf("%s:secret:%s", HMAC, key)
}
