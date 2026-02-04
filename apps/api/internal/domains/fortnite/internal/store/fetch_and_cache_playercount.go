package store

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/lesi97/lesi.dev/internal/httpapi"
)

type cachedPlayerCount struct {
	CachedAtUnix int64  `json:"cached_at_unix"`
	Value        string `json:"value"`
}


func (s *Store) fetchAndCachePlayerCount(ctx context.Context, cacheKey string, staleFor time.Duration, url string) (*string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	headers := map[string]string{
		"Accept": "*/*",
		"Accept-Encoding": "gzip, deflate, br, zstd",
		"Accept-Language": "en-GB,en-US;q=0.9,en;q=0.8,ja;q=0.7",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36",
		"Referer": s.BaseURL,
		"Origin": s.BaseURL,
	}
	body, statusCode, err := httpapi.DoRequest(ctx, s.HTTPClient, http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("failed to fetch fortnite playercount")
	}

	var res FortnitePlayerCount
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if len(res.Data.Values) == 0 {
		return nil, fmt.Errorf("no data in fortnite payload")
	}

	last := res.Data.Values[len(res.Data.Values)-1]

	message := fmt.Sprintf(
		"There are currently %v players in Fortnite across all platforms | Fortnite data pulled from %s",
		humanize.Comma(int64(last)),
		s.BaseURL,
	)

	wrap := cachedPlayerCount{
		CachedAtUnix: time.Now().Unix(),
		Value:        message,
	}
	if b, err := json.Marshal(wrap); err == nil {
		_ = s.Redis.Set(ctx, cacheKey, b, staleFor).Err()
	}

	return &message, nil
}
