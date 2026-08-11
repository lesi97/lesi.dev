package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	music_handler "github.com/lesi97/lesi.dev/internal/domains/spotify/handler"
	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type musicStoreStub struct {
	pollCalled bool
	pollInput  model.SpotifyPollInput
	pollResult *model.SpotifyPollResult
	pollErr    error

	enrichCalled bool
	enrichInput  model.SpotifyEnrichmentInput
	enrichResult *model.SpotifyEnrichmentResult
	enrichErr    error

	tagEnrichCalled bool
	tagEnrichInput  model.LastFMTagEnrichmentInput
	tagEnrichResult *model.LastFMTagEnrichmentResult
	tagEnrichErr    error
}

func (s *musicStoreStub) InsertScrobble(ctx context.Context, input model.ScrobbleInput) (*model.ScrobbleResult, error) {
	return nil, nil
}

func (s *musicStoreStub) GetLatestPlayedText(ctx context.Context) (*string, error) {
	return nil, nil
}

func (s *musicStoreStub) PollSpotifyRecentlyPlayed(ctx context.Context, input model.SpotifyPollInput) (*model.SpotifyPollResult, error) {
	s.pollCalled = true
	s.pollInput = input
	return s.pollResult, s.pollErr
}

func (s *musicStoreStub) EnrichSpotifyMetadata(ctx context.Context, input model.SpotifyEnrichmentInput) (*model.SpotifyEnrichmentResult, error) {
	s.enrichCalled = true
	s.enrichInput = input
	return s.enrichResult, s.enrichErr
}

func (s *musicStoreStub) EnrichLastFMTags(ctx context.Context, input model.LastFMTagEnrichmentInput) (*model.LastFMTagEnrichmentResult, error) {
	s.tagEnrichCalled = true
	s.tagEnrichInput = input
	return s.tagEnrichResult, s.tagEnrichErr
}

func TestPostSpotifyPollCallsStore(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	afterMS := int64(1786012626000)
	store := &musicStoreStub{
		pollResult: &model.SpotifyPollResult{
			Fetched:   2,
			Scrobbled: 2,
			Source:    "spotify-poll",
		},
	}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/poll?limit=5&after=1786012626000", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if !store.pollCalled {
		t.Fatal("expected PollSpotifyRecentlyPlayed to be called")
	}
	if got, want := store.pollInput.Limit, 5; got != want {
		t.Fatalf("limit = %d, want %d", got, want)
	}
	if store.pollInput.AfterMS == nil || *store.pollInput.AfterMS != afterMS {
		t.Fatalf("after ms = %v, want %d", store.pollInput.AfterMS, afterMS)
	}
}

func TestPostSpotifyPollReturnsRateLimitedResult(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	retryAfter := 120
	store := &musicStoreStub{
		pollResult: &model.SpotifyPollResult{
			Source:                     "spotify-poll",
			RateLimited:                true,
			RateLimitReason:            "spotify recently played quota exceeded",
			RateLimitRetryAfterSeconds: &retryAfter,
		},
	}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/poll", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"rate_limited":true`) {
		t.Fatalf("expected rate limited response, got %s", rec.Body.String())
	}
}

func TestPostSpotifyPollRejectsMissingApiKey(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/poll", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d got %d body %s", http.StatusForbidden, rec.Code, rec.Body.String())
	}
	if store.pollCalled {
		t.Fatal("poll should not be called when API key is missing")
	}
}

func TestPostSpotifyPollReturnsStoreErrorOutsideProduction(t *testing.T) {
	t.Setenv("GO_ENV", "development")
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{pollErr: errors.New("spotify token refresh failed")}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/poll", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d got %d body %s", http.StatusInternalServerError, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "spotify token refresh failed") {
		t.Fatalf("expected response to include store error, got %s", rec.Body.String())
	}
}

func TestPostSpotifyPollHidesStoreErrorInProduction(t *testing.T) {
	t.Setenv("GO_ENV", "production")
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{pollErr: errors.New("spotify token refresh failed")}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/poll", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d got %d body %s", http.StatusInternalServerError, rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "internal server error") {
		t.Fatalf("expected response to hide store error, got %s", rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "spotify token refresh failed") {
		t.Fatalf("expected response to hide store error, got %s", rec.Body.String())
	}
}

