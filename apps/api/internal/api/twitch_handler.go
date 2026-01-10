package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lesi97/lesi.dev/internal/store/twitch_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type TwitchHandler struct {
	logger         *utils.Logger
	store 			twitch_store.TwichStoreInterface
}

const TwitchContextKey = "Twitch"

func NewTwitchHandler(logger *utils.Logger, store twitch_store.TwichStoreInterface)  *TwitchHandler {
	return &TwitchHandler{
		logger: logger,
		store: store,
	}
}


func (h *TwitchHandler) HandleGetRandomChatter(w http.ResponseWriter, r *http.Request) {
	streamerName := chi.URLParam(r, "streamer")
	utils.Success(w, http.StatusOK, streamerName)
}
