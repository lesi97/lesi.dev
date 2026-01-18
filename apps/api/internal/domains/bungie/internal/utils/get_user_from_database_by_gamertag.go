package utils

import (
	"context"
	"database/sql"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type bungieDBData struct {
	BungieID          string `json:"bungie_id"`
	MembershipID      string `json:"membership_id"`
	PreferredPlatform int64  `json:"preferred_platform"`
	FriendlyName      string `json:"friendly_name"`
}

func getUserFromDatabaseByGamertag(ctx context.Context, database *db.DB, logger *utils.Logger, bungieID string) (*bungieDBData, error) {
	defer logger.LogExecutionTime("DATABASE CALL: getUserFromDatabaseByGamertag", time.Now(), ctx)
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
	err := database.QueryRow(ctx, query, bungieID).Scan(
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
