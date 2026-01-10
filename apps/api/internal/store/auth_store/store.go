package auth_store

import (
	"context"
	"time"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type AuthStoreInterface interface {
	GenerateAnilistAuthUrl() (*string, error)
	AnilistCallback(code string) error
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
	twitchApiDetails, _ := db.FetchApiDetails(context.Background(), "Twitch", logger)

	anilist := ApiApplication{
		api_details: anilistApiDetails,
		base_url: "https://anilist.co/api/v2/oauth",
	}

	twitch := ApiApplication{
		api_details: twitchApiDetails,
		base_url: "TODO",
	}

	return &AuthStore{
		StoreBase: store.NewStoreBase(db, logger),
		anilist: anilist,
		twitch: twitch,
	}
}
