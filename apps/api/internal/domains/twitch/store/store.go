package twitch_store

import (
	"fmt"
	"os"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type TwichStoreInterface interface {
	RandomViewer(streamer string) (*string, error)
}

type TwitchStore struct {
	store.StoreBase
	redis 			*redis.Client
	client_id		*string
	client_secret 	*string
	api_details 	*database.ApiDetails
	base_url		*string
	auth_url		*string
}


func NewStore(db *database.DB, logger *utils.Logger, redis *redis.Client) (*TwitchStore, error) {
	const userName = "Twitch_GO"
	twitchAuthUrl := "https://id.twitch.tv/oauth2/token"
	baseUrl := "https://api.twitch.tv/helix"

	twitchClientId := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientId == "" {
		message := "FATAL: ERROR GETTING TWITCH_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: "TWITCH STORE FATAL",
			Title: "TWITCH STORE FATAL",
		})
		return nil, fmt.Errorf("%s", message)
	}

	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if twitchClientSecret == "" {
		message := "FATAL: ERROR GETTING TWITCH_CLIENT_SECRET ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: "TWITCH STORE FATAL",
			Title: "TWITCH STORE FATAL",
		})
		return nil, fmt.Errorf("%s", message)
	}

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
		redis: redis,
		client_id: &twitchClientId,
		client_secret: &twitchClientSecret,
		base_url: &baseUrl,
		auth_url: &twitchAuthUrl,
		api_details: apiDetails,
	}, nil
}
