package bungie_store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lesi97/api.lesi.dev/internal/utils"
)

type bungieDBData struct {
	BungieID			string	`json:"bungie_id"`
	MembershipID 		string 	`json:"membership_id"`
	PreferredPlatform 	int64 	`json:"preferred_platform"`
	FriendlyName		string 	`json:"friendly_name"`
}

func (supabase *SupabaseBungieStoreStore) getUserFromDatabaseByGamertag(ctx context.Context, bungieID string) (*bungieDBData, error) {
	defer utils.LogExecutionTime("getUserFromDatabaseByGamertag", time.Now())
	query := `
		SELECT 
			membership_id, 
			preferred_platform, 
			friendly_name 
		FROM destiny_users 
		WHERE lower(bungie_id) = lower($1)
	`
	data := &bungieDBData{
		BungieID: bungieID,
	}
	err := supabase.db.QueryRow(ctx, query, bungieID).Scan(
		&data.MembershipID,
		&data.PreferredPlatform,
		&data.FriendlyName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}