package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/trials/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleGetLootReturnsMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"endDate": "2999-01-01 00:00:00",
			"maps": [{"name": "TestMap"}],
			"rewards": {"flawless": "TestLoot"}
		}`))
	}))
	defer server.Close()

	logger := utils.NewColourLogger("brightBlack")
	s := &store.Store{
		Logger:                 logger,
		URL:                    server.URL,
		SteamClientIDAvailable: false,
	}
	h := &Handler{
		logger: logger,
		store:  s,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/trials/loot", nil)
	rec := httptest.NewRecorder()

	h.HandleGetLoot(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "TestMap") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandleGetPlayerCountUsesTrialsReport(t *testing.T) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"platforms": {"0": {"recentStats": {"playerCount": 1234, "updatedAt": "` + now + `"}}}
		}`))
	}))
	defer server.Close()

	logger := utils.NewColourLogger("brightBlack")
	s := &store.Store{
		Logger:                 logger,
		URL:                    server.URL,
		SteamClientIDAvailable: false,
	}
	h := &Handler{
		logger: logger,
		store:  s,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/trials/player-count", nil)
	rec := httptest.NewRecorder()

	h.HandleGetPlayerCount(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Trials of Osiris") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
