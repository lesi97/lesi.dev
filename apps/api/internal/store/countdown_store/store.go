package countdown_store

import (
	"context"
	"log"
	"time"

	"github.com/lesi97/lesi.dev/internal/database"
)

type CountdownStore interface {
	GetCountdownByID(ctx context.Context, id string) (*string, error)
	InsertCountdown(ctx context.Context, data CountdownPostRequest) (*string, error)
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

