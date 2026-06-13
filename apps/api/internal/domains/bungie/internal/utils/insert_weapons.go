package utils

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func insertWeapons(database *db.DB, logger *utils.Logger, weapons []dbWeapons) {
	for _, weapon := range weapons {
		query := `
			INSERT INTO destiny_weapons 
				(id, display_name, item_type_display_name, flavor_text, bucket_type_hash, tier_type_hash, tier_type_name, tier_type, talent_grid_hash)
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id)
			DO UPDATE SET
				display_name = EXCLUDED.display_name,
				item_type_display_name = EXCLUDED.item_type_display_name,
				flavor_text = EXCLUDED.flavor_text,
				bucket_type_hash = EXCLUDED.bucket_type_hash,
				tier_type_hash = EXCLUDED.tier_type_hash,
				tier_type_name = EXCLUDED.tier_type_name,
				tier_type = EXCLUDED.tier_type,
				talent_grid_hash = EXCLUDED.talent_grid_hash
		`

		_, err := database.Exec(context.Background(), query,
			weapon.ID,
			weapon.DisplayName,
			weapon.ItemTypeDisplayName,
			weapon.FlavorText,
			weapon.BucketTypeHash,
			weapon.TierTypeHash,
			weapon.TierTypeName,
			weapon.TierType,
			weapon.TalentGridHash,
		)
		if err != nil {
			logger.Printf("insertWeapons failed for ID %s: %v\n", weapon.ID, err)
		}
	}
}
