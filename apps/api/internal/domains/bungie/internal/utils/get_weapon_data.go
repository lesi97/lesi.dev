package utils

import (
	"context"
	"database/sql"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func getWeaponData(ctx context.Context, database *db.DB, logger *utils.Logger, hashID string) (*weaponData, error) {
	defer logger.LogExecutionTime("DATABASE CALL: getWeaponData", time.Now(), ctx)
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

	err := database.QueryRow(ctx, query, hashID).Scan(
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
