package utils

import (
	"github.com/gin-gonic/gin"
)

// AuthTokenHMAC computes the HMAC of the token from the request header.
func AuthTokenHMAC(ctx *gin.Context) (string, error) {
	token := GetHeaderToken(ctx)
	hmac, err := HMAC(token)
	if err != nil {
		return "", err
	}

	return hmac, nil
}
