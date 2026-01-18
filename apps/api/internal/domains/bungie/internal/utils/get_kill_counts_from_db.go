package utils

import (
	"context"
	"database/sql"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type KillCountsDBData struct {
	PVPKills int `json:"pvp_kills"`
}

func GetKillCountsFromDB(ctx context.Context, database *db.DB, logger *utils.Logger, bungieID string, weaponID string) (*KillCountsDBData, error) {
	defer logger.LogExecutionTime("DATABASE CALL: getKillCountsFromDB", time.Now(), ctx)
	query := `
		SELECT pvp_kills 
		FROM destiny_weapon_kill_counts 
		WHERE membership_id = $1
		AND weapon_id = $2
	`
	var data KillCountsDBData
	err := database.QueryRow(ctx, query, bungieID, weaponID).Scan(&data.PVPKills)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &data, nil
}
