package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func GetBungieProfileByMembershipID(
	ctx context.Context,
	redis *redis.Client,
	logger *utils.Logger,
	clientID string,
	baseURL string,
	membershipID string,
	preferredPlatform string,
	components string,
) (*BungieProfile, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	freshFor := 15 * time.Second
	staleFor := 2 * time.Minute

	normalisedComponents := strings.ReplaceAll(strings.TrimSpace(components), " ", "")

	cacheKey := fmt.Sprintf(
		"bungie:profile:%v:%v:%v",
		preferredPlatform,
		membershipID,
		normalisedComponents,
	)

	now := time.Now()

	cached, err := redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var wrap cachedProfile
		if err := json.Unmarshal([]byte(cached), &wrap); err == nil {
			age := now.Sub(time.Unix(wrap.CachedAtUnix, 0))

			if age <= freshFor {
				logger.PrintCache("CACHE HIT fresh getBungieProfile %s", cacheKey)
				v := wrap.Value
				return &v, nil
			}

			if age <= staleFor {
				logger.PrintCache("CACHE HIT stale getBungieProfile %s", cacheKey)

				lockKey := cacheKey + ":lock"
				ok, _ := redis.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
				if ok {
					go func() {
						defer redis.Del(context.Background(), lockKey).Err()
						if ctx.Err() != nil {
							return
						}
						_, _ = fetchAndCacheBungieProfile(ctx, redis, logger, clientID, baseURL, membershipID, preferredPlatform, normalisedComponents, cacheKey)
					}()
				}

				v := wrap.Value
				return &v, nil
			}
		} else {
			_ = redis.Del(ctx, cacheKey).Err()
		}
	} else if !IsRedisNil(err) {
		return nil, err
	}

	return fetchAndCacheBungieProfile(ctx, redis, logger, clientID, baseURL, membershipID, preferredPlatform, normalisedComponents, cacheKey)
}
