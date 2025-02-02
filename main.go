package main

import (
	"context"
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
	// Initialize logger
	utils.InitializeLogger(slog.LevelDebug)

	// Load environment variables
	godotenv.Load()
	redisUrl := utils.GetEnv("REDIS_URL", "localhost:6379")
	appPort := utils.GetEnv("APP_PORT", "8888")

	// Set up Redis client
	redisClient, err := utils.SetupRedis(redisUrl)
	if err != nil {
		utils.LogError(context.Background(), "failed to connect to redis", "", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	// Generate and set server token
	serverToken := utils.RandomToken()
	config.SetServerToken(serverToken)

	// Set up server and routes
	router := gin.Default()
	router.Use(middlewares.LoggingMiddleware())
	router.Use(middlewares.RequestIDMiddleware())
	routes.TokenRoute(router)
	routes.SecretRoutes(router)

	// Start the server
	appPort = fmt.Sprintf(":%v", appPort)
	router.Run(appPort)
}
