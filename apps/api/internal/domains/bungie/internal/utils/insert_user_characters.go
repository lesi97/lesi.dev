package utils

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func InsertUserCharacters(database *db.DB, logger *utils.Logger, membershipID string, characterID string, characterType string, minutesPlayed string) {
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
	_, err := database.Exec(context.Background(), query, membershipID, characterID, characterType, minutesPlayed)

	if err != nil {
		logger.Printf("insertUserCharacters failed: %v", err)
	}
}
