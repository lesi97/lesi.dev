package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/twitch/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func GetStreamerData(
	redis *redis.Client,
	logger *utils.Logger,
	database *db.DB,
	apiDetails *db.ApiDetails,
	baseURL string,
	clientID string,
	clientSecret string,
	authURL string,
	streamer string,
) (*model.StreamerData, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("twitch:user:%s", strings.ToLower(streamer))
	cached, err := redis.Get(ctx, cacheKey).Result()
	if err == nil {
		result := &model.StreamerData{}
		if err := json.Unmarshal([]byte(cached), result); err == nil {
			logger.Printf("%vCACHE HIT getStreamerData %v%v", utils.Colours["brightBlack"], cacheKey, utils.Colours["reset"])
			return result, nil
		}
		_ = redis.Del(ctx, cacheKey).Err()
	} else if !IsRedisNil(err) {
		return nil, err
	}

	url := fmt.Sprintf("%v/users?login=%v", baseURL, streamer)
	body, err := TwitchGET(database, logger, apiDetails, clientID, clientSecret, authURL, url)
	if err != nil {
		return nil, err
	}

	result := &model.StreamerData{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}
	b, err := json.Marshal(result)
	if err == nil {
		day := 24 * time.Hour
		_ = redis.Set(ctx, cacheKey, b, 30*day).Err()
	}

	return result, nil
}
