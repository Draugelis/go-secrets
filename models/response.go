package models

// GetSecretResponse represents the response payload for retrieving a secret, containing the secret's value and its time-to-live (TTL).
type GetSecretResponse struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

// StoreSecretResponse represents the response payload for storing a secret, containing the generated key and its time-to-live (TTL).
type StoreSecretResponse struct {
	Key string `json:"key"`
	TTL int    `json:"ttl"`
}

// IssueTokenResponse represents the response payload for issuing a token, containing the token string and its time-to-live (TTL).
type IssueTokenResponse struct {
	Token string `json:"token"`
	TTL   int    `json:"ttl"`
}
