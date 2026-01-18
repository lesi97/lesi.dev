package store

import (
	"context"
	"fmt"
	"time"

	twitch_frontend_utils "github.com/lesi97/lesi.dev/internal/domains/auth/utils/twitch_frontend"
)

func (s *Store) TwitchFrontendCreateSession(ctx context.Context, userID string, ttl time.Duration) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("missing user id")
	}
	if ttl <= 0 {
		return "", fmt.Errorf("invalid ttl")
	}

	token, err := twitch_frontend_utils.RandomTokenB64URL(32)
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
