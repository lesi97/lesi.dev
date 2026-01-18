package store

import (
	"context"
	"fmt"
)

func (s *AuthStore) TwitchFrontendCallback(
	code string,
	state string,
	expectedState string,
	pkceVerifier string,
) (*TwitchFrontendIdentity, error) {
	if state == "" || expectedState == "" {
		return nil, fmt.Errorf("missing oauth state")
	}
	if state != expectedState {
		return nil, fmt.Errorf("oauth state mismatch")
	}
	if pkceVerifier == "" {
		return nil, fmt.Errorf("missing pkce verifier")
	}

	ctx := context.Background()

	token, err := s.twitchFrontendExchange(ctx, code, pkceVerifier)
	if err != nil {
		return nil, err
	}

	return s.fetchTwitchFrontendIdentity(ctx, token.AccessToken)
}
