package utils

import (
	"go-secrets/errors"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func AuthTokenHMAC(ctx *gin.Context) (string, error) {
	token := GetHeaderToken(ctx)
	hmac, err := HMAC(token)
	if err != nil {
		slog.Error("failed to get token hmac", slog.String("error", err.Error()))
		errors.ErrInternalServer.WithRequestID(ctx).JSON(ctx)
		return "", err
	}

	return hmac, nil
}
