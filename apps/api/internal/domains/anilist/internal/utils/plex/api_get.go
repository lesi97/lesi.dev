package plex

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/httpapi"
)

func (s *Store) plexGET(ctx context.Context, url string) ([]byte, error) {

	headers := map[string]string{
		"Accept": "application/json",
	}
	client := &http.Client{Timeout: 10 * time.Second}
	raw, statusCode, err := httpapi.DoRequest(ctx, client, http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("plex http %d %s", statusCode, string(raw))
	}

	return raw, nil
}
