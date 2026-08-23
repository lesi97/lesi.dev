package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	telemetry_handler "github.com/lesi97/lesi.dev/internal/domains/telemetry/handler"
	"github.com/lesi97/lesi.dev/internal/domains/telemetry/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type telemetryStoreStub struct {
	dateMemeClickInput *model.DateMemeClickInput
}

func (s *telemetryStoreStub) InsertTelemetry(_ context.Context, _ model.TelemetryPayload) error {
	return nil
}

func (s *telemetryStoreStub) UpsertDateMemeClick(_ context.Context, input model.DateMemeClickInput) error {
	s.dateMemeClickInput = &input
	return nil
}

func newTelemetryTestRouter(store *telemetryStoreStub) http.Handler {
	logger := utils.NewColourLogger("brightBlack")
	h := telemetry_handler.NewHandlerWithStore(logger, store)

	router := chi.NewRouter()
	h.RegisterRoutes(router)

	return router
}

func TestPostDateMemeClickRecordsServerDerivedIP(t *testing.T) {
	t.Setenv("WEB_URL", "https://example.com")
	t.Setenv("TELEMETRY_API_KEY", "secret")

	store := &telemetryStoreStub{}
	router := newTelemetryTestRouter(store)

	req := httptest.NewRequest(http.MethodPost, "/telemetry/date-meme-click", strings.NewReader(`{"route":"/date/sarah","action":"yes"}`))
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("X-Telemetry-Key", "secret")
	req.Header.Set("X-Forwarded-For", "203.0.113.10, 10.0.0.1")
	req.Header.Set("User-Agent", "date-test-agent")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d got %d: %s", http.StatusNoContent, rec.Code, rec.Body.String())
	}
	if store.dateMemeClickInput == nil {
		t.Fatal("expected store to record date meme click")
	}
	if store.dateMemeClickInput.Route != "/date/sarah" {
		t.Fatalf("unexpected route: %s", store.dateMemeClickInput.Route)
	}
	if store.dateMemeClickInput.Action != model.DateMemeClickActionYes {
		t.Fatalf("unexpected action: %s", store.dateMemeClickInput.Action)
	}
	if store.dateMemeClickInput.SecretEnding {
		t.Fatal("secret ending should default to false")
	}
	if store.dateMemeClickInput.IP != "203.0.113.10" {
		t.Fatalf("unexpected ip: %s", store.dateMemeClickInput.IP)
	}
	if store.dateMemeClickInput.UserAgent != "date-test-agent" {
		t.Fatalf("unexpected user agent: %s", store.dateMemeClickInput.UserAgent)
	}
}

func TestPostDateMemeClickRecordsSecretEnding(t *testing.T) {
	t.Setenv("WEB_URL", "https://example.com")
	t.Setenv("TELEMETRY_API_KEY", "secret")

	store := &telemetryStoreStub{}
	router := newTelemetryTestRouter(store)

	req := httptest.NewRequest(http.MethodPost, "/telemetry/date-meme-click", strings.NewReader(`{"route":"/audrey","action":"yes","secretEnding":true}`))
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("X-Telemetry-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d got %d: %s", http.StatusNoContent, rec.Code, rec.Body.String())
	}
	if store.dateMemeClickInput == nil {
		t.Fatal("expected store to record date meme click")
	}
	if !store.dateMemeClickInput.SecretEnding {
		t.Fatal("expected secret ending to be recorded")
	}
}

func TestPostDateMemeClickRejectsInvalidAction(t *testing.T) {
	t.Setenv("WEB_URL", "https://example.com")
	t.Setenv("TELEMETRY_API_KEY", "secret")

	store := &telemetryStoreStub{}
	router := newTelemetryTestRouter(store)

	req := httptest.NewRequest(http.MethodPost, "/telemetry/date-meme-click", strings.NewReader(`{"route":"/audrey","action":"maybe"}`))
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("X-Telemetry-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
	if store.dateMemeClickInput != nil {
		t.Fatal("store should not be called for invalid action")
	}
	if !strings.Contains(rec.Body.String(), "action must be yes or no") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestPostDateMemeClickRejectsSecretEndingForNo(t *testing.T) {
	t.Setenv("WEB_URL", "https://example.com")
	t.Setenv("TELEMETRY_API_KEY", "secret")

	store := &telemetryStoreStub{}
	router := newTelemetryTestRouter(store)

	req := httptest.NewRequest(http.MethodPost, "/telemetry/date-meme-click", strings.NewReader(`{"route":"/audrey","action":"no","secretEnding":true}`))
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("X-Telemetry-Key", "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
	if store.dateMemeClickInput != nil {
		t.Fatal("store should not be called when secret ending is sent for no")
	}
	if !strings.Contains(rec.Body.String(), "secret ending can only be tracked for yes clicks") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
