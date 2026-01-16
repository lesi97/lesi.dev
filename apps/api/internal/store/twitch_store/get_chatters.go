package twitch_store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type TwitchUser struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type TwitchChatters struct {
	Data []TwitchUser `json:"data"`
	Pagination struct {} `json:"pagination"`
	Total int `json:"total"`
}

func (s *TwitchStore) getChatters(streamerId string) (*TwitchChatters, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("twitch:chatters:%s", streamerId)

	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		result := &TwitchChatters{}
		if err := json.Unmarshal([]byte(cached), result); err == nil {
			s.Logger.Printf("%vCACHE HIT getChatters %v%v", utils.Colours["brightBlack"], cacheKey, utils.Colours["reset"])
			return result, nil
		}
		_ = s.redis.Del(ctx, cacheKey).Err()
	} else {
		if err != redis.Nil {
			return nil, err
		}
	}

	const modId = "101129910"
	url := fmt.Sprintf(
		"%v/chat/chatters?first=1000&broadcaster_id=%v&moderator_id=%v",
		s.base_url,
		streamerId,
		modId,
	)

	body, err := s.twitchGET(url)
	if err != nil {
		return nil, err
	}

	result := &TwitchChatters{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	b, err := json.Marshal(result)
	if err == nil {
		_ = s.redis.Set(ctx, cacheKey, b, time.Minute).Err()
	}

	return result, nil
}