package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/model"
	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/utils"
)

func (h *Handler) updateLeaderboard(w http.ResponseWriter, r *http.Request) {
	allowedOrigin := os.Getenv("WEB_URL")

	if r.Header.Get("Origin") != allowedOrigin {
		utils.WriteJSON(w, http.StatusForbidden, model.ErrorResponse{Error: "Rejected due to CORS policy"})
		return
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil || len(raw) == 0 {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "No Data Provided"})
		return
	}

	var body model.PostBody
	if err := json.Unmarshal(raw, &body); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "No Data Provided"})
		return
	}

	body.Username = strings.TrimSpace(body.Username)
	if body.Username == "" {
		utils.WriteJSON(w, http.StatusOK, model.ErrorResponse{Error: "Invalid Data"})
		return
	}

	var presence map[string]json.RawMessage
	_ = json.Unmarshal(raw, &presence)
	hasNiceTriggers := presence["nice_triggers"] != nil
	hasMiloAttacks := presence["milo_attacks"] != nil

	now := time.Now().UTC().Format(time.RFC3339)

	update, err := h.store.UpsertUpdateUser(r.Context(), model.UpdateInput{
		Username:        body.Username,
		Score:           body.Score,
		Accuracy:        body.Accuracy,
		HasNiceTriggers: hasNiceTriggers,
		HasMiloAttacks:  hasMiloAttacks,
		UpdatedAtISO:    now,
	})
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, model.ErrorResponse{Error: "An Error Has Occurred."})
		return
	}

	utils.WriteJSON(w, http.StatusOK, update)
}