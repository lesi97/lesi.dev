package store

import (
	"context"
	"fmt"

	twitch_frontend_utils "github.com/lesi97/lesi.dev/internal/domains/auth/utils/twitch_frontend"
)

type TwitchFrontendIdentity struct {
	ID          string
	Login       string
	DisplayName string
	AvatarURL   string
}

func (s *Store) TwitchFrontendCallback(
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

	redirectURL := fmt.Sprintf("%v/api/auth/twitch/callback", *s.webUrl)
	accessToken, err := twitch_frontend_utils.ExchangeCodeForToken(
		ctx,
		s.twitch.base_url,
		*s.twitch.client_id,
		*s.twitch.client_secret,
		redirectURL,
		code,
		pkceVerifier,
	)
	if err != nil {
		return nil, err
	}

	identity, err := twitch_frontend_utils.FetchIdentity(ctx, accessToken, *s.twitch.client_id)
	if err != nil {
		return nil, err
	}

	return &TwitchFrontendIdentity{
		ID:          identity.ID,
		Login:       identity.Login,
		DisplayName: identity.DisplayName,
		AvatarURL:   identity.AvatarURL,
	}, nil
}
