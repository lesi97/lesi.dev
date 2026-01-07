package countdown_store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type CountdownPostRequest struct {
	TargetDate      time.Time `json:"target_date"`
	Message         string    `json:"message"`
	FallbackMessage string    `json:"fallback_message"`
}

func (supabase *SupabaseCountdownStore) InsertCountdown(ctx context.Context, data CountdownPostRequest) (*string, error) {
	if data.TargetDate.IsZero() {
		return nil, errors.New("target_date is required")
	}
	if data.Message == "" {
		return nil, errors.New("message is required")
	}
	if data.FallbackMessage == "" {
		return nil, errors.New("fallback_message is required")
	}

	query := `
		INSERT INTO countdown (target_date, message, fallback_message)
		VALUES ($1, $2, $3)
		RETURNING uuid
	`

	var uuid string
	err := supabase.db.QueryRow(ctx, query,
		data.TargetDate,
		data.Message,
		data.FallbackMessage,
	).Scan(&uuid)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &uuid, nil
}