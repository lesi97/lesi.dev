package anilist_store

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type refreshTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

type refreshTokenResponse struct {
	AccessToken        *string `json:"access_token"`
	RefreshToken       *string `json:"refresh_token"`
	RefreshTokenExpiry *int64  `json:"refresh_token_expiry"`
}

type refreshTokenResult struct {
	AccessToken        string
	RefreshToken       string
	RefreshTokenExpiry *int64
}

func refreshToken(
	clientID string,
	clientSecret string,
	refreshTokenValue string,
) (*refreshTokenResult, error) {
	payload := refreshTokenRequest{
		GrantType:    "refresh_token",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshTokenValue,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://anilist.co/api/v2/oauth/token",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("failed to refresh AniList token")
	}

	var data refreshTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if data.AccessToken == nil || data.RefreshToken == nil {
		return nil, errors.New("failed to refresh AniList token")
	}

	return &refreshTokenResult{
		AccessToken:        *data.AccessToken,
		RefreshToken:       *data.RefreshToken,
		RefreshTokenExpiry: data.RefreshTokenExpiry,
	}, nil
}
