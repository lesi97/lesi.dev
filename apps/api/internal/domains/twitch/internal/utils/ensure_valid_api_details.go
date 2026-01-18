package utils

import (
	"context"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func EnsureValidApiDetails(
	database *db.DB,
	logger *utils.Logger,
	application string,
	clientID string,
	clientSecret string,
	authURL string,
	apiDetails *db.ApiDetails,
) error {
	if apiDetails.RefreshTokenExpiry == nil || apiDetails.RefreshToken == nil {
		return fmt.Errorf("missing refresh token details")
	}

	expired := utils.IsRefreshTokenExpired(*apiDetails.RefreshTokenExpiry)
	if !expired {
		return nil
	}

	api, err := utils.RefreshToken(application, clientID, clientSecret, *apiDetails.RefreshToken, authURL)
	if err != nil {
		return fmt.Errorf("failed to refresh twitch api details")
	}

	apiDetails.AccessToken = &api.AccessToken
	apiDetails.RefreshToken = &api.RefreshToken
	apiDetails.RefreshTokenExpiry = api.RefreshTokenExpiry

	return database.UpdateApiDetails(
		context.Background(),
		application,
		&clientID,
		&clientSecret,
		&api.AccessToken,
		&api.RefreshToken,
		api.RefreshTokenExpiry,
		apiDetails.BaseURL,
		apiDetails.RedirectURL,
		logger,
	)
}
