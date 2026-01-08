package bungie_store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
)

type weaponData struct {
	HashID				string 	`json:"hash_id"`
	DisplayName			string	`json:"display_name"`
	TierTypeName 		string 	`json:"tier_type_name"`
}

func (s *BungieStore) getWeaponData(ctx context.Context, hashID string) (*weaponData, error) {
	defer utils.LogExecutionTime("getWeaponData", time.Now())
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
	
	err := s.DB.QueryRow(ctx, query, hashID).Scan(
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