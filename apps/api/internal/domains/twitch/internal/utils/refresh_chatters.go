package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/twitch/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

func RefreshChatters(
	redis *redis.Client,
	logger *utils.Logger,
	database *db.DB,
	apiDetails *db.ApiDetails,
	baseURL string,
	clientID string,
	clientSecret string,
	authURL string,
	cacheKey string,
	streamerID string,
) {
	ctx := context.Background()

	result, err := FetchChatters(database, logger, apiDetails, baseURL, clientID, clientSecret, authURL, streamerID)
	if err != nil {
		return
	}

	wrapper := model.CachedChatters{
		FetchedAt: time.Now(),
		Data:      *result,
	}

	b, err := json.Marshal(wrapper)
	if err != nil {
		return
	}

	_ = redis.Set(ctx, cacheKey, b, 5*time.Minute).Err()
}
