package bungie_store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lesi97/api.lesi.dev/internal/utils"
)

type killCountsDBData struct {
	PVPKills			int		`json:"pvp_kills"`
}

func (supabase *SupabaseBungieStore) getKillCountsFromDB(ctx context.Context, bungieID string, weaponID string) (*killCountsDBData, error) {
	defer utils.LogExecutionTime("getKillCountsFromDB", time.Now())
	query := `
		SELECT pvp_kills 
		FROM destiny_weapon_kill_counts 
		WHERE membership_id = $1
		AND weapon_id = $2
	`
	var data killCountsDBData
	err := supabase.db.QueryRow(ctx, query, bungieID, weaponID).Scan(&data.PVPKills)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &data, nil
}