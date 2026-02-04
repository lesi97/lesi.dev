package store

import (
	"context"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Methods interface {
	GetPlayerCount(ctx context.Context) (*string, error)
}

type Store struct {
	DB         *db.DB
	Logger     *utils.Logger
	Redis      *redis.Client
	HTTPClient *http.Client
	BaseURL    string
	ClientID   string
}

func NewStore(db *db.DB, logger *utils.Logger, redis *redis.Client, httpClient *http.Client) (*Store, error) {
	return &Store{
		DB:         db,
		Logger:     logger,
		Redis:      redis,
		HTTPClient: httpClient,
		BaseURL:    "https://fortnite.gg",
	}, nil
}
