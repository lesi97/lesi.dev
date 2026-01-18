package handler

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/domains/local/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleDbDump(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeUpdateApiDetailsRequest(r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	utils.PrintPrettyJSON(req)

	err = h.store.UpdateApiDetails(r.Context(), model.UpdateApiDetailsInput{
		Name:               req.Name,
		ClientID:           req.ClientID,
		ClientSecret:       req.ClientSecret,
		AccessToken:        req.AccessToken,
		RefreshToken:       req.RefreshToken,
		RefreshTokenExpiry: req.RefreshTokenExpiry,
		BaseURL:            req.BaseURL,
		RedirectURL:        req.RedirectURL,
	})
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	utils.Success(w, http.StatusNoContent, nil)
}
