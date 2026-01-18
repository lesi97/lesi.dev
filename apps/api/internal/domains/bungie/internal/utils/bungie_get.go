package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func BungieGET(logger *utils.Logger, clientID string, url string) ([]byte, error) {
	defer logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", url), time.Now(), nil)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-API-Key", clientID)
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

	isErr, err := isBungieError(body)
	if isErr {
		return nil, err
	}

	return body, nil
}
