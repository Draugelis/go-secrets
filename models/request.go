package models

// StoreSecretRequest represents the request payload for storing a secret, containing the value of the secret.
// @Description Store secret request format
// @Example { "value": "my_secret_value" }
type StoreSecretRequest struct {
	Value string `json:"value" binding:"required"`
}
