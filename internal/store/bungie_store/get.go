package bungie_store

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type errorResponse struct {
	ErrorCode       int               `json:"ErrorCode"`
	ThrottleSeconds int               `json:"ThrottleSeconds"`
	ErrorStatus     string            `json:"ErrorStatus"`
	Message         string            `json:"Message"`
	MessageData     map[string]string `json:"MessageData"`
}


func bungieGET(url string) ([]byte, error) {
	apiKey := os.Getenv("BUNGIE_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing BUNGIE_KEY in environment")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-API-Key", apiKey)
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