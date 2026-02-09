package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Methods interface {
	RandomViewer(ctx context.Context, streamer string) (*string, error)
}

type Store struct {
	DB           *db.DB
	Logger       *utils.Logger
	Redis        *redis.Client
	ClientID     string
	ClientSecret string
	ApiDetails   *db.ApiDetails
	BaseURL      string
	AuthURL      string
}
