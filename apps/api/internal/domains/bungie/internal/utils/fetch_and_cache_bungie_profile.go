package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func fetchAndCacheBungieProfile(
	ctx context.Context,
	redis *redis.Client,
	logger *utils.Logger,
	httpClient *http.Client,
	clientID string,
	baseURL string,
	membershipID string,
	preferredPlatform string,
	components string,
	cacheKey string,
) (*BungieProfile, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	reqURL := fmt.Sprintf("%s/Platform/Destiny2/%s/Profile/%s/?components=%s", baseURL, preferredPlatform, membershipID, components)

	body, err := BungieGET(ctx, logger, httpClient, clientID, reqURL)
	if err != nil {
		return nil, err
	}

	result := &BungieProfile{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	wrap := cachedProfile{
		CachedAtUnix: time.Now().Unix(),
		Value:        *result,
	}

	b, err := json.Marshal(wrap)
	if err == nil {
		_ = redis.Set(ctx, cacheKey, b, 10*time.Minute).Err()
	}

	return result, nil
}
