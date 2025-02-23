package models

// GetSecretResponse represents the response payload for retrieving a secret, containing the secret's value and its time-to-live (TTL).
// @Description Get secret response format
// @Example { "value": "my_secret_value", "ttl": 3600 }
type GetSecretResponse struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

// StoreSecretResponse represents the response payload for storing a secret, containing the generated key and its time-to-live (TTL).
// @Description Store secret response format
// @Example { "key": "abc123", "ttl": 3600 }
type StoreSecretResponse struct {
	Key string `json:"key"`
	TTL int    `json:"ttl"`
}

// IssueTokenResponse represents the response payload for issuing a token, containing the token string and its time-to-live (TTL).
// @Description Issue token response format
// @Example { "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", "ttl": 7200 }
type IssueTokenResponse struct {
	Token string `json:"token"`
	TTL   int    `json:"ttl"`
}

// TokenValidationResponse represents the response payload for validating a token.
// @Description Token validation response format
type TokenValidationResponse struct {
	Valid bool `json:"valid"`
}
