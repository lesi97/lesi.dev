package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/httpapi"
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

	headers := map[string]string{
		"Client-ID":     clientID,
		"Authorization": fmt.Sprintf("Bearer %v", *apiDetails.AccessToken),
	}
	body, statusCode, err := httpapi.DoRequest(nil, &http.Client{}, http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d: %s", statusCode, string(body))
	}

	return body, nil
}
