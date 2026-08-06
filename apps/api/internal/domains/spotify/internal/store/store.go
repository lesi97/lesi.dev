package store

import (
	"context"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Methods interface {
	InsertScrobble(ctx context.Context, input model.ScrobbleInput) (*model.ScrobbleResult, error)
	GetLatestPlayedText(ctx context.Context) (*string, error)
	PollSpotifyRecentlyPlayed(ctx context.Context, input model.SpotifyPollInput) (*model.SpotifyPollResult, error)
}

type Store struct {
	DB         *db.DB
	Logger     *utils.Logger
	Redis      *redis.Client
	HTTPClient *http.Client
}

func NewStore(db *db.DB, logger *utils.Logger, redis *redis.Client) *Store {
	return &Store{
		DB:         db,
		Logger:     logger,
		Redis:      redis,
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	}
}
