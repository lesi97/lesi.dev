package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	bungie_handler "github.com/lesi97/lesi.dev/internal/domains/bungie/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleGetPlayTimeMissingID(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := bungie_handler.NewHandlerWithStore(logger, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/d2//time", nil)
	rec := httptest.NewRecorder()

	h.HandleGetPlayTime(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid ID") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
