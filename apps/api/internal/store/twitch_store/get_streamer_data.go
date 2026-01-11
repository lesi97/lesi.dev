package twitch_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
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
	url := fmt.Sprintf("%v/users?login=%v", s.base_url, streamer)
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

	return result, nil
}
