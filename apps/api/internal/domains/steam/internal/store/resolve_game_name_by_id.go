package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lesi97/lesi.dev/internal/cache"
	"github.com/lesi97/lesi.dev/internal/domains/steam/internal/model"
	"github.com/lesi97/lesi.dev/internal/httpapi"
)

func (s *Store) ResolveGameNameByID(ctx context.Context, gameID string) (string, error) {
	idToNameCacheKey := GetIDToNameCacheKey(gameID)
	cachedGameName, cachedNameErr := s.Redis.Get(ctx, idToNameCacheKey).Result()
	if cachedNameErr == nil {
		if cachedGameName == steamGameNotFoundCacheValue {
			_ = s.Redis.Expire(ctx, idToNameCacheKey, steamGameNegativeCacheTTL).Err()
			return "", errors.New("gameId not found on Steam Store")
		}
		if strings.TrimSpace(cachedGameName) != "" {
			_ = s.Redis.Expire(ctx, idToNameCacheKey, steamGameNameCacheTTL).Err()
			_ = s.Redis.Expire(ctx, GetNameToIDCacheKey(NormaliseGameName(cachedGameName)), steamGameNameCacheTTL).Err()
			s.Logger.PrintCache("CACHE HIT resolveGameNameByID %s", idToNameCacheKey)
			return cachedGameName, nil
		}
	}
	if !cache.IsRedisNil(cachedNameErr) {
		s.Logger.Printf("WARN: cache get failed for %s: %v", idToNameCacheKey, cachedNameErr)
	}

	url := fmt.Sprintf("%s/api/appdetails?appids=%s&l=english&cc=us", s.StoreURL, gameID)
	body, statusCode, err := httpapi.DoRequest(ctx, s.HTTPClient, http.MethodGet, url, nil, nil)
	if err != nil {
		return "", err
	}

	if statusCode < 200 || statusCode >= 300 {
		return "", errors.New("failed to fetch game details from Steam Store")
	}

	var detailsData model.SteamAppDetailsByIDData
	err = json.Unmarshal(body, &detailsData)
	if err != nil {
		return "", err
	}

	details, ok := detailsData[gameID]
	if !ok || !details.Success {
		_ = s.Redis.Set(ctx, idToNameCacheKey, steamGameNotFoundCacheValue, steamGameNegativeCacheTTL).Err()
		return "", errors.New("gameId not found on Steam Store")
	}

	trimmedName := strings.TrimSpace(details.Data.Name)
	if trimmedName == "" {
		_ = s.Redis.Set(ctx, idToNameCacheKey, steamGameNotFoundCacheValue, steamGameNegativeCacheTTL).Err()
		return "", errors.New("gameId has no game name on Steam Store")
	}

	_ = s.Redis.Set(ctx, idToNameCacheKey, trimmedName, steamGameNameCacheTTL).Err()
	normalisedGameName := NormaliseGameName(trimmedName)
	_ = s.Redis.Set(ctx, GetNameToIDCacheKey(normalisedGameName), gameID, steamGameNameCacheTTL).Err()

	return trimmedName, nil
}
