package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func fetchAndCacheBungieProfile(
	redis *redis.Client,
	logger *utils.Logger,
	clientID string,
	baseURL string,
	membershipID string,
	preferredPlatform string,
	components string,
	cacheKey string,
) (*BungieProfile, error) {
	ctx := context.Background()

	reqURL := fmt.Sprintf("%s/Platform/Destiny2/%s/Profile/%s/?components=%s", baseURL, preferredPlatform, membershipID, components)

	body, err := BungieGET(logger, clientID, reqURL)
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
