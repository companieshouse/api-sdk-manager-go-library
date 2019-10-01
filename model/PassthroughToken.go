package model

type PassthroughToken struct{
	TokenType string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn int64 `json:"expires_in"`
}