package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/httpapi"
)

func (s *Store) anilistPOST(
	ctx context.Context,
	url string,
	query string,
	variables map[string]interface{},
) ([]byte, error) {

	err := s.validateRefresh(ctx)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.env.AccessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	respBody, statusCode, err := httpapi.DoRequest(
		ctx,
		client,
		req.Method,
		req.URL.String(),
		req.Body,
		map[string]string{
			"Content-Type":  req.Header.Get("Content-Type"),
			"Authorization": req.Header.Get("Authorization"),
		},
	)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, errors.New("anilist request failed")
	}

	return respBody, nil
}
