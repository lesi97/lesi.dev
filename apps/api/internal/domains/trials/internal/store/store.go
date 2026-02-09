package store

import (
	"errors"
	"os"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Methods interface {
	GetLoot() *string
	GetPlayerCount() *string
}

type Store struct {
	DB                     *db.DB
	Logger                 *utils.Logger
	Redis                  *redis.Client
	URL                    string
	SteamClientID          string
	SteamURL               string
}


func NewStore(db *db.DB, logger *utils.Logger, redis *redis.Client) (*Store, error) {
	steamApiKey := os.Getenv("STEAM_CLIENT_ID")
	if steamApiKey == "" {
		err := "FATAL: ERROR GETTING STEAM_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  err,
			Username: "STEAM STORE FATAL",
			Title:    "STEAM STORE FATAL",
		})
		return nil, errors.New(err)
	}
	return &Store{
		DB:            db,
		Logger:        logger,
		Redis:         redis,
		URL:           "https://api.trialsofthenine.com/weeks/0",
		SteamClientID: steamApiKey,
		SteamURL:      "https://api.steampowered.com",
	}, nil
}
