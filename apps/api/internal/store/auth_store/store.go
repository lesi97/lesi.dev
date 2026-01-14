package auth_store

import (
	"context"
	"time"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type AuthStoreInterface interface {
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
	api_details *database.ApiDetails
	base_url string
}

type AuthStore struct {
	store.StoreBase
	anilist ApiApplication
	twitch ApiApplication
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


func NewStore(db *database.DB, logger *utils.Logger) *AuthStore {
	const userName = "Auth Store"

	anilistApiDetails, _ := db.FetchApiDetails(context.Background(), "Anilist", logger)
	twitchApiDetails, _ := db.FetchApiDetails(context.Background(), "Twitch_GO", logger)

	anilist := ApiApplication{
		api_details: anilistApiDetails,
		base_url: "https://anilist.co/api/v2/oauth",
	}

	twitch := ApiApplication{
		api_details: twitchApiDetails,
		base_url: "https://id.twitch.tv/oauth2",
	}

	return &AuthStore{
		StoreBase: store.NewStoreBase(db, logger),
		anilist: anilist,
		twitch: twitch,
	}
}
