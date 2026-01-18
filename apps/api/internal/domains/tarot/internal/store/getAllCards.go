package store

import tarot_model "github.com/lesi97/lesi.dev/internal/domains/tarot/internal/model"

func (s *TarotStore) GetAllCards() *[]tarot_model.Card {
	return &s.Cards
}
