package internal

import "fmt"

type SecretPathService interface {
	FormatSecretPath(hmac, key string) (string, error)
}

type SecretPathServiceImpl struct{}

// FormatSecretPath formats the secret path by combining HMAC and key into a single string.
func (s *SecretPathServiceImpl) FormatSecretPath(hmac string, key string) (string, error) {
	if hmac == "" || key == "" {
		return "", fmt.Errorf("hmac and key cannot be empty")
	}
	return fmt.Sprintf("%s:secret:%s", hmac, key), nil
}
