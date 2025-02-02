package utils

import (
	"github.com/gin-gonic/gin"
)

func AuthTokenHMAC(ctx *gin.Context) (string, error) {
	token := GetHeaderToken(ctx)
	hmac, err := HMAC(token)
	if err != nil {
		return "", err
	}

	return hmac, nil
}
