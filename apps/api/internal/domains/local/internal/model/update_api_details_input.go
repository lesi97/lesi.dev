package model

type UpdateApiDetailsInput struct {
	Name               string
	ClientID           string
	ClientSecret       string
	AccessToken        string
	RefreshToken       string
	RefreshTokenExpiry int64
	BaseURL            string
	RedirectURL        string
}
