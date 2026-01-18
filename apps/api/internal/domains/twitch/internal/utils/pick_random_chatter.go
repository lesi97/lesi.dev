package utils

import (
	"math/rand"

	"github.com/lesi97/lesi.dev/internal/domains/twitch/internal/model"
)

func PickRandomChatter(chatters *model.TwitchChatters) *string {
	if chatters.Total == 0 {
		return nil
	}
	max := chatters.Total
	index := rand.Intn(max)
	user := chatters.Data[index].UserName
	return &user
}
