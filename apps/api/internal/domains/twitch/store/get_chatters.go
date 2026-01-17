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

type cachedChatters struct {
	FetchedAt time.Time      `json:"fetched_at"`
	Data      TwitchChatters `json:"data"`
}


func (s *TwitchStore) getChatters(streamerId string) (*TwitchChatters, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("twitch:chatters:%s", streamerId)

	const freshFor = 2 * time.Minute
	const staleFor = 5 * time.Minute


	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var wrapper cachedChatters
		if err := json.Unmarshal([]byte(cached), &wrapper); err == nil {
			age := time.Since(wrapper.FetchedAt)

			if age <= freshFor {
				s.Logger.Printf("%vCACHE HIT getChatters %v%v", utils.Colours["brightBlack"], cacheKey, utils.Colours["reset"])
				return &wrapper.Data, nil
			}

			if age <= staleFor {
				s.Logger.Printf("%vCACHE STALE getChatters %v%v", utils.Colours["brightBlack"], cacheKey, utils.Colours["reset"])

				go func() {
					s.refreshChatters(cacheKey, streamerId)
				}()

				return &wrapper.Data, nil
			}
		}

		_ = s.redis.Del(ctx, cacheKey).Err()
	} else {
		if err != redis.Nil {
			return nil, err
		}
	}

	result, err := s.fetchChatters(streamerId)
	if err != nil {
		return nil, err
	}

	wrapper := cachedChatters{
		FetchedAt: time.Now(),
		Data:      *result,
	}

	b, err := json.Marshal(wrapper)
	if err == nil {
		_ = s.redis.Set(ctx, cacheKey, b, staleFor).Err()
	}

	return result, nil
}

func (s *TwitchStore) fetchChatters(streamerId string) (*TwitchChatters, error) {
	const modId = "101129910"
	url := fmt.Sprintf(
		"%v/chat/chatters?first=1000&broadcaster_id=%v&moderator_id=%v",
		*s.base_url,
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

	return result, nil
}

func (s *TwitchStore) refreshChatters(cacheKey string, streamerId string) {
	ctx := context.Background()

	result, err := s.fetchChatters(streamerId)
	if err != nil {
		return
	}

	type cachedChatters struct {
		FetchedAt time.Time      `json:"fetched_at"`
		Data      TwitchChatters `json:"data"`
	}

	wrapper := cachedChatters{
		FetchedAt: time.Now(),
		Data:      *result,
	}

	b, err := json.Marshal(wrapper)
	if err != nil {
		return
	}

	_ = s.redis.Set(ctx, cacheKey, b, 5 * time.Minute).Err()
}