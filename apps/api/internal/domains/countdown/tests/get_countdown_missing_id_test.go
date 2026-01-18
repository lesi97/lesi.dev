package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	countdown_handler "github.com/lesi97/lesi.dev/internal/domains/countdown/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestGetCountdownMissingID(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := countdown_handler.NewHandlerWithStore(logger, nil)

	req := httptest.NewRequest(http.MethodGet, "/countdown", nil)
	rec := httptest.NewRecorder()

	h.HandleGetCountdown(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid ID") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
