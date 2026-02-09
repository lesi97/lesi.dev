package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/twitch/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func FetchChatters(
	ctx context.Context,
	database *db.DB,
	logger *utils.Logger,
	apiDetails *db.ApiDetails,
	baseURL string,
	clientID string,
	clientSecret string,
	authURL string,
	streamerID string,
) (*model.TwitchChatters, error) {
	const modID = "101129910"
	url := fmt.Sprintf(
		"%v/chat/chatters?first=1000&broadcaster_id=%v&moderator_id=%v",
		baseURL,
		streamerID,
		modID,
	)

	body, err := TwitchGET(ctx, database, logger, apiDetails, clientID, clientSecret, authURL, url)
	if err != nil {
		return nil, err
	}

	result := &model.TwitchChatters{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
