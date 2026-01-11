package twitch_store

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
)


func (s *TwitchStore) twitchGET(url string) ([]byte, error) {
	defer s.Logger.LogExecutionTime(url, time.Now())

	expired := utils.IsRefreshTokenExpired(*s.api_details.RefreshTokenExpiry)
	if expired {
		api, err :=utils.RefreshToken("Twitch_GO", *s.api_details.ClientID, *s.api_details.ClientSecret, *s.api_details.RefreshToken, s.auth_url)
		if err != nil {
			return nil, fmt.Errorf("Failed to refresh Twitch API Details")
		}
		s.api_details.AccessToken = &api.AccessToken
		s.api_details.RefreshToken = &api.RefreshToken
		s.api_details.RefreshTokenExpiry = api.RefreshTokenExpiry
		s.DB.UpdateApiDetails(
			context.Background(),
			"Twitch_GO",
			s.api_details.ClientID,
			s.api_details.ClientSecret,
			&api.AccessToken,
			&api.RefreshToken,
			api.RefreshTokenExpiry,
			s.api_details.BaseURL,
			s.api_details.RedirectURL,
			s.Logger,
		)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Client-ID", *s.api_details.ClientID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", *s.api_details.AccessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}