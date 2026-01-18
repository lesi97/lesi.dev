package utils

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func InsertKillCounts(
	database *db.DB,
	logger *utils.Logger,
	membershipID string,
	weaponID string,
	weaponHash string,
	pvpKills *int,
	pveKills *int,
	trialsKills *int,
) {
	query := `
		INSERT INTO destiny_weapon_kill_counts 
			(membership_id, weapon_id, pvp_kills, pve_kills, trials_kills, weapon_hash)
		VALUES 
			($1, $2, $3, $4, $5, $6)
		ON CONFLICT (membership_id, weapon_id) DO UPDATE
		SET 
			pvp_kills = COALESCE(EXCLUDED.pvp_kills, destiny_weapon_kill_counts.pvp_kills),
			pve_kills = COALESCE(EXCLUDED.pve_kills, destiny_weapon_kill_counts.pve_kills),
			trials_kills = COALESCE(EXCLUDED.trials_kills, destiny_weapon_kill_counts.trials_kills),
			weapon_hash = EXCLUDED.weapon_hash
	`
	_, err := database.Exec(context.Background(), query,
		membershipID,
		weaponID,
		pvpKills,
		pveKills,
		trialsKills,
		weaponHash,
	)
	if err != nil {
		logger.Printf("insertKillCounts failed: %v", err)
	}
}