func TestPostSpotifyEnrichCallsStore(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{
		enrichResult: &model.SpotifyEnrichmentResult{
			EntityType: model.SpotifyEnrichmentTypeAlbum,
			Status:     "complete",
		},
	}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/enrich?type=album", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if !store.enrichCalled {
		t.Fatal("expected EnrichSpotifyMetadata to be called")
	}
	if got, want := store.enrichInput.EntityType, model.SpotifyEnrichmentTypeAlbum; got != want {
		t.Fatalf("entity type = %q, want %q", got, want)
	}
}

func TestPostSpotifyEnrichDefaultsToTrack(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{
		enrichResult: &model.SpotifyEnrichmentResult{
			EntityType: model.SpotifyEnrichmentTypeTrack,
			Status:     "complete",
		},
	}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/enrich", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if got, want := store.enrichInput.EntityType, model.SpotifyEnrichmentTypeTrack; got != want {
		t.Fatalf("entity type = %q, want %q", got, want)
	}
}

func TestPostSpotifyEnrichRejectsInvalidType(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/enrich?type=playlist", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d body %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if store.enrichCalled {
		t.Fatal("enrichment should not be called when type is invalid")
	}
}

func TestPostSpotifyEnrichRejectsMissingApiKey(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/enrich?type=track", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d got %d body %s", http.StatusForbidden, rec.Code, rec.Body.String())
	}
	if store.enrichCalled {
		t.Fatal("enrichment should not be called when API key is missing")
	}
}

func TestPostSpotifyEnrichTagsCallsStore(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{
		tagEnrichResult: &model.LastFMTagEnrichmentResult{
			EntityType: model.SpotifyEnrichmentTypeAlbum,
			Status:     "tagged",
			TagsFound:  3,
			Updated:    true,
		},
	}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/enrich/tags?type=album", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if !store.tagEnrichCalled {
		t.Fatal("expected EnrichLastFMTags to be called")
	}
	if got, want := store.tagEnrichInput.EntityType, model.SpotifyEnrichmentTypeAlbum; got != want {
		t.Fatalf("entity type = %q, want %q", got, want)
	}
}

func TestPostSpotifyEnrichTagsDefaultsToTrack(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{
		tagEnrichResult: &model.LastFMTagEnrichmentResult{
			EntityType: model.SpotifyEnrichmentTypeTrack,
			Status:     "complete",
		},
	}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/enrich/tags", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if got, want := store.tagEnrichInput.EntityType, model.SpotifyEnrichmentTypeTrack; got != want {
		t.Fatalf("entity type = %q, want %q", got, want)
	}
}

func TestPostSpotifyEnrichTagsAcceptsEntityID(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{
		tagEnrichResult: &model.LastFMTagEnrichmentResult{
			EntityType: model.SpotifyEnrichmentTypeTrack,
			Status:     "tagged",
		},
	}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/enrich/tags?type=track&id=121", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if store.tagEnrichInput.EntityID == nil || *store.tagEnrichInput.EntityID != 121 {
		t.Fatalf("entity id = %v, want 121", store.tagEnrichInput.EntityID)
	}
	if !store.tagEnrichInput.Force {
		t.Fatal("expected id lookup to force enrichment")
	}
}

func TestPostSpotifyEnrichTagsRejectsInvalidType(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/enrich/tags?type=playlist", nil)
	req.Header.Set("X-API-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d body %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
	if store.tagEnrichCalled {
		t.Fatal("tag enrichment should not be called when type is invalid")
	}
}

func TestPostSpotifyEnrichTagsRejectsMissingApiKey(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/spotify/enrich/tags?type=track", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d got %d body %s", http.StatusForbidden, rec.Code, rec.Body.String())
	}
	if store.tagEnrichCalled {
		t.Fatal("tag enrichment should not be called when API key is missing")
	}
}
