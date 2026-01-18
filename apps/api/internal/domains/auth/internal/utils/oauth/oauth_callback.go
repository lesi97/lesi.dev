package oauth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func OAuthCallback(
	ctx context.Context,
	database *db.DB,
	logger *utils.Logger,
	code string,
	cfg ProviderConfig,
) error {
	if cfg.Application == "" {
		return fmt.Errorf("missing application")
	}
	if cfg.TokenURL == "" {
		return fmt.Errorf("missing token url")
	}
	if cfg.ClientID == nil || *cfg.ClientID == "" {
		return fmt.Errorf("missing client id")
	}
	if cfg.ClientSecret == nil || *cfg.ClientSecret == "" {
		return fmt.Errorf("missing client secret")
	}
	if cfg.RedirectURL == "" {
		return fmt.Errorf("missing redirect url")
	}

	httpClient := &http.Client{Timeout: 15 * time.Second}

	tokenResp, err := ExchangeCodeForToken(
		ctx,
		httpClient,
		cfg.TokenURL,
		*cfg.ClientID,
		*cfg.ClientSecret,
		cfg.RedirectURL,
		code,
	)
	if err != nil {
		return err
	}

	accessToken := tokenResp.AccessToken
	refreshToken := tokenResp.RefreshToken

	var refreshTokenExpiry *int64
	if tokenResp.ExpiresIn > 0 {
		expiryUnix := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Unix()
		refreshTokenExpiry = &expiryUnix
	}

	clientIDPtr := cfg.ClientID
	clientSecretPtr := cfg.ClientSecret
	redirectPtr := cfg.RedirectURL

	accessTokenPtr := &accessToken
	var refreshTokenPtr *string
	if refreshToken != "" {
		refreshTokenPtr = &refreshToken
	}

	_ = database.UpdateApiDetails(
		ctx,
		cfg.Application,
		clientIDPtr,
		clientSecretPtr,
		accessTokenPtr,
		refreshTokenPtr,
		refreshTokenExpiry,
		nil,
		&redirectPtr,
		logger,
	)

	return nil
}
