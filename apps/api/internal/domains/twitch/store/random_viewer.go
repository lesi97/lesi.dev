package store

import (
	"fmt"

	twitch_utils "github.com/lesi97/lesi.dev/internal/domains/twitch/utils"
)

func (s *Store) RandomViewer(streamer string) (*string, error) {
	streamerData, err := twitch_utils.GetStreamerData(
		s.Redis,
		s.Logger,
		s.DB,
		s.ApiDetails,
		s.BaseURL,
		s.ClientID,
		s.ClientSecret,
		s.AuthURL,
		streamer,
	)
	if err != nil {
		return nil, err
	}

	if len(streamerData.Data) == 0 {
		return nil, fmt.Errorf("streamer not found")
	}

	streamerID := streamerData.Data[0].ID
	chatters, err := twitch_utils.GetChatters(
		s.Redis,
		s.Logger,
		s.DB,
		s.ApiDetails,
		s.BaseURL,
		s.ClientID,
		s.ClientSecret,
		s.AuthURL,
		streamerID,
	)
	if err != nil {
		return nil, err
	}

	chatter := twitch_utils.PickRandomChatter(chatters)
	if chatter == nil {
		return nil, fmt.Errorf("no chatters available")
	}

	return chatter, nil
}
