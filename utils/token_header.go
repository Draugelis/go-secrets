package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Get token from Authorization header
func GetHeaderToken(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	token := parts[1]

	return token
}
