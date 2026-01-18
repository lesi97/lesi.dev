package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Methods interface {
	AnilistAuthUrl() (*string, error)
	AnilistCallback(code string) error
	TwitchModAuthUrl() (*string, error)
	TwitchModCallback(code string) error
	
	TwitchFrontendAuthStart() (*TwitchFrontendAuthStartResult, error)
	TwitchFrontendCallback(code string, state string, expectedState string, pkceVerifier string) (*TwitchFrontendIdentity, error)
	TwitchFrontendUpsertUser(ctx context.Context, identity TwitchFrontendIdentity) (string, error)
	TwitchFrontendCreateSession(ctx context.Context, userID string, ttl time.Duration) (string, error)
	TwitchFrontendGetUserBySession(ctx context.Context, sessionToken string) (*TwitchFrontendUser, error)
	TwitchFrontendDeleteSessionByToken(ctx context.Context, sessionToken string) error
}

type ApiApplication struct {
	client_id *string
	client_secret *string
	base_url string
}

type Store struct {
	store.StoreBase
	anilist ApiApplication
	twitch ApiApplication
	apiUrl *string
	webUrl *string
}

type OAuthClientConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type OAuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    *time.Time
	TokenType    string
}


func NewStore(db *database.DB, logger *utils.Logger) (*Store, error) {
	const userName = "Auth Store"
	const discord = "AUTH STORE FATAL"

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		message := "FATAL: ERROR GETTING API_URL ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: discord,
			Title: discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	webUrl := os.Getenv("WEB_URL")
	if webUrl == "" {
		message := "FATAL: ERROR GETTING WEB_URL ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: discord,
			Title: discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	anilistClientId := os.Getenv("ANILIST_CLIENT_ID")
	if anilistClientId == "" {
		message := "FATAL: ERROR GETTING ANILIST_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: discord,
			Title: discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	anilistClientSecret := os.Getenv("ANILIST_CLIENT_SECRET")
	if anilistClientSecret == "" {
		message := "FATAL: ERROR GETTING ANILIST_CLIENT_SECRET ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: discord,
			Title: discord,
		})
		return nil, fmt.Errorf("%s", message)
	}


	twitchClientId := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientId == "" {
		message := "FATAL: ERROR GETTING TWITCH_CLIENT_ID ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: discord,
			Title: discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if twitchClientSecret == "" {
		message := "FATAL: ERROR GETTING TWITCH_CLIENT_SECRET ENV VAR"
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: discord,
			Title: discord,
		})
		return nil, fmt.Errorf("%s", message)
	}

	anilist := ApiApplication{
		client_id: &anilistClientId,
		client_secret: &anilistClientSecret,
		base_url: "https://anilist.co/api/v2/oauth",
	}

	twitch := ApiApplication{
		client_id: &twitchClientId,
		client_secret: &twitchClientSecret,
		base_url: "https://id.twitch.tv/oauth2",
	}

	return &Store{
		StoreBase: store.NewStoreBase(db, logger),
		anilist: anilist,
		twitch: twitch,
		apiUrl: &apiUrl,
		webUrl: &webUrl,
	}, nil
}
