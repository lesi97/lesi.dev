package twitch_store

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	const modId = "101129910" // me :)
	url := fmt.Sprintf("%v/chat/chatters?first=1000&broadcaster_id=%v&moderator_id=%v", s.base_url, streamerId, modId)
	body, err := s.twitchGET(url)
	if err != nil {
		return nil, err
	}

	result := &TwitchChatters{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}

	return result, nil
}
