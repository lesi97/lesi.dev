package bungie_store

import (
	"context"
)

type userCharacters struct {
	MembershipID   string
	CharacterID    string
	CharacterType  string
	MinutesPlayed  string
}

func (s *BungieStore) insertUserCharacters(user *userCharacters) {
	query := `
		INSERT INTO destiny_user_characters
			(membership_id, character_id, character_type, minutes_played)
		VALUES 
			($1, $2, $3, $4)
		ON CONFLICT (character_id)
		DO UPDATE SET 
			membership_id = EXCLUDED.membership_id,
			character_type = EXCLUDED.character_type,
			minutes_played = EXCLUDED.minutes_played
	`
	_, err := s.DB.Exec(context.Background(), query,
		user.MembershipID, user.CharacterID, user.CharacterType, user.MinutesPlayed)

	if err != nil {
		s.Logger.Printf("insertUserCharacters failed: %v", err)
	}
}
