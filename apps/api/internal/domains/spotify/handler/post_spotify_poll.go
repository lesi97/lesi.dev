package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
	domain_utils "github.com/lesi97/lesi.dev/internal/domains/spotify/internal/utils"
	shared_utils "github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) PostSpotifyPoll(w http.ResponseWriter, r *http.Request) {
	if err := domain_utils.ValidateScrobbleAccess(r); err != nil {
		h.logger.Printf("Spotify scrobble poll access denied: %v", err)
		shared_utils.Error(w, http.StatusForbidden, "scrobble access denied")
		return
	}

	input, err := spotifyPollInputFromRequest(r)
	if err != nil {
		shared_utils.Error(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.store.PollSpotifyRecentlyPlayed(r.Context(), input)
	if err != nil {
		h.logger.Printf("ERROR: Spotify scrobble poll: %v", err)
		shared_utils.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	shared_utils.Success(w, http.StatusOK, result)
}

func spotifyPollInputFromRequest(r *http.Request) (model.SpotifyPollInput, error) {
	input := model.SpotifyPollInput{
		Limit: 50,
	}
	query := r.URL.Query()

	if rawLimit := query.Get("limit"); rawLimit != "" {
		limit, err := strconv.Atoi(rawLimit)
		if err != nil || limit < 1 || limit > 50 {
			return model.SpotifyPollInput{}, fmt.Errorf("limit must be between 1 and 50")
		}
		input.Limit = limit
	}

	rawAfter := query.Get("after")
	if rawAfter == "" {
		rawAfter = query.Get("after_ms")
	}
	if rawAfter != "" {
		afterMS, err := strconv.ParseInt(rawAfter, 10, 64)
		if err != nil || afterMS < 0 {
			return model.SpotifyPollInput{}, fmt.Errorf("after must be a unix timestamp in milliseconds")
		}
		input.AfterMS = &afterMS
	}

	return input, nil
}
