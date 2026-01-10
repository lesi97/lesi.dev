package auth_store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type anilistTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Code         string `json:"code"`
}

type anilistTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func fetchAniListToken(ctx context.Context, httpClient *http.Client, clientID string, clientSecret string, redirectURI string, code string) (anilistTokenResponse, error) {
	reqBody := anilistTokenRequest{
		GrantType:    "authorization_code",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Code:         code,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return anilistTokenResponse{}, fmt.Errorf("marshal token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://anilist.co/api/v2/oauth/token", bytes.NewReader(b))
	if err != nil {
		return anilistTokenResponse{}, fmt.Errorf("create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return anilistTokenResponse{}, fmt.Errorf("do token request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return anilistTokenResponse{}, fmt.Errorf("read token response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return anilistTokenResponse{}, fmt.Errorf("anilist token request failed: status=%d body=%s", resp.StatusCode, string(bodyBytes))
	}

	var tr anilistTokenResponse
	if err := json.Unmarshal(bodyBytes, &tr); err != nil {
		return anilistTokenResponse{}, fmt.Errorf("unmarshal token response: %w", err)
	}

	if tr.AccessToken == "" {
		return anilistTokenResponse{}, fmt.Errorf("token response missing access_token")
	}

	return tr, nil
}

func (s *AuthStore) AnilistCallback(code string) error {
	clientID := s.anilist.api_details.ClientID
	secret := s.anilist.api_details.ClientSecret
	redirect := s.anilist.api_details.RedirectURL

	httpClient := &http.Client{Timeout: 15 * time.Second}

	tokenResp, err := fetchAniListToken(context.Background(), httpClient, *clientID, *secret, *redirect, code)
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

	application := "Anilist"
	var clientIDPtr *string = clientID
	var clientSecretPtr *string = secret
	var accessTokenPtr *string = &accessToken
	var refreshTokenPtr *string = &refreshToken
	var redirectPtr *string = redirect

	s.DB.UpdateApiDetails(
		context.Background(),
		application,
		clientIDPtr,
		clientSecretPtr,
		accessTokenPtr,
		refreshTokenPtr,
		refreshTokenExpiry,
		nil,
		redirectPtr,
		s.Logger,
	)

	return nil
}

