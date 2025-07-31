package bungie_store

import (
	"context"
	"database/sql"
)

type weaponData struct {
	HashID				int64 	`json:"preferred_platform"`
	DisplayName			string	`json:"display_name"`
	TierTypeName 		string 	`json:"tier_type_name"`
}

func (supabase *SupabaseBungieStoreStore) getWeaponData(ctx context.Context, hashID int64) (*weaponData, error) {
	query := `
		SELECT 
			display_name, 
			tier_type_name
		FROM destiny_weapons 
		WHERE id = $1
	`
	data := &weaponData{
		HashID: hashID,
	}
	
	err := supabase.db.QueryRow(ctx, query, hashID).Scan(
		&data.DisplayName,
		&data.TierTypeName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}