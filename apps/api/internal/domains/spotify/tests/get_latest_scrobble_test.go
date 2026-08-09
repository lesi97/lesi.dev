package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	spotify_handler "github.com/lesi97/lesi.dev/internal/domains/spotify/handler"
	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type latestScrobbleStoreStub struct {
	text *string
	err  error
}

func (s *latestScrobbleStoreStub) InsertScrobble(ctx context.Context, input model.ScrobbleInput) (*model.ScrobbleResult, error) {
	return nil, nil
}

func (s *latestScrobbleStoreStub) GetLatestPlayedText(ctx context.Context) (*string, error) {
	return s.text, s.err
}

func (s *latestScrobbleStoreStub) PollSpotifyRecentlyPlayed(ctx context.Context, input model.SpotifyPollInput) (*model.SpotifyPollResult, error) {
	return nil, nil
}

func (s *latestScrobbleStoreStub) EnrichSpotifyMetadata(ctx context.Context, input model.SpotifyEnrichmentInput) (*model.SpotifyEnrichmentResult, error) {
	return nil, nil
}

func TestGetLatestScrobbleReturnsPlainTextWithoutApiKey(t *testing.T) {
	text := "Hozier - Talk"
	store := &latestScrobbleStoreStub{text: &text}

	router := chi.NewRouter()
	handler := spotify_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/spotify/latest", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d body %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if got, want := rec.Body.String(), text; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
	if got, want := rec.Header().Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Fatalf("content type = %q, want %q", got, want)
	}
}

func TestGetLatestScrobbleReturnsNotFoundWhenEmpty(t *testing.T) {
	store := &latestScrobbleStoreStub{}

	router := chi.NewRouter()
	handler := spotify_handler.NewHandlerWithStore(utils.NewColourLogger("brightBlack"), store)
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/spotify/latest", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d got %d body %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}
