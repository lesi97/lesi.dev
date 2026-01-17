package auth_store

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (s *AuthStore) oauthCallback(
	ctx context.Context,
	code string,
	cfg oauthProviderConfig,
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

	tokenResp, err := exchangeCodeForToken(
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

	s.DB.UpdateApiDetails(
		ctx,
		cfg.Application,
		clientIDPtr,
		clientSecretPtr,
		accessTokenPtr,
		refreshTokenPtr,
		refreshTokenExpiry,
		nil,
		&redirectPtr,
		s.Logger,
	)

	return nil
}