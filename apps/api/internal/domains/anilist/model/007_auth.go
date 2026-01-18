package model

import "time"

type OAuthClientConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type OAuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    *time.Time
	TokenType    string
}
