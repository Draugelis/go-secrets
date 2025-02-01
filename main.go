package main

import (
	"fmt"
	"go-secrets/config"
	"go-secrets/middlewares"
	"go-secrets/routes"
	"go-secrets/utils"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Start logging
	utils.InitializeLogger(slog.LevelDebug)

	// Load environment variables
	godotenv.Load()
	redisUrl := utils.GetEnv("REDIS_URL", "localhost:6379")
	appPort := utils.GetEnv("APP_PORT", "8888")

	// Initialize redis client
	if _, err := utils.SetupRedis(redisUrl); err != nil {
		slog.Error("failed to connect to Redis", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize server token
	serverToken := utils.RandomToken()
	config.SetServerToken(serverToken)

	// Initialize server
	router := gin.Default()
	router.Use(middlewares.LoggingMiddleware())
	routes.TokenRoute(router)
	routes.SecretRoutes(router)

	appPort = fmt.Sprintf(":%v", appPort)
	router.Run(appPort)
}
