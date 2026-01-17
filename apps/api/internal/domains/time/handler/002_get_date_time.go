package handler

import (
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/time/model"
	"github.com/lesi97/lesi.dev/internal/domains/time/utils"
	su "github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) getDateTime(w http.ResponseWriter, r *http.Request) {

	zone := utils.GetGeoLocationZone(r)
	loc, err := time.LoadLocation(*zone)
	if err != nil {
		su.Error(w, http.StatusBadRequest, err)
		return
	}

	now := time.Now().In(loc)
	date := now.Format("02/01/2006")
	clock := now.Format("15:04:05")
	output := model.TimeResponseFormat{
		Date: date,
		Time: clock,
	}

	su.Success(w, http.StatusOK, output)
}