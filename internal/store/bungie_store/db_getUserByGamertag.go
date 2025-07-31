package bungie_store

import (
	"context"
	"database/sql"
)

type bungieDBData struct {
	BungieID			string	`json:"bungie_id"`
	MembershipID 		string 	`json:"membership_id"`
	PreferredPlatform 	int64 	`json:"preferred_platform"`
	FriendlyName		string 	`json:"friendly_name"`
}

func (supabase *SupabaseBungieStoreStore) getUserFromDatabaseByGamertag(ctx context.Context, bungieID string) (*bungieDBData, error) {
	query := `
		SELECT 
			membership_id, 
			preferred_platform, 
			friendly_name 
		FROM destiny_users 
		WHERE bungie_id ILIKE $1
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