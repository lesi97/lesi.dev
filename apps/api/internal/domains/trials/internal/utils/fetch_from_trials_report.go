package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lesi97/lesi.dev/internal/cache"
	"github.com/lesi97/lesi.dev/internal/domains/trials/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

const trialsReportCacheKey = "trials:report"

func FetchFromTrialsReport(ctx context.Context, logger *utils.Logger, url string, redis *redis.Client) (*model.TrialsData, error) {
	now := time.Now()
	shouldLog := false
	startedAt := time.Now()
	defer func() {
		if shouldLog {
			logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", url), startedAt, nil)
		}
	}()

	cacheKey := trialsReportCacheKey
	if redis != nil {
		cached, err := redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var wrap cachedTrialsReport
			if err := json.Unmarshal([]byte(cached), &wrap); err == nil {
				freshFor, staleFor := GetTrialsReportCacheDurations(now, &wrap.Value)
				age := now.Sub(time.Unix(wrap.CachedAtUnix, 0))
				if freshFor > 0 && age <= freshFor {
					logger.PrintCache("CACHE HIT fresh fetchFromTrialsReport %s", cacheKey)
					return &wrap.Value, nil
				}
				if staleFor > 0 && age <= staleFor {
					logger.PrintCache("CACHE HIT stale fetchFromTrialsReport %s", cacheKey)
					lockKey := cacheKey + ":lock"
					ok, _ := redis.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
					if ok {
						go func() {
							defer redis.Del(context.Background(), lockKey).Err()
							data, err := FetchTrialsReportFromAPI(context.Background(), url)
							if err != nil {
								return
							}
							_, ttl := GetTrialsReportCacheDurations(time.Now(), data)
							if ttl <= 0 {
								return
							}
							wrap := cachedTrialsReport{
								CachedAtUnix: time.Now().Unix(),
								Value:        *data,
							}
							if b, err := json.Marshal(wrap); err == nil {
								_ = redis.Set(context.Background(), cacheKey, b, ttl).Err()
							}
						}()
					}
					return &wrap.Value, nil
				}
			} else {
				_ = redis.Del(ctx, cacheKey).Err()
			}
		} else if !cache.IsRedisNil(err) {
			return nil, err
		}
	}

	if !IsTrialsReportAvailable(now) {
		return nil, fmt.Errorf("trials report is not available")
	}

	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	shouldLog = true
	result, err := FetchTrialsReportFromAPI(reqCtx, url)
	if err != nil {
		return nil, err
	}

	if redis != nil {
		_, ttl := GetTrialsReportCacheDurations(now, result)
		if ttl > 0 {
			wrap := cachedTrialsReport{
				CachedAtUnix: time.Now().Unix(),
				Value:        *result,
			}
			if b, err := json.Marshal(wrap); err == nil {
				_ = redis.Set(ctx, cacheKey, b, ttl).Err()
			}
		}
	}

	return result, nil
}
