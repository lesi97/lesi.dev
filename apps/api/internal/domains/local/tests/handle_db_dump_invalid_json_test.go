package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	local_handler "github.com/lesi97/lesi.dev/internal/domains/local/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleDbDumpInvalidJSON(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := local_handler.NewHandlerWithStore(logger, nil)

	req := httptest.NewRequest(http.MethodPost, "/local/db-dump", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	h.HandleDbDump(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid json body") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
