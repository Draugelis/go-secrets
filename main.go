package main

import (
	"go-secrets/config"
	"go-secrets/routes"
	"go-secrets/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()
	redisUrl := utils.GetEnv("REDIS_URL", "localhost:6379")

	// Initialize redis client
	_, err := utils.SetupRedis(redisUrl)
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		os.Exit(1)
	}

	// Initialize server
	router := gin.Default()
	routes.TokenRoute(router)
	routes.SecretRoutes(router)

	// Initialize server token
	serverToken := utils.RandomToken()
	config.SetServerToken(serverToken)

	router.Run(":8888")
}
