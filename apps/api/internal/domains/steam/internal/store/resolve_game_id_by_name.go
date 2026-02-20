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
	decodedGameName := DecodeGameName(gameName)
	trimmedGameName := strings.TrimSpace(decodedGameName)
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
			_ = s.Redis.Expire(ctx, nameToIDCacheKey, steamGameNameCacheTTL).Err()
			_ = s.Redis.Expire(ctx, GetIDToNameCacheKey(cachedGameID), steamGameNameCacheTTL).Err()
			s.Logger.PrintCache("CACHE HIT resolveGameIDByName %s", nameToIDCacheKey)
			return cachedGameID, cachedGameName, nil
		}

		_ = s.Redis.Expire(ctx, nameToIDCacheKey, steamGameNameCacheTTL).Err()
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

	normalisedTarget := NormaliseGameName(trimmedGameName)
	for _, item := range searchData.Items {
		trimmedAppName := strings.TrimSpace(item.Name)
		if trimmedAppName == "" {
			continue
		}

		normalisedAppName := NormaliseGameName(trimmedAppName)
		if normalisedAppName == normalisedTarget {
			resolvedGameID := fmt.Sprintf("%d", item.ID)
			_ = s.Redis.Set(ctx, nameToIDCacheKey, resolvedGameID, steamGameNameCacheTTL).Err()
			_ = s.Redis.Set(ctx, GetIDToNameCacheKey(resolvedGameID), trimmedAppName, steamGameNameCacheTTL).Err()
			return resolvedGameID, trimmedAppName, nil
		}
	}

	for _, item := range searchData.Items {
		trimmedAppName := strings.TrimSpace(item.Name)
		normalisedAppName := NormaliseGameName(trimmedAppName)
		if strings.Contains(normalisedAppName, normalisedTarget) || strings.Contains(normalisedTarget, normalisedAppName) {
			resolvedGameID := fmt.Sprintf("%d", item.ID)
			_ = s.Redis.Set(ctx, nameToIDCacheKey, resolvedGameID, steamGameNameCacheTTL).Err()
			_ = s.Redis.Set(ctx, GetIDToNameCacheKey(resolvedGameID), trimmedAppName, steamGameNameCacheTTL).Err()
			return resolvedGameID, trimmedAppName, nil
		}
	}

	targetTokens := strings.Fields(normalisedTarget)
	bestMatchID := int64(0)
	bestMatchName := ""
	bestScore := 0
	for _, item := range searchData.Items {
		trimmedAppName := strings.TrimSpace(item.Name)
		if trimmedAppName == "" {
			continue
		}

		normalisedAppName := NormaliseGameName(trimmedAppName)
		if normalisedAppName == "" {
			continue
		}

		matchedTokenCount := 0
		for _, token := range targetTokens {
			if token != "" && strings.Contains(normalisedAppName, token) {
				matchedTokenCount++
			}
		}

		if matchedTokenCount > bestScore {
			bestScore = matchedTokenCount
			bestMatchID = item.ID
			bestMatchName = trimmedAppName
		}
	}

	if bestScore > 0 && bestMatchName != "" {
		resolvedGameID := fmt.Sprintf("%d", bestMatchID)
		_ = s.Redis.Set(ctx, nameToIDCacheKey, resolvedGameID, steamGameNameCacheTTL).Err()
		_ = s.Redis.Set(ctx, GetIDToNameCacheKey(resolvedGameID), bestMatchName, steamGameNameCacheTTL).Err()
		return resolvedGameID, bestMatchName, nil
	}

	fallbackGameID, fallbackGameName, fallbackErr := s.ResolveGameIDByNameFromAppList(ctx, trimmedGameName)
	if fallbackErr == nil && fallbackGameID != "" && fallbackGameName != "" {
		_ = s.Redis.Set(ctx, nameToIDCacheKey, fallbackGameID, steamGameNameCacheTTL).Err()
		_ = s.Redis.Set(ctx, GetIDToNameCacheKey(fallbackGameID), fallbackGameName, steamGameNameCacheTTL).Err()
		return fallbackGameID, fallbackGameName, nil
	}

	knownGameID, hasKnownGameID := GetKnownGameIDByName(normalisedTarget)
	if hasKnownGameID {
		knownGameName, knownNameErr := s.ResolveGameNameByID(ctx, knownGameID)
		if knownNameErr != nil {
			return "", "", knownNameErr
		}

		_ = s.Redis.Set(ctx, nameToIDCacheKey, knownGameID, steamGameNameCacheTTL).Err()
		_ = s.Redis.Set(ctx, GetIDToNameCacheKey(knownGameID), knownGameName, steamGameNameCacheTTL).Err()
		return knownGameID, knownGameName, nil
	}

	_ = s.Redis.Set(ctx, nameToIDCacheKey, steamGameNotFoundCacheValue, steamGameNegativeCacheTTL).Err()
	return "", "", errors.New("game not found on Steam Store")
}
