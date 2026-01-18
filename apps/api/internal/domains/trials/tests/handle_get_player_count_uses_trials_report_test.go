package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	trials_handler "github.com/lesi97/lesi.dev/internal/domains/trials/handler"
	trials_store "github.com/lesi97/lesi.dev/internal/domains/trials/internal/store"
	trials_utils "github.com/lesi97/lesi.dev/internal/domains/trials/internal/utils"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleGetPlayerCountUsesTrialsReport(t *testing.T) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"platforms": {"0": {"recentStats": {"playerCount": 1234, "updatedAt": "` + now + `"}}}
		}`))
	}))
	defer server.Close()

	trials_utils.ClearTrialsReportCache()

	if !trials_utils.IsTrialsReportAvailable(time.Now()) {
		t.Skip("trials report is not available")
	}

	logger := utils.NewColourLogger("brightBlack")
	s := &trials_store.Store{
		Logger:        logger,
		URL:           server.URL,
		SteamClientID: "test",
		SteamURL:      server.URL,
	}
	h := trials_handler.NewHandlerWithStore(logger, s)

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
