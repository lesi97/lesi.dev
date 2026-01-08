package api

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/store/anilist_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type AnilistHandler struct {
	logger         *utils.Logger
	store 			anilist_store.AnilistStoreInterface
}

const anilistContextKey = "anilist"

func NewAnilistHandler(logger *utils.Logger, store anilist_store.AnilistStoreInterface)  *AnilistHandler {
	return &AnilistHandler{
		logger: logger,
		store: store,
	}
}


func (h *AnilistHandler) HandleAnilistTest(w http.ResponseWriter, r *http.Request) {
	h.store.Test()

	utils.TextResponse(w, http.StatusOK, "test")
}