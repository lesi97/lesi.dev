package store

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/httpapi"
)

type spotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Store) SpotifyCallback(ctx context.Context, code string) error {
	creds, err := s.spotifyCredentials(ctx)
	if err != nil {
		return err
	}

	tokenResp, err := exchangeSpotifyCode(ctx, creds, code)
	if err != nil {
		return err
	}

	if tokenResp.RefreshToken == "" {
		return fmt.Errorf("spotify token response missing refresh_token")
	}

	accessToken := tokenResp.AccessToken
	refreshToken := tokenResp.RefreshToken
	refreshTokenExpiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UnixMilli()
	baseURL := creds.APIURL
	redirectURL := creds.RedirectURL

	return s.DB.UpdateApiDetails(
		ctx,
		"Spotify",
		&creds.ClientID,
		&creds.ClientSecret,
		&accessToken,
		&refreshToken,
		&refreshTokenExpiry,
		&baseURL,
		&redirectURL,
		s.Logger,
	)
}

func exchangeSpotifyCode(ctx context.Context, creds spotifyCredentials, code string) (spotifyTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", creds.RedirectURL)

	tokenURL := fmt.Sprintf("%v/api/token", creds.AccountsURL)
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(creds.ClientID+":"+creds.ClientSecret))
	httpClient := &http.Client{Timeout: 15 * time.Second}

	bodyBytes, statusCode, err := httpapi.DoRequest(
		ctx,
		httpClient,
		http.MethodPost,
		tokenURL,
		strings.NewReader(form.Encode()),
		map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/x-www-form-urlencoded",
			"Accept":        "application/json",
		},
	)
	if err != nil {
		return spotifyTokenResponse{}, err
	}

	if statusCode < 200 || statusCode >= 300 {
		return spotifyTokenResponse{}, fmt.Errorf("spotify token exchange failed: status=%d body=%s", statusCode, string(bodyBytes))
	}

	var tokenResp spotifyTokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return spotifyTokenResponse{}, err
	}
	if tokenResp.AccessToken == "" {
		return spotifyTokenResponse{}, fmt.Errorf("spotify token response missing access_token")
	}

	return tokenResp, nil
}
