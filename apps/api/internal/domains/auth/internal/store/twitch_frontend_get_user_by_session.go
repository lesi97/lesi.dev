package store

import (
	"context"
	"fmt"
	"time"
)

type TwitchFrontendUser struct {
	ID               string  `json:"id"`
	TwitchUserID     string  `json:"twitchUserId"`
	TwitchLogin      string  `json:"twitchLogin"`
	TwitchDisplayName *string `json:"twitchDisplayName"`
	TwitchAvatarURL  *string `json:"twitchAvatarUrl"`
}

func (s *Store) TwitchFrontendGetUserBySession(ctx context.Context, sessionToken string) (*TwitchFrontendUser, error) {
	if sessionToken == "" {
		return nil, fmt.Errorf("missing session token")
	}

	var u TwitchFrontendUser
	var expiresAt time.Time

	err := s.DB.QueryRow(
		ctx,
		`
		select
			u.id,
			u.twitch_user_id,
			u.twitch_login,
			u.twitch_display_name,
			u.twitch_avatar_url,
			s.expires_at
		from public.app_session s
		join public.app_user u on u.id = s.user_id
		where s.session_token = $1
		`,
		sessionToken,
	).Scan(&u.ID, &u.TwitchUserID, &u.TwitchLogin, &u.TwitchDisplayName, &u.TwitchAvatarURL, &expiresAt)

	if err != nil {
		return nil, err
	}

	if time.Now().After(expiresAt) {
		_ = s.TwitchFrontendDeleteSessionByToken(ctx, sessionToken)
		return nil, fmt.Errorf("session expired")
	}

	return &u, nil
}
