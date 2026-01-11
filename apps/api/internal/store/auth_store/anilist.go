package auth_store

import (
	"context"
	"fmt"
	"net/url"
)

func (s *AuthStore) AnilistCallback(code string) error {

	cfg := oauthProviderConfig{
		Application: "Anilist",
		TokenURL: fmt.Sprintf("%v/token", s.anilist.base_url),
		ClientID: s.anilist.api_details.ClientID,
		ClientSecret: s.anilist.api_details.ClientSecret,
		RedirectURL: s.anilist.api_details.RedirectURL,
	}

	return s.oauthCallback(context.Background(), code, cfg)
}


func (s *AuthStore) AnilistAuthUrl() (*string, error) {

	baseURL, err := url.Parse(fmt.Sprintf("%v/authorize", s.anilist.base_url))
	if err != nil {
		return nil, err
	}
	query := baseURL.Query()
	query.Set("client_id", *s.anilist.api_details.ClientID)
	query.Set("redirect_uri", *s.anilist.api_details.RedirectURL)
	query.Set("response_type", "code")

	baseURL.RawQuery = query.Encode()
	str := baseURL.String()
	return &str, nil
}