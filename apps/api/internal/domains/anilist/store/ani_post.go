package anilist_store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (s *AnilistStore) anilistPOST(url string, query string, variables map[string]interface{}) ([]byte, error) {
	defer s.Logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", url), time.Now(), nil)

	body, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.graphql_url,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if s.api_details != nil && s.api_details.AccessToken != nil {
		req.Header.Set("Authorization", "Bearer "+*s.api_details.AccessToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s.Logger.Errorf("anilist request failed\n")
		s.Logger.Errorf("status_code: %v\n", resp.StatusCode)
		s.Logger.Errorf("response_body: %v\n", string(respBody))
		return nil, errors.New("anilist request failed")
	}

	return io.ReadAll(resp.Body)
}
