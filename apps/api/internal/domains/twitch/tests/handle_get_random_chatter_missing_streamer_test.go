package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	twitch_handler "github.com/lesi97/lesi.dev/internal/domains/twitch/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleGetRandomChatterMissingStreamer(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := twitch_handler.NewHandlerWithStore(logger, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/twitch//chatters", nil)
	rec := httptest.NewRecorder()

	h.HandleGetRandomChatter(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "You must declare a streamer") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
