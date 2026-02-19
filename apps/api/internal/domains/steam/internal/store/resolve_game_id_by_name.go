package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/lesi97/lesi.dev/internal/cache"
	"github.com/lesi97/lesi.dev/internal/domains/steam/internal/model"
	"github.com/lesi97/lesi.dev/internal/httpapi"
)

func (s *Store) ResolveGameIDByName(ctx context.Context, gameName string) (string, string, error) {
	trimmedGameName := strings.TrimSpace(gameName)
	if trimmedGameName == "" {
		return "", "", errors.New("you must provide gameName when gameId is not set")
	}
	normalisedGameName := NormaliseGameName(trimmedGameName)

	nameToIDCacheKey := GetNameToIDCacheKey(normalisedGameName)
	cachedGameID, cachedIDErr := s.Redis.Get(ctx, nameToIDCacheKey).Result()
	if cachedIDErr == nil {
		if cachedGameID == steamGameNotFoundCacheValue {
			return "", "", errors.New("game not found on Steam Store")
		}

		cachedGameName, cachedNameErr := s.Redis.Get(ctx, GetIDToNameCacheKey(cachedGameID)).Result()
		if cachedNameErr == nil && strings.TrimSpace(cachedGameName) != "" && cachedGameName != steamGameNotFoundCacheValue {
			s.Logger.PrintCache("CACHE HIT resolveGameIDByName %s", nameToIDCacheKey)
			return cachedGameID, cachedGameName, nil
		}

		s.Logger.PrintCache("CACHE HIT resolveGameIDByName %s", nameToIDCacheKey)
		return cachedGameID, trimmedGameName, nil
	}
	if !cache.IsRedisNil(cachedIDErr) {
		s.Logger.Printf("WARN: cache get failed for %s: %v", nameToIDCacheKey, cachedIDErr)
	}

	searchURL := fmt.Sprintf(
		"%s/api/storesearch/?term=%s&l=english&cc=us",
		s.StoreURL,
		url.QueryEscape(trimmedGameName),
	)
	body, statusCode, err := httpapi.DoRequest(ctx, s.HTTPClient, http.MethodGet, searchURL, nil, nil)
	if err != nil {
		return "", "", err
	}

	if statusCode < 200 || statusCode >= 300 {
		return "", "", errors.New("failed to search Steam by gameName")
	}

	var searchData model.SteamStoreSearchData
	err = json.Unmarshal(body, &searchData)
	if err != nil {
		return "", "", err
	}

	target := strings.ToLower(strings.TrimSpace(trimmedGameName))
	for _, item := range searchData.Items {
		trimmedAppName := strings.TrimSpace(item.Name)
		if trimmedAppName == "" {
			continue
		}

		lowerAppName := strings.ToLower(trimmedAppName)
		if lowerAppName == target {
			resolvedGameID := fmt.Sprintf("%d", item.ID)
			_ = s.Redis.Set(ctx, nameToIDCacheKey, resolvedGameID, steamGameNameCacheTTL).Err()
			_ = s.Redis.Set(ctx, GetIDToNameCacheKey(resolvedGameID), trimmedAppName, steamGameNameCacheTTL).Err()
			return resolvedGameID, trimmedAppName, nil
		}
	}

	for _, item := range searchData.Items {
		trimmedAppName := strings.TrimSpace(item.Name)
		if strings.Contains(strings.ToLower(trimmedAppName), target) {
			resolvedGameID := fmt.Sprintf("%d", item.ID)
			_ = s.Redis.Set(ctx, nameToIDCacheKey, resolvedGameID, steamGameNameCacheTTL).Err()
			_ = s.Redis.Set(ctx, GetIDToNameCacheKey(resolvedGameID), trimmedAppName, steamGameNameCacheTTL).Err()
			return resolvedGameID, trimmedAppName, nil
		}
	}

	_ = s.Redis.Set(ctx, nameToIDCacheKey, steamGameNotFoundCacheValue, steamGameNegativeCacheTTL).Err()
	return "", "", errors.New("game not found on Steam Store")
}
