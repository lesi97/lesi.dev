package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type oauthTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Code         string `json:"code"`
}

type OAuthTokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	Scope        []string `json:"scope"`
}

func ExchangeCodeForToken(
	ctx context.Context,
	httpClient *http.Client,
	tokenURL string,
	clientID string,
	clientSecret string,
	redirectURI string,
	code string,
) (OAuthTokenResponse, error) {
	reqBody := oauthTokenRequest{
		GrantType:    "authorization_code",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Code:         code,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return OAuthTokenResponse{}, fmt.Errorf("marshal token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewReader(body))
	if err != nil {
		return OAuthTokenResponse{}, fmt.Errorf("create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return OAuthTokenResponse{}, fmt.Errorf("do token request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return OAuthTokenResponse{}, fmt.Errorf("read token response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return OAuthTokenResponse{}, fmt.Errorf("token request failed: status=%d body=%s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResponse OAuthTokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
		return OAuthTokenResponse{}, fmt.Errorf("unmarshal token response: %w", err)
	}

	if tokenResponse.AccessToken == "" {
		return OAuthTokenResponse{}, fmt.Errorf("token response missing access_token")
	}

	return tokenResponse, nil
}
