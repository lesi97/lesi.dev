package store

import (
	tarot_model "github.com/lesi97/lesi.dev/internal/domains/tarot/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type TarotStoreInterface interface {
	GetRandomTarot() *string
	GetAllCards() *[]tarot_model.Card
}

type TarotStore struct {
	logger *utils.Logger
	Cards  []tarot_model.Card
}

func NewStore(logger *utils.Logger) *TarotStore {
	return &TarotStore{
		logger: logger,
		Cards: tarotCards,
	}
}
