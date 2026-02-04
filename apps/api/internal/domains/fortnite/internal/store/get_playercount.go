package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lesi97/lesi.dev/internal/cache"
)

type FortnitePlayerCount = struct {
	Success bool `json:"success"`
	Data    struct {
		Start     int   `json:"start"`
		Step      int   `json:"step"`
		Values    []int `json:"values"`
		ValuesAvg []int `json:"values_avg"`
		Max24H    int   `json:"max_24h"`
	} `json:"data"`
}

func (s *Store) GetPlayerCount(ctx context.Context) (*string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	url := s.BaseURL + "/player-count-graph?range=1d&id=-1"
	startedAt := time.Now()
	shouldLog := true
	defer func() {
		if shouldLog && ctx.Err() == nil {
			s.Logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", url), startedAt, ctx)
		}
	}()

	freshFor := 30 * time.Second
	staleFor := 5 * time.Minute

	cacheKey := "fortnite:playercount"

	now := time.Now()

	cached, err := s.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var wrap cachedPlayerCount
		if err := json.Unmarshal([]byte(cached), &wrap); err == nil {
			age := now.Sub(time.Unix(wrap.CachedAtUnix, 0))

			if age <= freshFor {
				s.Logger.PrintCache("CACHE HIT fresh getPlayerCount %s", cacheKey)
				shouldLog = false
				v := wrap.Value
				return &v, nil
			}

			if age <= staleFor {
				s.Logger.PrintCache("CACHE HIT stale getPlayerCount %s", cacheKey)
				shouldLog = false

				lockKey := cacheKey + ":lock"
				ok, _ := s.Redis.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
				if ok {
					go func() {
						defer s.Redis.Del(context.Background(), lockKey).Err()
						if ctx.Err() != nil {
							return
						}
						_, _ = s.fetchAndCachePlayerCount(ctx, cacheKey, staleFor, url)
					}()
				}

				v := wrap.Value
				return &v, nil
			}
		} else {
			_ = s.Redis.Del(ctx, cacheKey).Err()
		}
	} else if !cache.IsRedisNil(err) {
		return nil, err
	}

	return s.fetchAndCachePlayerCount(ctx, cacheKey, staleFor, url)
}
