package anilist_store

import (
	"context"
	"fmt"
	"os"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type AnilistStoreInterface interface {
	UpdateAnilist(ctx context.Context, data PlexWebhookPayload) error
}

type AnilistStore struct {
	store.StoreBase
	graphql_url 	string
	auth_url		string
	api_details 	*database.ApiDetails
	discord_username string
	plex			PlexInfo	
	seasonCountCache map[string]map[int]int

}

type PlexInfo struct {
	baseUrl		string
	username 	string
	xtoken 		string
}

func NewStore(db *database.DB, logger *utils.Logger) (*AnilistStore, error) {
	const userName = "Anilist Updater"
	apiDetails, err := db.FetchApiDetails(context.Background(),"Anilist", logger)
	if err != nil || apiDetails == nil {
		message := fmt.Sprintf("FATAL: ERROR GETTING ANILIST API DETAILS %v\n", err)
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: userName,
			Title: "ANILIST FATAL",
		})
		return nil, fmt.Errorf("FATAL: ERROR GETTING ANILIST API DETAILS %v\n", err)
	}

	noRefreshToken := apiDetails.RefreshToken == nil
	if noRefreshToken {
		logger.Error("No Anilist Refresh Token Found")
		return nil, fmt.Errorf("No Anilist Refresh Token Found")
	}

	expired := utils.IsRefreshTokenExpired(*apiDetails.RefreshTokenExpiry)
	if expired {
		newAuthData, err := refreshToken(*apiDetails.ClientID, *apiDetails.ClientSecret, *apiDetails.RefreshToken)
		if err != nil {
			return nil, err
		}
		apiDetails.AccessToken = &newAuthData.AccessToken
		apiDetails.RefreshToken = &newAuthData.RefreshToken
		apiDetails.RefreshTokenExpiry = newAuthData.RefreshTokenExpiry
		_ = db.UpdateApiDetails(
			context.Background(),
			"Anilist",
			nil,
			nil,
			apiDetails.AccessToken,
			apiDetails.RefreshToken,
			apiDetails.RefreshTokenExpiry,
			nil,
			nil,
			logger,
		)
	}

	plexInfo := PlexInfo{
		baseUrl: os.Getenv("PLEX_API_URL"),
		username: os.Getenv("PLEX_USERNAME"),
		xtoken: os.Getenv("PLEX_X_TOKEN"),
	}

	return &AnilistStore{
		StoreBase: store.NewStoreBase(db, logger),
		graphql_url: "https://graphql.anilist.co",
		auth_url: "https://anilist.co/api/v2/oauth/token",
		api_details: apiDetails,
		discord_username: userName,
		plex: plexInfo,
		seasonCountCache: make(map[string]map[int]int),
	}, nil
}
