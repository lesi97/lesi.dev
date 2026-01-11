package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	RefreshTokenExpiry2 *int64  `json:"expires_in"`
}

type RefreshTokenResult struct {
	AccessToken        string
	RefreshToken       string
	RefreshTokenExpiry *int64
}

func RefreshToken(
	application string,
	clientID string,
	clientSecret string,
	refreshTokenValue string,
	authUrl string,
) (*RefreshTokenResult, error) {
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
		authUrl,
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
		errMessage := fmt.Sprintf("failed to refresh %v token | status code: %v", application, resp.StatusCode)
		return nil, errors.New(errMessage)
	}

	var data refreshTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if data.AccessToken == nil || data.RefreshToken == nil {
		errMessage := fmt.Sprintf("failed to refresh %v token, no access token provided", application)
		return nil, errors.New(errMessage)
	}

	if data.RefreshTokenExpiry == nil || *data.RefreshTokenExpiry == int64(0) {
		data.RefreshTokenExpiry = data.RefreshTokenExpiry2
	}

	return &RefreshTokenResult{
		AccessToken:        *data.AccessToken,
		RefreshToken:       *data.RefreshToken,
		RefreshTokenExpiry: data.RefreshTokenExpiry,
	}, nil
}