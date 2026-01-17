package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
)

type FetchData struct {
	UUID            *string    `json:"uuid"`
	TargetDate      time.Time `json:"target_date"`
	Message         string    `json:"message"`
	FallbackMessage string    `json:"fallback_message"`
}

func (d *FetchData) Select(db *db.DB, ctx *context.Context, uuid *string) error {
	query := `
		SELECT
			target_date,
			message,
			fallback_message
		FROM countdown
		WHERE uuid = $1
	`
	err := db.QueryRow(*ctx, query, *uuid).Scan(
		&d.TargetDate,
		&d.Message,
		&d.FallbackMessage,
	)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}