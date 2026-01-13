package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/store/aim_trainer_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type AimTrainerHandler struct {
	logger         *utils.Logger
	store 			aim_trainer_store.AimTrainerStoreInterface
}

const aimTrainerContextKey = "aim-trainer"

func NewAimTrainerHandler(logger *utils.Logger, store aim_trainer_store.AimTrainerStoreInterface)  *AimTrainerHandler {
	return &AimTrainerHandler{
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

// Legacy support for aim-trainer game which expects a specific resposne
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}


func (h *AimTrainerHandler) HandleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	rows, err := h.store.GetLeaderboard(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Error: err.Error(),
		})
		return
	}

	if len(rows) == 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "No Data",
		})
		return
	}

	writeJSON(w, http.StatusOK, leaderboardResponse{
		Data: rows,
	})
}

type aimTrainerpostBody struct {
	Username string   `json:"username"`
	Score    *float64   `json:"score,omitempty"`
	Accuracy *float64 `json:"accuracy,omitempty"`
}


func (h *AimTrainerHandler) HandleUpdateLeaderboard(w http.ResponseWriter, r *http.Request) {
	allowedOrigin := "https://lesi.dev"

	if r.Header.Get("Origin") != allowedOrigin {
		writeJSON(w, http.StatusForbidden, errorResponse{Error: "Rejected due to CORS policy"})
		return
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil || len(raw) == 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "No Data Provided"})
		return
	}

	var body aimTrainerpostBody
	if err := json.Unmarshal(raw, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "No Data Provided"})
		return
	}

	body.Username = strings.TrimSpace(body.Username)
	if body.Username == "" {
		writeJSON(w, http.StatusOK, errorResponse{Error: "Invalid Data"})
		return
	}

	var presence map[string]json.RawMessage
	_ = json.Unmarshal(raw, &presence)
	hasNiceTriggers := presence["nice_triggers"] != nil
	hasMiloAttacks := presence["milo_attacks"] != nil

	now := time.Now().UTC().Format(time.RFC3339)

	update, err := h.store.UpsertUpdateUser(r.Context(), aim_trainer_store.UpdateInput{
		Username:        body.Username,
		Score:           body.Score,
		Accuracy:        body.Accuracy,
		HasNiceTriggers: hasNiceTriggers,
		HasMiloAttacks:  hasMiloAttacks,
		UpdatedAtISO:    now,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "An Error Has Occurred."})
		return
	}

	writeJSON(w, http.StatusOK, update)
}