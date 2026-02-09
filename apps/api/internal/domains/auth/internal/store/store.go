package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Methods interface {
	AnilistAuthUrl() (*string, error)
	AnilistCallback(ctx context.Context, code string) error
	TwitchModAuthUrl() (*string, error)
	TwitchModCallback(ctx context.Context, code string) error

	TwitchFrontendAuthStart() (*TwitchFrontendAuthStartResult, error)
	TwitchFrontendCallback(ctx context.Context, code string, state string, expectedState string, pkceVerifier string) (*TwitchFrontendIdentity, error)
	TwitchFrontendUpsertUser(ctx context.Context, identity TwitchFrontendIdentity) (string, error)
	TwitchFrontendCreateSession(ctx context.Context, userID string, ttl time.Duration) (string, error)
	TwitchFrontendGetUserBySession(ctx context.Context, sessionToken string) (*TwitchFrontendUser, error)
	TwitchFrontendDeleteSessionByToken(ctx context.Context, sessionToken string) error
}

type ApiApplication struct {
	client_id     *string
	client_secret *string
	base_url      string
}

type Store struct {
	DB      *db.DB
	Logger  *utils.Logger
	anilist ApiApplication
	twitch  ApiApplication
	apiUrl  *string
	webUrl  *string
}

func NewStore(db *db.DB, logger *utils.Logger) (*Store, error) {
	const discord = "AUTH STORE FATAL"

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		message := "FATAL: ERROR GETTING API_URL ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: discord,
			Title:    discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	webUrl := os.Getenv("WEB_URL")
	if webUrl == "" {
		message := "FATAL: ERROR GETTING WEB_URL ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: discord,
			Title:    discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	anilistClientId := os.Getenv("ANILIST_CLIENT_ID")
	if anilistClientId == "" {
		message := "FATAL: ERROR GETTING ANILIST_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: discord,
			Title:    discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	anilistClientSecret := os.Getenv("ANILIST_CLIENT_SECRET")
	if anilistClientSecret == "" {
		message := "FATAL: ERROR GETTING ANILIST_CLIENT_SECRET ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: discord,
			Title:    discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	twitchClientId := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientId == "" {
		message := "FATAL: ERROR GETTING TWITCH_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: discord,
			Title:    discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if twitchClientSecret == "" {
		message := "FATAL: ERROR GETTING TWITCH_CLIENT_SECRET ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: discord,
			Title:    discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	anilist := ApiApplication{
		client_id:     &anilistClientId,
		client_secret: &anilistClientSecret,
		base_url:      "https://anilist.co/api/v2/oauth",
	}

	twitch := ApiApplication{
		client_id:     &twitchClientId,
		client_secret: &twitchClientSecret,
		base_url:      "https://id.twitch.tv/oauth2",
	}

	return &Store{
		DB:      db,
		Logger:  logger,
		anilist: anilist,
		twitch:  twitch,
		apiUrl:  &apiUrl,
		webUrl:  &webUrl,
	}, nil
}
