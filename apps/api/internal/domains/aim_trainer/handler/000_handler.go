package handler

import (
	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Handler struct {
	logger         *utils.Logger
	store 			store.Methods
}

const aimTrainerContextKey = "aim-trainer"

func NewHandler(logger *utils.Logger, db *db.DB)  *Handler {

	store := store.NewStore(db, logger)
	

	return &Handler{
		logger: logger,
		store: store,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

type leaderboardResponse struct {
	Data []aim_trainer_store.LeaderboardRow `json:"data"`
}





