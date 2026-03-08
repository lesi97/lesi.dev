package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	bungie_utils "github.com/lesi97/lesi.dev/internal/domains/bungie/internal/utils"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleReset(w http.ResponseWriter, r *http.Request) {
	resetLocation, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		utils.TextResponse(w, http.StatusInternalServerError, "failed to load reset timezone")
		return
	}

	now := time.Now().In(resetLocation)
	currentWeekday := now.Weekday()

	daysUntil := (int(time.Tuesday) - int(currentWeekday) + 7) % 7
	if daysUntil == 0 && (now.Hour() > 10 || (now.Hour() == 10 && (now.Minute() > 0 || now.Second() > 0))) {
		daysUntil = 7
	}

	future := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		10, 0, 0, 0,
		resetLocation,
	).AddDate(0, 0, daysUntil)

	diff := future.Sub(now)
	if diff < 0 {
		utils.TextResponse(w, http.StatusOK, "Reset's here! Maybe gift a sub ????")
		return
	}

	totalSeconds := int64(diff.Seconds())
	days := totalSeconds / (24 * 3600)
	hours := (totalSeconds % (24 * 3600)) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	parts := make([]string, 0, 4)
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d day%s", days, bungie_utils.Plural(int(days))))
	}
	parts = append(parts, fmt.Sprintf("%d hour%s", hours, bungie_utils.Plural(int(hours))))
	parts = append(parts, fmt.Sprintf("%d minute%s", minutes, bungie_utils.Plural(int(minutes))))
	parts = append(parts, fmt.Sprintf("%d second%s", seconds, bungie_utils.Plural(int(seconds))))

	final := fmt.Sprintf("%s until reset!", strings.Join(parts, ", "))

	utils.TextResponse(w, http.StatusOK, final)
}
