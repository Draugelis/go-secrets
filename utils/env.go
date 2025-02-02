package utils

import "os"

// GetEnv retrieves the value of the environment variable specified by the key.
func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
