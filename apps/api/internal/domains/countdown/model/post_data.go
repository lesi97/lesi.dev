package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
)

type PostData struct {
	TargetDate      time.Time `json:"target_date"`
	Message         string    `json:"message"`
	FallbackMessage string    `json:"fallback_message"`
}

func (d *PostData) Validate() error {
	if d.TargetDate.IsZero() {
		return errors.New("target_date is required")
	}

	if d.Message == "" {
		return errors.New("message is required")
	}

	if d.FallbackMessage == "" {
		return errors.New("fallback_message is required")
	}

	if d.TargetDate.Before(time.Now().UTC()) {		
		return errors.New("target_date must be in the future")
	}

	return nil
}

func (d *PostData) Insert(db *db.DB, ctx context.Context) (*string, error) {
	query := `
		INSERT INTO countdown (target_date, message, fallback_message)
		VALUES ($1, $2, $3)
		RETURNING uuid
	`

	var uuid string
	err := db.QueryRow(ctx, query,
		d.TargetDate,
		d.Message,
		d.FallbackMessage,
	).Scan(&uuid)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &uuid, nil
}


