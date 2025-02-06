package internal

import "os"

type EnvService interface {
	Get(key string, defaultValue string) string
}

type EnvServiceImpl struct{}

// Get retrieves the value of the environment variable specified by the key.
func (e *EnvServiceImpl) Get(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
