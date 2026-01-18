package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	bungie_handler "github.com/lesi97/lesi.dev/internal/domains/bungie/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleGetTerrorKillCountMissingName(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := bungie_handler.NewHandlerWithStore(logger, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/d2/terror/weapons", nil)
	rec := httptest.NewRecorder()

	h.HandleGetTerrorKillCount(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
}
