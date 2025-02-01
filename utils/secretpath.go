package utils

import "fmt"

func FormatSecretPath(HMAC string, key string) string {
	return fmt.Sprintf("%s:secret:%s", HMAC, key)
}
