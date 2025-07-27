package store

import (
	"database/sql"
	"time"
)

type CountdownStore interface {
	GetCountdownByID(id string) (*CountdownData, error)
}

type SupabaseCountdownStore struct {
	db *sql.DB
}

type CountdownData struct {
	UUID 			string 		`json:"uuid"`
	TargetDate 		time.Time 	`json:"target_date"`
	Message 		string 		`json:"message"`
	FallbackMessage string 		`json:"fallback_message"`
}

func NewSupabaseCountdownStore(db *sql.DB) *SupabaseCountdownStore {
	return &SupabaseCountdownStore{db: db}
}

func (supabase *SupabaseCountdownStore) GetCountdownByID(uuid string) (*CountdownData, error) {
	countdown := &CountdownData{
		UUID: uuid,
	}

	query := `
		SELECT
			target_date,
			message,
			fallback_message
		FROM countdown
		WHERE uuid = $1
	`
	err := supabase.db.QueryRow(query, uuid).Scan(
		&countdown.TargetDate,
		&countdown.Message,
		&countdown.FallbackMessage,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return countdown, nil
}