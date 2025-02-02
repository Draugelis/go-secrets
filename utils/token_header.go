package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GetHeaderToken retrieves the token from the Authorization header in the request.
func GetHeaderToken(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	token := parts[1]

	return token
}
