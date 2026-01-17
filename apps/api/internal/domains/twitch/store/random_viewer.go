package twitch_store

import (
	"math/rand"
)

func (s *TwitchStore) RandomViewer(streamer string) (*string, error) {
	sd, err := s.getStreamerData(streamer)
	if err != nil {
		return nil, err
	}
	id := sd.Data[0].ID
	c, err := s.getChatters(id)
	if err != nil {
		return nil, err
	}
	return c.pickRandom(), nil
}

func (c *TwitchChatters) pickRandom() *string {
	if c.Total == 0 {
		return nil
	}
	max := c.Total
	index := rand.Intn(max)
	user := c.Data[index].UserName
	return &user
}