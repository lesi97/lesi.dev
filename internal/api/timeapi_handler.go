package api

import (
	"log"
	"net/http"
	"time"

	"github.com/lesi97/api.lesi.dev/internal/utils"
)

type TimeapiHandler struct {
	logger         *log.Logger
}

type timeResponseFormat struct {
	Date string `json:"date"`
	Time string `json:"time"`
}


func NewTimeapiHandler(logger *log.Logger)  *TimeapiHandler {
	return &TimeapiHandler{
		logger: logger,
	}
}

func (h *TimeapiHandler) HandleGetDateTime(w http.ResponseWriter, r *http.Request) {

	zone := r.URL.Query().Get("zone")
	if zone == "" {
		zone = "Europe/London"
	}

	loc, err := time.LoadLocation(zone)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err)
		return
	}

	now := time.Now().In(loc)

	date := now.Format("02/01/2006")
	clock := now.Format("15:04:05")

	output := timeResponseFormat{
		Date: date,
		Time: clock,
	}

	utils.Success(w, http.StatusOK, output)
}