package twitch_store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type StreamerData struct {
	Data []struct {
		ID              string    `json:"id"`
		Login           string    `json:"login"`
		DisplayName     string    `json:"display_name"`
		Type            string    `json:"type"`
		BroadcasterType string    `json:"broadcaster_type"`
		Description     string    `json:"description"`
		ProfileImageURL string    `json:"profile_image_url"`
		OfflineImageURL string    `json:"offline_image_url"`
		ViewCount       int       `json:"view_count"`
		CreatedAt       time.Time `json:"created_at"`
	} `json:"data"`
}

func (s *TwitchStore) getStreamerData(streamer string) (*StreamerData, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("twitch:user:%s", strings.ToLower(streamer))
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		result := &StreamerData{}
		if err := json.Unmarshal([]byte(cached), result); err == nil {
			s.Logger.Printf("%vCACHE HIT getStreamerData %v%v", utils.Colours["brightBlack"], cacheKey, utils.Colours["reset"])
			return result, nil
		}
		_ = s.redis.Del(ctx, cacheKey).Err()
	} else {
		if err != redis.Nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%v/users?login=%v", *s.base_url, streamer)
	body, err := s.twitchGET(url)
	if err != nil {
		return nil, err
	}

	result := &StreamerData{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}
	b, err := json.Marshal(result)
	if err == nil {
		day := 24 * time.Hour
		_ = s.redis.Set(ctx, cacheKey, b, 30 * day).Err()
	}

	return result, nil
}
