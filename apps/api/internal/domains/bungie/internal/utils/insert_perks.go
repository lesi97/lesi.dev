package utils

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func insertPerks(database *db.DB, logger *utils.Logger, perks []dbPerk) {
	for _, perk := range perks {
		query := `
			INSERT INTO destiny_weapon_perks
				(name, description, item_type, hash_id)
			VALUES 
				($1, $2, $3, $4)
			ON CONFLICT (hash_id)
			DO UPDATE SET
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				item_type = EXCLUDED.item_type
		`
		_, err := database.Exec(context.Background(), query, perk.Name, perk.Description, perk.ItemType, perk.HashID)
		if err != nil {
			logger.Printf("insertPerks failed for %s: %v\n", perk.Name, err)
		}
	}
}
