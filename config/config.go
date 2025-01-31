package config

import (
	"log"
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

func GetServerToken() string {
	if instance == nil || instance.ServerToken == "" {
		log.Fatal("Server token is not set.")
	}
	return instance.ServerToken
}
