package store

import (
	"os"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func NewStore(db *db.DB, logger *utils.Logger) *Store {
	steamApiKey := os.Getenv("STEAM_CLIENT_ID")
	if steamApiKey == "" {
		err := "FATAL: ERROR GETTING STEAM_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  err,
			Username: "STEAM STORE FATAL",
			Title:    "STEAM STORE FATAL",
		})
	}
	return &Store{
		DB:                     db,
		Logger:                 logger,
		URL:                    "https://api.trialsofthenine.com/weeks/0",
		SteamClientID:          steamApiKey,
		SteamURL:               "https://api.steampowered.com",
		SteamClientIDAvailable: steamApiKey != "",
	}
}
