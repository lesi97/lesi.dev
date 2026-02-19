package store

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Methods interface {
	GetPlayerCount(ctx context.Context, gameID string, gameName string) (*string, error)
}

type Store struct {
	DB         *db.DB
	Logger     *utils.Logger
	Redis      *redis.Client
	HTTPClient *http.Client
	BaseURL    string
	StoreURL   string
	SteamKey   string
}

func NewStore(db *db.DB, logger *utils.Logger, redis *redis.Client, httpClient *http.Client) (*Store, error) {
	steamAPIKey := os.Getenv("STEAM_CLIENT_ID")
	if steamAPIKey == "" {
		err := "FATAL: ERROR GETTING STEAM_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  err,
			Username: "STEAM STORE FATAL",
			Title:    "STEAM STORE FATAL",
		})
		return nil, errors.New(err)
	}

	return &Store{
		DB:         db,
		Logger:     logger,
		Redis:      redis,
		HTTPClient: httpClient,
		BaseURL:    "https://api.steampowered.com",
		StoreURL:   "https://store.steampowered.com",
		SteamKey:   steamAPIKey,
	}, nil
}
