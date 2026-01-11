package twitch_store

import (
	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type TwichStoreInterface interface {
	RandomViewer(streamer string) (*string, error)
}

type TwitchStore struct {
	store.StoreBase
	api_details 	*database.ApiDetails
	base_url		string
	auth_url		string
}


func NewStore(db *database.DB, logger *utils.Logger) (*TwitchStore, error) {
	const userName = "Twitch_GO"
	const twitchAuthUrl = "https://id.twitch.tv/oauth2/token"

	apiDetails, err := db.ValidateAndFetchApiDetails(database.ValidateApiDetailsArgs{
		Application: "Twitch_GO",
		Logger: logger,
		DiscordUsername: userName,
		AuthUrl: twitchAuthUrl,
	})
	if err != nil {
		return nil, err
	}

	return &TwitchStore{
		StoreBase: store.NewStoreBase(db, logger),
		base_url: "https://api.twitch.tv/helix",
		auth_url: twitchAuthUrl,
		api_details: apiDetails,
	}, nil
}
