package controllers

import (
	"go-secrets/internal"

	"github.com/gin-gonic/gin"
)

type SecretsController interface {
	Get(ctx *gin.Context)
	Set(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type SecretsControllerImpl struct {
	Logger internal.LoggerService
	Crypto internal.CryptoService
	Redis  internal.RedisService
	Token  internal.TokenService
}
