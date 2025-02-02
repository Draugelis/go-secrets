package models

type StoreSecretRequest struct {
	Value string `json:"value" binding:"required"`
}
