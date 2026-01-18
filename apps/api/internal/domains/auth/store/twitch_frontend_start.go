package store

import (
	"fmt"
	"net/url"

	twitch_frontend_utils "github.com/lesi97/lesi.dev/internal/domains/auth/utils/twitch_frontend"
)

type TwitchFrontendAuthStartResult struct {
	URL          string
	State        string
	PKCEVerifier string
}

func (s *Store) TwitchFrontendAuthStart() (*TwitchFrontendAuthStartResult, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%v/authorize", s.twitch.base_url))
	if err != nil {
		return nil, err
	}

	state, err := twitch_frontend_utils.RandomB64URL(32)
	if err != nil {
		return nil, err
	}

	verifier, err := twitch_frontend_utils.RandomB64URL(64)
	if err != nil {
		return nil, err
	}

	q := baseURL.Query()
	q.Set("client_id", *s.twitch.client_id)
	q.Set("redirect_uri", fmt.Sprintf("%v/api/auth/twitch/callback", *s.webUrl))
	q.Set("response_type", "code")
	q.Set("state", state)
	q.Set("code_challenge", twitch_frontend_utils.PKCEChallenge(verifier))
	q.Set("code_challenge_method", "S256")

	baseURL.RawQuery = q.Encode()

	return &TwitchFrontendAuthStartResult{
		URL:          baseURL.String(),
		State:        state,
		PKCEVerifier: verifier,
	}, nil
}
