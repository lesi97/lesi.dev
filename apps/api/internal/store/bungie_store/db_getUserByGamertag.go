package bungie_store

import (
	"context"
	"database/sql"
	"time"
)

type bungieDBData struct {
	BungieID			string	`json:"bungie_id"`
	MembershipID 		string 	`json:"membership_id"`
	PreferredPlatform 	int64 	`json:"preferred_platform"`
	FriendlyName		string 	`json:"friendly_name"`
}

func (s *BungieStore) getUserFromDatabaseByGamertag(ctx context.Context, bungieID string) (*bungieDBData, error) {
	defer s.Logger.LogExecutionTime("DATABASE CALL: getUserFromDatabaseByGamertag", time.Now(), ctx)
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
	err := s.DB.QueryRow(ctx, query, bungieID).Scan(
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