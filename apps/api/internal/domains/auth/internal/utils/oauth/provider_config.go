package oauth

type ProviderConfig struct {
	Application  string
	TokenURL     string
	ClientID     *string
	ClientSecret *string
	RedirectURL  string
}
