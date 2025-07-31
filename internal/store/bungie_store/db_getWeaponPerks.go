package bungie_store

import (
	"context"
	"database/sql"
)

type weaponPerks struct {
	Name				string 	`json:"name"`
	ItemType			string	`json:"item_type"`
}

func (supabase *SupabaseBungieStoreStore) getWeaponPerks(ctx context.Context, hashIDs []int64) (*weaponPerks, error) {
	query := `
		SELECT 
			name, 
			item_type
		FROM destiny_weapon_perks 
		WHERE hash_id = ANY($1)
	`
	data := &weaponPerks{}
	
	err := supabase.db.QueryRow(ctx, query, hashIDs).Scan(
		&data.Name,
		&data.ItemType,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}