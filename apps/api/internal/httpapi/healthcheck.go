package httpapi

import (
	"net/http"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	utils.TextResponse(w, http.StatusOK, "healthy")
}
