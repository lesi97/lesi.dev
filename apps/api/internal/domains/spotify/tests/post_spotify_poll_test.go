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
