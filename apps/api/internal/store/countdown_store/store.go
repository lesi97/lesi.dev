package countdown_store

import (
	"context"
	"time"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type CountdownStoreInterface interface {
	GetCountdownByID(ctx context.Context, id string) (*string, error)
	InsertCountdown(ctx context.Context, data CountdownPostRequest) (*string, error)
}

type CountdownStore struct {
	store.StoreBase
}

type CountdownData struct {
	UUID 			string 		`json:"uuid"`
	TargetDate 		time.Time 	`json:"target_date"`
	Message 		string 		`json:"message"`
	FallbackMessage string 		`json:"fallback_message"`
}

func NewStore(db *database.DB, logger *utils.Logger) *CountdownStore {
	return &CountdownStore{
		StoreBase: store.NewStoreBase(db, logger),
	}
}

