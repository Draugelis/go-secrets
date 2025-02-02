package utils

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthTokenHMAC(ctx *gin.Context) (string, error) {
	token := GetHeaderToken(ctx)
	hmac, err := HMAC(token)
	if err != nil {
		slog.Error("failed to get token hmac", slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get token hmac"})
		return "", err
	}

	return hmac, nil
}
