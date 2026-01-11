package anilist_store

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

func (s *AnilistStore) anilistPOST(url string, query string, variables map[string]interface{}) ([]byte, error) {
	defer s.Logger.LogExecutionTime(url, time.Now())

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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("anilist request failed")
	}

	return io.ReadAll(resp.Body)
}
