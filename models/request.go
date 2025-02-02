package models

// StoreSecretRequest represents the request payload for storing a secret, containing the value of the secret.
type StoreSecretRequest struct {
	Value string `json:"value" binding:"required"`
}
