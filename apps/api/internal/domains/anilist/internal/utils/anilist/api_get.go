package anilist

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/httpapi"
)

func (s *Store) anilistGET(ctx context.Context, url string) ([]byte, error) {

	err := s.validateRefresh(ctx)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + s.env.AccessToken,
	}
	body, statusCode, err := httpapi.DoRequest(ctx, &http.Client{}, http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d: %s", statusCode, string(body))
	}

	return body, nil
}
