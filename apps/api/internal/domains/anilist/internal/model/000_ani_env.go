package model

import (
	"errors"
	"os"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type AnilistEnv struct {
	DiscordUsername string
	BaseUrl         string
	GraphqlUrl      string
	AuthUrl         string
	ClientId        string
	ClientSecret    string
	AccessToken     string
	RefreshToken    string
	RefreshExpires  int64
}

func (e *AnilistEnv) Validate(database *db.DB, logger *utils.Logger) error {
	e.DiscordUsername = "Anilist"
	e.BaseUrl = "https://anilist.co"
	e.GraphqlUrl = "https://graphql.anilist.co"
	e.AuthUrl = "https://anilist.co/api/v2/oauth/token"

	clientId := os.Getenv("ANILIST_CLIENT_ID")
	if clientId == "" {
		return errors.New("ANILIST_CLIENT_ID not found in env")
	}

	clientSecret := os.Getenv("ANILIST_CLIENT_SECRET")
	if clientSecret == "" {
		return errors.New("ANILIST_CLIENT_SECRET not found in env")
	}

	apiDetails, err := database.ValidateAndFetchApiDetails(db.ValidateApiDetailsArgs{
		Application:     "Anilist",
		Logger:          logger,
		DiscordUsername: e.DiscordUsername,
		AuthUrl:         e.AuthUrl,
	})
	if err != nil {
		return err
	}

	e.ClientId = clientId
	e.ClientSecret = clientSecret

	e.AccessToken = *apiDetails.AccessToken
	e.RefreshToken = *apiDetails.RefreshToken
	e.RefreshExpires = *apiDetails.RefreshTokenExpiry

	return nil

}
