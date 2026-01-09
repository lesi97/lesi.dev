package api

import (
	"encoding/json"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type LocalHandler struct {
	logger 	*utils.Logger
	db 		*database.DB
}

type UpdateApiDetailsRequest struct {
	Name               string `json:"name"`
	ClientID           string `json:"client_id"`
	ClientSecret       string `json:"client_secret"`
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	RefreshTokenExpiry int64  `json:"refresh_token_expiry"`
	BaseURL            string `json:"base_url"`
	RedirectURL        string `json:"redirect_url"`
}

func NewLocalHandler(logger *utils.Logger, db *database.DB) *LocalHandler {
	return &LocalHandler{
		logger: logger,
		db: db,
	}
}


func (h *LocalHandler) HandleDbDump(w http.ResponseWriter, r *http.Request) {

	var req UpdateApiDetailsRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	utils.PrintPrettyJSON(req)

	err := h.db.UpdateApiDetails(
		r.Context(),
		req.Name,
		&req.ClientID,
		&req.ClientSecret,
		&req.AccessToken,
		&req.RefreshToken,
		&req.RefreshTokenExpiry,
		&req.BaseURL,
		&req.RedirectURL,
		h.logger,
	)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.Success(w, http.StatusNoContent, nil)

}