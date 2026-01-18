package twitch_frontend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Identity struct {
	ID          string
	Login       string
	DisplayName string
	AvatarURL   string
}

func FetchIdentity(
	ctx context.Context,
	accessToken string,
	clientID string,
) (*Identity, error) {
	httpClient := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.twitch.tv/helix/users",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", clientID)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("twitch users request failed")
	}

	var body struct {
		Data []struct {
			ID              string `json:"id"`
			Login           string `json:"login"`
			DisplayName     string `json:"display_name"`
			ProfileImageURL string `json:"profile_image_url"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	if len(body.Data) != 1 {
		return nil, fmt.Errorf("unexpected twitch user response")
	}

	u := body.Data[0]

	return &Identity{
		ID:          u.ID,
		Login:       u.Login,
		DisplayName: u.DisplayName,
		AvatarURL:   u.ProfileImageURL,
	}, nil
}
