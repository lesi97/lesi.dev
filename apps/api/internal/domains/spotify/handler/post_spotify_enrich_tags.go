package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
	domain_utils "github.com/lesi97/lesi.dev/internal/domains/spotify/internal/utils"
	shared_utils "github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) PostSpotifyEnrichTags(w http.ResponseWriter, r *http.Request) {
	if err := domain_utils.ValidateScrobbleAccess(r); err != nil {
		h.logger.Printf("Last.fm tag enrichment access denied: %v", err)
		shared_utils.Error(w, http.StatusForbidden, "scrobble access denied")
		return
	}

	input, err := lastFMTagEnrichmentInputFromRequest(r)
	if err != nil {
		shared_utils.Error(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.store.EnrichLastFMTags(r.Context(), input)
	if err != nil {
		h.logger.Printf("ERROR: Last.fm tag enrichment: %v", err)
		if os.Getenv("GO_ENV") == "production" {
			shared_utils.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}
		shared_utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	shared_utils.Success(w, http.StatusOK, result)
}

func lastFMTagEnrichmentInputFromRequest(r *http.Request) (model.LastFMTagEnrichmentInput, error) {
	query := r.URL.Query()
	entityType := strings.ToLower(strings.TrimSpace(query.Get("type")))
	if entityType == "" {
		entityType = strings.ToLower(strings.TrimSpace(query.Get("entity_type")))
	}
	if entityType == "" {
		entityType = model.SpotifyEnrichmentTypeTrack
	}

	switch entityType {
	case model.SpotifyEnrichmentTypeTrack, model.SpotifyEnrichmentTypeAlbum, model.SpotifyEnrichmentTypeArtist:
		input := model.LastFMTagEnrichmentInput{
			EntityType: entityType,
			Force:      queryBool(query.Get("force")),
		}

		rawID := strings.TrimSpace(query.Get("id"))
		if rawID == "" {
			rawID = strings.TrimSpace(query.Get("entity_id"))
		}
		if rawID != "" {
			entityID, err := strconv.ParseInt(rawID, 10, 64)
			if err != nil || entityID < 1 {
				return model.LastFMTagEnrichmentInput{}, fmt.Errorf("id must be a positive integer")
			}
			input.EntityID = &entityID
			input.Force = true
		}

		return input, nil
	default:
		return model.LastFMTagEnrichmentInput{}, fmt.Errorf("type must be track, album, or artist")
	}
}

func queryBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
