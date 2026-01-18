package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	countdown_handler "github.com/lesi97/lesi.dev/internal/domains/countdown/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestCreateCountdownInvalidJSON(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := countdown_handler.NewHandlerWithStore(logger, nil)

	router := chi.NewRouter()
	h.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/countdown", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid request body") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
