package auth_store

import (
	"fmt"
	"net/url"
)

func (s *AuthStore) GenerateAnilistAuthUrl() (*string, error) {
	
	baseURL, err := url.Parse(fmt.Sprintf("%v/authorize",s.anilist.base_url))
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