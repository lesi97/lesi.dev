package tarot_store

import (
	"fmt"
	"math/rand"

	"github.com/lesi97/lesi.dev/internal/utils"
)

type TarotStoreInterface interface {
	GetRandomTarot() *string
	GetAllCards() *[]Card
}

type TarotStore struct {
	logger *utils.Logger
	Cards  []Card
}

func NewStore(logger *utils.Logger) *TarotStore {
	return &TarotStore{
		logger: logger,
		Cards: tarotCards,
	}
}

func (s *TarotStore) GetRandomTarot() *string {
	max := len(s.Cards)
	index := rand.Intn(max)
	card := s.Cards[index]
	message := fmt.Sprintf("%s | %s", card.Name, card.Desc)
	utils.TruncateString(&message, 400)
	return &message
}

func (s *TarotStore) GetAllCards() *[]Card {
	return &s.Cards
}

