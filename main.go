package main

import (
	"context"
	"fmt"
	"go-secrets/config"
	"go-secrets/helpers"
	"go-secrets/internal"
	"go-secrets/middlewares"
	"go-secrets/routes"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	internal.SetupLogger(slog.LevelDebug, os.Stdout)
	logger, err := internal.GetLoggerService()
	if err != nil {
		os.Exit(1)
	}

	// Load environment variables from .env file (if present)
	_ = godotenv.Load()
	redisURL, _ := helpers.GetEnv("REDIS_URL", "localhost:6379")
	appPort, _ := helpers.GetEnv("APP_PORT", "8888")

	// Set up Redis client
	if err := internal.SetupRedis(redisURL); err != nil {
		logger.LogError(context.Background(), "Failed to connect to Redis", "", err)
		os.Exit(1)
	}

	// Generate and set server token
	tokenService := internal.NewTokenService()
	serverToken, err := tokenService.GenerateToken()
	if err != nil {
		logger.LogError(context.Background(), "Failed to generate server token", "", err)
		os.Exit(1)
	}
	config.SetServerToken(serverToken)

	// Initialize services
	cryptoService := internal.NewCryptoService()
	redisClient, err := internal.GetRedisService()
	if err != nil {
		logger.LogError(context.Background(), "Failed to get Redis client", "", err)
		os.Exit(1)
	}

	// Set up router and middleware
	router := gin.Default()
	router.Use(middlewares.LoggingMiddleware())
	router.Use(middlewares.RequestIDMiddleware())

	// Register routes
	routes.TokenRoute(router, logger, cryptoService, redisClient, tokenService)
	routes.SecretRoutes(router, logger, cryptoService, redisClient, tokenService)

	// Start the server
	serverAddress := fmt.Sprintf(":%s", appPort)
	if err := router.Run(serverAddress); err != nil {
		logger.LogError(context.Background(), "Failed to start server", "", err)
		os.Exit(1)
	}
}
