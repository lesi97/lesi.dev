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

	trials_utils.ClearTrialsReportCache()

	if !trials_utils.IsTrialsReportAvailable(time.Now()) {
		t.Skip("trials report is not available")
	}

	logger := utils.NewColourLogger("brightBlack")
	s := &trials_store.Store{
		Logger:        logger,
		URL:           server.URL,
		SteamClientID: "",
	}
	h := trials_handler.NewHandlerWithStore(logger, s)

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
