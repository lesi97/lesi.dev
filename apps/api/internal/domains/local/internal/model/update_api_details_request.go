package model

type UpdateApiDetailsRequest struct {
	Name               string `json:"name"`
	ClientID           string `json:"client_id"`
	ClientSecret       string `json:"client_secret"`
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	RefreshTokenExpiry int64  `json:"refresh_token_expiry"`
	BaseURL            string `json:"base_url"`
	RedirectURL        string `json:"redirect_url"`
}
