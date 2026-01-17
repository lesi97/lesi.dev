package auth_store

import (
	"context"
	"fmt"
)

func (s *AuthStore) TwitchFrontendDeleteSessionByToken(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return fmt.Errorf("missing session token")
	}

	_, err := s.DB.Exec(
		ctx,
		`delete from public.app_session where session_token = $1`,
		sessionToken,
	)
	return err
}
