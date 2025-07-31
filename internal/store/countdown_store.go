package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lesi97/api.lesi.dev/internal/database"
	"github.com/lesi97/api.lesi.dev/internal/utils"
)

type CountdownStore interface {
	GetCountdownByID(ctx context.Context, id string) (*string, error)
}

type SupabaseCountdownStore struct {
	db *database.Supabase
	logger *log.Logger
}

type CountdownData struct {
	UUID 			string 		`json:"uuid"`
	TargetDate 		time.Time 	`json:"target_date"`
	Message 		string 		`json:"message"`
	FallbackMessage string 		`json:"fallback_message"`
}

func NewSupabaseCountdownStore(db *database.Supabase) *SupabaseCountdownStore {
	return &SupabaseCountdownStore{db: db}
}

func (supabase *SupabaseCountdownStore) GetCountdownByID(ctx context.Context, uuid string) (*string, error) {
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
	err := supabase.db.QueryRow(ctx, query, uuid).Scan(
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

	result := countdownToString(countdown.TargetDate)
	if result == "Passed" {
		message := countdown.FallbackMessage
		utils.TruncateString(&message, 400)
		return &message, nil
	}

	message := fmt.Sprintf("%s %s", result, countdown.Message)
	utils.TruncateString(&message, 400)

	return &message, nil
}

/*
Utility function to convert check if the date from the countdown api route has been reached or not

Returns string "Passed" if the date has passed

Returns string difference in time if the countdown has not yet been reached
*/
func countdownToString(target time.Time) string {
	now := time.Now()
	if now.After(target) {
		return "Passed"
	}

	diff := target.Sub(now).Round(time.Second)

	days := int(diff.Hours()) / 24
	hours := int(diff.Hours()) % 24
	minutes := int(diff.Minutes()) % 60
	seconds := int(diff.Seconds()) % 60

	return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, seconds)
}
