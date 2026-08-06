package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
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

	req := httptest.NewRequest(http.MethodPost, "/scrobbles/poll/spotify?limit=5&after=1786012626000", nil)
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

func TestPostSpotifyPollRejectsMissingApiKey(t *testing.T) {
	t.Setenv("SCROBBLE_API_KEY", "secret")
	store := &musicStoreStub{}

	router := chi.NewRouter()
	handler := music_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/scrobbles/poll/spotify", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d got %d body %s", http.StatusForbidden, rec.Code, rec.Body.String())
	}
	if store.pollCalled {
		t.Fatal("poll should not be called when API key is missing")
	}
}
