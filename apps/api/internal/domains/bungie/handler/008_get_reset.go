package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) HandleReset(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	loc := now.Location()
	currentWeekday := now.Weekday()

	daysUntil := (int(time.Tuesday) - int(currentWeekday) + 7) % 7
	if daysUntil == 0 && now.Hour() >= 17 {
		daysUntil = 7
	}

	future := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		17, 0, 0, 0,
		loc,
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
		parts = append(parts, fmt.Sprintf("%d days", days))
	}
	parts = append(parts, fmt.Sprintf("%d hours", hours))
	parts = append(parts, fmt.Sprintf("%d minutes", minutes))
	parts = append(parts, fmt.Sprintf("%d seconds", seconds))

	final := fmt.Sprintf("%s until reset!", strings.Join(parts, ", "))

	utils.TextResponse(w, http.StatusOK, final)
}
