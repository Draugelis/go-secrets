package controllers

import (
	"go-secrets/internal"

	"github.com/gin-gonic/gin"
)

type TokenController interface {
	Generate(ctx *gin.Context)
	Validate(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type TokenControllerImpl struct {
	Logger internal.LoggerService
	Crypto internal.CryptoService
	Redis  internal.RedisService
	Token  internal.TokenService
}
