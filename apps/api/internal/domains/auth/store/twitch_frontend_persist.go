package store

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

func randomTokenB64URL(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *AuthStore) TwitchFrontendUpsertUser(ctx context.Context, identity TwitchFrontendIdentity) (string, error) {
	if identity.ID == "" || identity.Login == "" {
		return "", fmt.Errorf("missing twitch identity")
	}

	var userID string
	err := s.DB.QueryRow(
		ctx,
		`
		insert into public.app_user (twitch_user_id, twitch_login, twitch_display_name, twitch_avatar_url)
		values ($1, $2, $3, $4)
		on conflict (twitch_user_id) do update set
			twitch_login = excluded.twitch_login,
			twitch_display_name = excluded.twitch_display_name,
			twitch_avatar_url = excluded.twitch_avatar_url
		returning id
		`,
		identity.ID,
		identity.Login,
		identity.DisplayName,
		identity.AvatarURL,
	).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *AuthStore) TwitchFrontendCreateSession(ctx context.Context, userID string, ttl time.Duration) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("missing user id")
	}
	if ttl <= 0 {
		return "", fmt.Errorf("invalid ttl")
	}

	token, err := randomTokenB64URL(32)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(ttl)

	_, err = s.DB.Exec(
		ctx,
		`
		insert into public.app_session (user_id, session_token, expires_at)
		values ($1, $2, $3)
		`,
		userID,
		token,
		expiresAt,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}
