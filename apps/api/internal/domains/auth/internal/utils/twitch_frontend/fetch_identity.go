package twitch_frontend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/httpapi"
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

	bodyBytes, statusCode, err := httpapi.DoRequest(
		ctx,
		httpClient,
		http.MethodGet,
		"https://api.twitch.tv/helix/users",
		nil,
		map[string]string{
			"Authorization": "Bearer " + accessToken,
			"Client-Id":     clientID,
		},
	)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode >= 300 {
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

	if err := json.Unmarshal(bodyBytes, &body); err != nil {
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
