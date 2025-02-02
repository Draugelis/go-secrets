package errors

import (
	"go-secrets/models"
	"net/http"
)

var (
	ErrInvalidRequest = models.NewErrorResponse(http.StatusBadRequest, "invalid request")
	ErrUnauthorized   = models.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
	ErrNotFound       = models.NewErrorResponse(http.StatusNotFound, "resource not found")
	ErrInternalServer = models.NewErrorResponse(http.StatusInternalServerError, "internal server error")
	ErrAPIMissingPath = models.NewErrorResponse(http.StatusBadRequest, "missing key path")
)
