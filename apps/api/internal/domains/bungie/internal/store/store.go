package store

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Methods interface {
	GetCharacterPlayTime(ctx context.Context) (*string, error)
	GetEquippedWeapon(ctx context.Context) (*string, error)
	GetTerrorWeapon(ctx context.Context) (*string, error)
	GetYungerWeapon(ctx context.Context) (*string, error)
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
	bungieApiKey := os.Getenv("BUNGIE_CLIENT_ID")
	if bungieApiKey == "" {
		message := "FATAL: ERROR GETTING BUNGIE_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: "BUNGIE STORE FATAL",
			Title:    "BUNGIE STORE FATAL",
		})
		return nil, fmt.Errorf("%s", message)
	}

	return &Store{
		DB:         db,
		Logger:     logger,
		Redis:      redis,
		HTTPClient: httpClient,
		BaseURL:    "https://www.bungie.net",
		ClientID:   bungieApiKey,
	}, nil
}
