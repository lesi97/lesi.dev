package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type ApiDetails struct {
	Name               *string
	ClientID           *string
	ClientSecret       *string
	AccessToken        *string
	RefreshToken       *string
	RefreshTokenExpiry *int64
	BaseURL            *string
	RedirectURL        *string
}

var allowedApplications = map[string]struct{}{
	"Anilist":  {},
	"Nightbot": {},
	"Twitch_GO":   {},
	"Spotify":  {},
	"Nasa":     {},
}

func isAllowedApplication(app string) bool {
	_, ok := allowedApplications[app]
	return ok
}

func (db *DB) FetchApiDetails(ctx context.Context, application string, logger *utils.Logger) (*ApiDetails, error) {
	if !isAllowedApplication(application) {
		logger.Errorf("Invalid application: " + application)
		return nil, errors.New("invalid application")
	}

	var (
		name         *string
		clientID     *string
		clientSecret *string
		accessToken  *string
		refreshToken *string
		expiry       *int64
		baseURL      *string
		redirectURL  *string
	)

	err := db.QueryRow(ctx, `
		select
			name,
			client_id,
			client_secret,
			access_token,
			refresh_token,
			refresh_token_expiry,
			base_url,
			redirect_url
		from api_keys
		where name = $1
	`, application).Scan(
		&name,
		&clientID,
		&clientSecret,
		&accessToken,
		&refreshToken,
		&expiry,
		&baseURL,
		&redirectURL,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if clientID != nil {
		v, err := Decrypt(*clientID)
		if err != nil {
			return nil, err
		}
		clientID = &v
	}
	if clientSecret != nil {
		v, err := Decrypt(*clientSecret)
		if err != nil {
			return nil, err
		}
		clientSecret = &v
	}
	if accessToken != nil {
		v, err := Decrypt(*accessToken)
		if err != nil {
			return nil, err
		}
		accessToken = &v
	}
	if refreshToken != nil {
		v, err := Decrypt(*refreshToken)
		if err != nil {
			return nil, err
		}
		refreshToken = &v
	}


	return &ApiDetails{
		Name:               name,
		ClientID:           clientID,
		ClientSecret:       clientSecret,
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: expiry,
		BaseURL:            baseURL,
		RedirectURL:        redirectURL,
	}, nil
}