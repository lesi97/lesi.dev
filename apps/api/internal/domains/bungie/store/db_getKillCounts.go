package bungie_store

import (
	"context"
	"database/sql"
	"time"
)

type killCountsDBData struct {
	PVPKills			int		`json:"pvp_kills"`
}

func (s *BungieStore) getKillCountsFromDB(ctx context.Context, bungieID string, weaponID string) (*killCountsDBData, error) {
	defer s.Logger.LogExecutionTime("DATABASE CALL: getKillCountsFromDB", time.Now(), ctx)
	query := `
		SELECT pvp_kills 
		FROM destiny_weapon_kill_counts 
		WHERE membership_id = $1
		AND weapon_id = $2
	`
	var data killCountsDBData
	err := s.DB.QueryRow(ctx, query, bungieID, weaponID).Scan(&data.PVPKills)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &data, nil
}