package utils

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func InsertDestinyUser(database *db.DB, logger *utils.Logger, membershipID string, bungieID string, preferredPlatform int64, friendlyName string) {
	query := `
		INSERT INTO destiny_users 
			(membership_id, bungie_id, preferred_platform, friendly_name)
		VALUES 
			($1, $2, $3, $4)
		ON CONFLICT (bungie_id) DO NOTHING
	`
	_, err := database.Exec(context.Background(), query, membershipID, bungieID, preferredPlatform, friendlyName)
	if err != nil {
		logger.Printf("insertUser failed: %v", err)
	}
}
