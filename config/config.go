package config

import (
	"errors"
	"sync"
)

type Config struct {
	ServerToken string
}

var once sync.Once
var instance *Config

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func SetServerToken(token string) {
	GetConfig().ServerToken = token
}

func GetServerToken() (string, error) {
	if instance == nil || instance.ServerToken == "" {
		return "", errors.New("server token is not set")
	}
	return instance.ServerToken, nil
}
