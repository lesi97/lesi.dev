package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TwitchGET(
	database *db.DB,
	logger *utils.Logger,
	apiDetails *db.ApiDetails,
	clientID string,
	clientSecret string,
	authURL string,
	url string,
) ([]byte, error) {
	defer logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", url), time.Now(), nil)

	err := EnsureValidApiDetails(database, logger, "Twitch_GO", clientID, clientSecret, authURL, apiDetails)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", *apiDetails.AccessToken))
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
