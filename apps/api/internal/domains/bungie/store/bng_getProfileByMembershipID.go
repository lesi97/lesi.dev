package bungie_store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type cachedProfile struct {
	CachedAtUnix int64         `json:"cached_at_unix"`
	Value        BungieProfile `json:"value"`
}

func (s *BungieStore) getBungieProfileByMembershipID(membershipID string, preferredPlatform string, components string) (*BungieProfile, error) {
	ctx := context.Background()

	freshFor := 15 * time.Second
	staleFor := 2 * time.Minute

	cacheKey := fmt.Sprintf(
		"bungie:profile:%v",
		membershipID,
	)

	now := time.Now()

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var wrap cachedProfile
		if err := json.Unmarshal([]byte(cached), &wrap); err == nil {
			age := now.Sub(time.Unix(wrap.CachedAtUnix, 0))

			if age <= freshFor {
				s.Logger.Printf("CACHE HIT fresh getBungieProfile %s", cacheKey)
				v := wrap.Value
				return &v, nil
			}

			if age <= staleFor {
				s.Logger.Printf("CACHE HIT stale getBungieProfile %s", cacheKey)

				lockKey := cacheKey + ":lock"
				ok, _ := s.redis.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
				if ok {
					go func() {
						defer s.redis.Del(context.Background(), lockKey).Err()
						_, _ = s.fetchAndCacheBungieProfile(membershipID, preferredPlatform, components, cacheKey)
					}()
				}

				v := wrap.Value
				return &v, nil
			}
		} else {
			_ = s.redis.Del(ctx, cacheKey).Err()
		}
	} else {
		if err != redis.Nil {
			return nil, err
		}
	}

	return s.fetchAndCacheBungieProfile(membershipID, preferredPlatform, components, cacheKey)
}

func (s *BungieStore) fetchAndCacheBungieProfile(membershipID string, preferredPlatform string, components string, cacheKey string) (*BungieProfile, error) {
	ctx := context.Background()

	reqURL := fmt.Sprintf("%s/Platform/Destiny2/%s/Profile/%s/?components=%s", s.url, preferredPlatform, membershipID, components)

	body, err := s.bungieGET(reqURL)
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
		_ = s.redis.Set(ctx, cacheKey, b, 10*time.Minute).Err()
	}

	return result, nil
}