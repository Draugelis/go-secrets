package helpers

import (
	"errors"
	"os"
)

// GetEnv retrieves the environment variable or returns a default value.
func GetEnv(key string, defaultValue ...string) (string, error) {
	value := os.Getenv(key)
	if value != "" {
		return value, nil
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return "", errors.New("environment variable " + key + " is not set and no default value provided")
}
