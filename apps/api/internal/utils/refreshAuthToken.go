package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	requestmetrics "github.com/lesi97/lesi.dev/internal/request_metrics"
)

type refreshTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

type refreshTokenResponse struct {
	AccessToken         *string `json:"access_token"`
	RefreshToken        *string `json:"refresh_token"`
	RefreshTokenExpiry  *int64  `json:"refresh_token_expiry"`
	RefreshTokenExpiry2 *int64  `json:"expires_in"`
}

type RefreshTokenResult struct {
	AccessToken        string
	RefreshToken       string
	RefreshTokenExpiry *int64
}

func RefreshToken(
	ctx context.Context,
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

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		authUrl,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		requestmetrics.AddFetchCallsDuration(ctx, time.Since(start), err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		requestmetrics.AddFetchCallsDuration(ctx, time.Since(start), nil)
		errMessage := fmt.Sprintf("failed to refresh %v token | status code: %v", application, resp.StatusCode)
		return nil, errors.New(errMessage)
	}

	var data refreshTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		requestmetrics.AddFetchCallsDuration(ctx, time.Since(start), err)
		return nil, err
	}

	if data.AccessToken == nil || data.RefreshToken == nil {
		errMessage := fmt.Sprintf("failed to refresh %v token, no access token provided", application)
		return nil, errors.New(errMessage)
	}

	if data.RefreshTokenExpiry == nil || *data.RefreshTokenExpiry == int64(0) {
		data.RefreshTokenExpiry = data.RefreshTokenExpiry2
	}

	requestmetrics.AddFetchCallsDuration(ctx, time.Since(start), nil)

	return &RefreshTokenResult{
		AccessToken:        *data.AccessToken,
		RefreshToken:       *data.RefreshToken,
		RefreshTokenExpiry: data.RefreshTokenExpiry,
	}, nil
}
