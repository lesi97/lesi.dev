package db

import (
	"context"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/utils"
)

type ValidateApiDetailsArgs struct {
	Application string
	Logger *utils.Logger
	DiscordUsername string
	AuthUrl string
}

func (db *DB) ValidateAndFetchApiDetails(data ValidateApiDetailsArgs) (*ApiDetails, error) {
	application := data.Application
	logger := data.Logger
	discordUsername := data.DiscordUsername
	authUrl := data.AuthUrl

	apiDetails, err := db.FetchApiDetails(context.Background(), application, logger)
	if err != nil || apiDetails == nil {
		message := fmt.Sprintf("FATAL: ERROR GETTING %v API DETAILS %v\n", application, err)
		logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: discordUsername,
			Title:    fmt.Sprintf("%v FATAL", application),
		})
		return nil, fmt.Errorf("%s", message)
	}

	noRefreshToken := apiDetails.RefreshToken == nil
	if noRefreshToken {
		message := fmt.Sprintf("No %v Refresh Token Found", application)
		logger.Error(message)
		return nil, fmt.Errorf("%s", message)
	}
	
	if apiDetails.RefreshTokenExpiry == nil {
		message := fmt.Sprintf("No %v Refresh Token Expiry Found", application)
		logger.Error(message)
		return nil, fmt.Errorf("%s", message)
	}

	expired := utils.IsRefreshTokenExpired(*apiDetails.RefreshTokenExpiry)
	if expired {
		newAuthData, err := utils.RefreshToken(application, *apiDetails.ClientID, *apiDetails.ClientSecret, *apiDetails.RefreshToken, authUrl)
		if err != nil {
			fmt.Printf("NEW AUTH DATA ERR: %v\n", err)
			return nil, err
		}
		apiDetails.AccessToken = &newAuthData.AccessToken
		apiDetails.RefreshToken = &newAuthData.RefreshToken
		apiDetails.RefreshTokenExpiry = newAuthData.RefreshTokenExpiry
		_ = db.UpdateApiDetails(
			context.Background(),
			application,
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
	return apiDetails, nil
}
