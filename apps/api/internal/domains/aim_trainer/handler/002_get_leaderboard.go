package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/model"
	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/utils"
)

func (h *Handler) getLeaderboard(w http.ResponseWriter, r *http.Request) {
	rows, err := h.store.GetLeaderboard(r.Context())
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if len(rows) == 0 {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{
			Error: "No Data",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, model.LeaderboardResponse{
		Data: rows,
	})
}