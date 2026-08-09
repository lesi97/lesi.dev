package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
	domain_utils "github.com/lesi97/lesi.dev/internal/domains/spotify/internal/utils"
	shared_utils "github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) PostSpotifyEnrich(w http.ResponseWriter, r *http.Request) {
	if err := domain_utils.ValidateScrobbleAccess(r); err != nil {
		h.logger.Printf("Spotify enrichment access denied: %v", err)
		shared_utils.Error(w, http.StatusForbidden, "scrobble access denied")
		return
	}

	input, err := spotifyEnrichmentInputFromRequest(r)
	if err != nil {
		shared_utils.Error(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.store.EnrichSpotifyMetadata(r.Context(), input)
	if err != nil {
		h.logger.Printf("ERROR: Spotify enrichment: %v", err)
		if os.Getenv("GO_ENV") == "production" {
			shared_utils.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}
		shared_utils.Error(w, http.StatusInternalServerError, err)
		return
	}

	shared_utils.Success(w, http.StatusOK, result)
}

func spotifyEnrichmentInputFromRequest(r *http.Request) (model.SpotifyEnrichmentInput, error) {
	entityType := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("type")))
	if entityType == "" {
		entityType = strings.ToLower(strings.TrimSpace(r.URL.Query().Get("entity_type")))
	}
	if entityType == "" {
		entityType = model.SpotifyEnrichmentTypeTrack
	}

	switch entityType {
	case model.SpotifyEnrichmentTypeTrack, model.SpotifyEnrichmentTypeAlbum, model.SpotifyEnrichmentTypeArtist:
		return model.SpotifyEnrichmentInput{EntityType: entityType}, nil
	default:
		return model.SpotifyEnrichmentInput{}, fmt.Errorf("type must be track, album, or artist")
	}
}
