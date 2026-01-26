package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/lesi97/lesi.dev/internal/cache"
	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func getUserFromBungieByGamertag(
	ctx context.Context,
	database *db.DB,
	logger *utils.Logger,
	redis *redis.Client,
	baseURL string,
	clientID string,
	id string,
	platform string,
) (*bungieSearch, error) {
	platformEnum := switchPlatforms(platform)
	cacheKey := fmt.Sprintf("bungie:searchdestinyplayer:%s:%s", id, platformEnum)

	cached, err := redis.Get(ctx, cacheKey).Result()
	if err == nil {
		result := &bungieSearch{}
		if err := json.Unmarshal([]byte(cached), result); err == nil {
			logger.PrintCache("CACHE HIT getUserFromBungieByGamertag %s", cacheKey)
			return result, nil
		}
		_ = redis.Del(ctx, cacheKey).Err()
	} else if !cache.IsRedisNil(err) {
		return nil, err
	}

	escapedID := url.PathEscape(id)
	reqURL := fmt.Sprintf("%s/Platform/Destiny2/SearchDestinyPlayer/%s/%s/", baseURL, platformEnum, escapedID)

	body, err := BungieGET(ctx, logger, clientID, reqURL)
	if err != nil {
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil {
			fmt.Printf("ERROR in getUserFromBungieByGamertag: %v\n", err.Error())
		}
		return nil, err
	}

	result := &bungieSearch{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}

	b, err := json.Marshal(result)
	if err == nil {
		_ = redis.Set(ctx, cacheKey, b, 168*time.Hour).Err()
	}

	go func() {
		if len(result.Response) > 0 {
			user := result.Response[0]
			InsertDestinyUser(database, logger, user.MembershipID, id, int64(user.MembershipType), user.BungieGlobalDisplayName)
		}
	}()

	return result, nil
}
