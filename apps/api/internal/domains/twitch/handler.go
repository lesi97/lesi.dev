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

const TwitchContextKey = "Twitch_GO"

func NewTwitchHandler(logger *utils.Logger, store twitch_store.TwichStoreInterface)  *TwitchHandler {
	return &TwitchHandler{
		logger: logger,
		store: store,
	}
}


func (h *TwitchHandler) HandleGetRandomChatter(w http.ResponseWriter, r *http.Request) {
	streamerName := chi.URLParam(r, "streamer")
	if streamerName == "" {
		utils.TextResponse(w, http.StatusOK, "You must declare a streamer for this to work")
		return
	}

	username, err := h.store.RandomViewer(streamerName)
	if err != nil {
		h.logger.Errorf("ERROR: %v", err)
		// have to return ok, used by nightbot in twitch and if not ok, nightbot wont output a response
		utils.TextResponse(w, http.StatusOK, "An error has occurred, take note of the current time and tell Lesi when this happened")
		return
	}
	h.logger.Printf("%v%v hit%v\n", utils.Colours["cyan"], *username, utils.Colours["reset"])
	utils.TextResponse(w, http.StatusOK, *username)
}
