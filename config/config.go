package config

import (
	"errors"
	"sync"
)

type ConfigService interface {
	SetServerToken(token string)
	GetServerToken() (string, error)
}

// Config holds the configuration data for the application, including the server token.
type Config struct {
	ServerToken string
}

var once sync.Once
var instance *Config

// GetConfig returns a singleton instance of the Config struct.
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

// SetServerToken sets the server token in the configuration.
func SetServerToken(token string) {
	GetConfig().ServerToken = token
}

// GetServerToken retrieves the server token from the configuration, returning an error if it is not set.
func GetServerToken() (string, error) {
	if instance == nil || instance.ServerToken == "" {
		return "", errors.New("server token is not set")
	}
	return instance.ServerToken, nil
}
