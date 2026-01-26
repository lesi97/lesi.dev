package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lesi97/lesi.dev/internal/cache"
	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/twitch/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func GetChatters(
	redis *redis.Client,
	logger *utils.Logger,
	database *db.DB,
	apiDetails *db.ApiDetails,
	baseURL string,
	clientID string,
	clientSecret string,
	authURL string,
	streamerID string,
) (*model.TwitchChatters, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("twitch:chatters:%s", streamerID)

	const freshFor = 2 * time.Minute
	const staleFor = 5 * time.Minute

	cached, err := redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var wrapper model.CachedChatters
		if err := json.Unmarshal([]byte(cached), &wrapper); err == nil {
			age := time.Since(wrapper.FetchedAt)

			if age <= freshFor {
				logger.PrintCache("CACHE HIT getChatters %v", cacheKey)
				return &wrapper.Data, nil
			}

			if age <= staleFor {
				logger.PrintCache("CACHE STALE getChatters %v", cacheKey)

				go func() {
					RefreshChatters(
						redis,
						logger,
						database,
						apiDetails,
						baseURL,
						clientID,
						clientSecret,
						authURL,
						cacheKey,
						streamerID,
					)
				}()

				return &wrapper.Data, nil
			}
		}

		_ = redis.Del(ctx, cacheKey).Err()
	} else if !cache.IsRedisNil(err) {
		return nil, err
	}

	result, err := FetchChatters(database, logger, apiDetails, baseURL, clientID, clientSecret, authURL, streamerID)
	if err != nil {
		return nil, err
	}

	wrapper := model.CachedChatters{
		FetchedAt: time.Now(),
		Data:      *result,
	}

	b, err := json.Marshal(wrapper)
	if err == nil {
		_ = redis.Set(ctx, cacheKey, b, staleFor).Err()
	}

	return result, nil
}
