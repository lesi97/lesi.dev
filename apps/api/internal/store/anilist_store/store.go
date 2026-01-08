package anilist_store

import (
	"context"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type AnilistStoreInterface interface {}

type AnilistStore struct {
	store.StoreBase
	graphql_url 	string
	auth_url		string
	api_details 	*database.ApiDetails
}

func NewStore(db *database.DB, logger *utils.Logger) *AnilistStore {
	apiDetails, err := db.FetchApiDetails(context.Background(),"Anilist", logger)
	if err != nil {
		message := fmt.Sprintf("FATAL: ERROR GETTING ANILIST API DETAILS %v\n", err)
		utils.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: "Anilist FATAL",
			Logger: logger,
		})
		logger.Fatal(message)
		return nil
	}

	noRefreshToken := apiDetails.RefreshToken == nil
	if noRefreshToken {
		logger.Error("No Anilist Refresh Token Found")
		return nil
	}

	expired := utils.IsRefreshTokenExpired(*apiDetails.RefreshTokenExpiry)
	if expired {
		newAuthData, err := refreshToken(*apiDetails.ClientID, *apiDetails.ClientSecret, *apiDetails.RefreshToken)
		if err != nil {
			logger.Errorf("Failed to refresh Anilist Auth Token: %v\n", err)
			return nil
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

	return &AnilistStore{
		StoreBase: store.NewStoreBase(db, logger),
		graphql_url: "https://graphql.anilist.co",
		auth_url: "https://anilist.co/api/v2/oauth/token",
		api_details: apiDetails,
	}
}


