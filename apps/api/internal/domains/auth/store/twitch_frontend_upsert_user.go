package store

import (
	"context"
	"fmt"
)

func (s *Store) TwitchFrontendUpsertUser(ctx context.Context, identity TwitchFrontendIdentity) (string, error) {
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
