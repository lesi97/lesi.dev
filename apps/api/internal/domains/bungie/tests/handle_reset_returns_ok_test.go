package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	bungie_handler "github.com/lesi97/lesi.dev/internal/domains/bungie/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleResetReturnsOK(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := bungie_handler.NewHandlerWithStore(logger, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/d2/reset", nil)
	rec := httptest.NewRecorder()

	h.HandleReset(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
}
