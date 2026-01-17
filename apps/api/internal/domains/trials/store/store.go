package trials_store

import (
	"os"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type TrialsStoreInterface interface {
	GetLoot() *string
	GetPlayerCount() *string
}

type TrialsStore struct {
	store.StoreBase
	url				string
	steamClientId 	string
	steamUrl		string
	steamClientIdAvailable bool
}

func NewStore(db *database.DB, logger *utils.Logger) *TrialsStore{
	steamApiKey := os.Getenv("STEAM_CLIENT_ID")
	if steamApiKey == "" {
		err := "FATAL: ERROR GETTING STEAM_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: err,
			Username: "STEAM STORE FATAL",
			Title: "STEAM STORE FATAL",
		})
	}
	return &TrialsStore{
		StoreBase: store.NewStoreBase(db, logger),
		url: "https://api.trialsofthenine.com/weeks/0",
		steamClientId: steamApiKey,
		steamUrl: "https://api.steampowered.com",
		steamClientIdAvailable: steamApiKey != "",
	}
}


