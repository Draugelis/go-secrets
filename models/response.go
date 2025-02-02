package models

type GetSecretResponse struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

type StoreSecretResponse struct {
	Key string `json:"key"`
	TTL int    `json:"ttl"`
}

type IssueTokenResponse struct {
	Token string `json:"token"`
	TTL   int    `json:"ttl"`
}
