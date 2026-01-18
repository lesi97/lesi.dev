package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
)

type TwitchFrontendAuthStartResult struct {
	URL          string
	State        string
	PKCEVerifier string
}

func randomB64URL(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func (s *AuthStore) TwitchFrontendAuthStart() (*TwitchFrontendAuthStartResult, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%v/authorize", s.twitch.base_url))
	if err != nil {
		return nil, err
	}

	state, err := randomB64URL(32)
	if err != nil {
		return nil, err
	}

	verifier, err := randomB64URL(64)
	if err != nil {
		return nil, err
	}

	q := baseURL.Query()
	q.Set("client_id", *s.twitch.client_id)
	q.Set("redirect_uri", fmt.Sprintf("%v/api/auth/twitch/callback", *s.webUrl))
	q.Set("response_type", "code")
	q.Set("state", state)
	q.Set("code_challenge", pkceChallenge(verifier))
	q.Set("code_challenge_method", "S256")

	baseURL.RawQuery = q.Encode()

	return &TwitchFrontendAuthStartResult{
		URL:          baseURL.String(),
		State:        state,
		PKCEVerifier: verifier,
	}, nil
}
