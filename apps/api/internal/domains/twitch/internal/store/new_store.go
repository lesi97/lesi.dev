package store

import (
	"fmt"
	"os"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func NewStore(database *db.DB, logger *utils.Logger, redis *redis.Client) (*Store, error) {
	const userName = "Twitch_GO"
	twitchAuthURL := "https://id.twitch.tv/oauth2/token"
	baseURL := "https://api.twitch.tv/helix"

	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientID == "" {
		message := "FATAL: ERROR GETTING TWITCH_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: "TWITCH STORE FATAL",
			Title:    "TWITCH STORE FATAL",
		})
		return nil, fmt.Errorf("%s", message)
	}

	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if twitchClientSecret == "" {
		message := "FATAL: ERROR GETTING TWITCH_CLIENT_SECRET ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: "TWITCH STORE FATAL",
			Title:    "TWITCH STORE FATAL",
		})
		return nil, fmt.Errorf("%s", message)
	}

	apiDetails, err := database.ValidateAndFetchApiDetails(db.ValidateApiDetailsArgs{
		Application:     "Twitch_GO",
		Logger:          logger,
		DiscordUsername: userName,
		AuthUrl:         twitchAuthURL,
	})
	if err != nil {
		return nil, err
	}

	return &Store{
		DB:           database,
		Logger:       logger,
		Redis:        redis,
		ClientID:     twitchClientID,
		ClientSecret: twitchClientSecret,
		BaseURL:      baseURL,
		AuthURL:      twitchAuthURL,
		ApiDetails:   apiDetails,
	}, nil
}
