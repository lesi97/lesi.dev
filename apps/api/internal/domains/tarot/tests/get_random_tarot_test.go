package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	tarot_handler "github.com/lesi97/lesi.dev/internal/domains/tarot/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestGetRandomTarotReturnsMessage(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h, err := tarot_handler.NewHandler(logger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	router := chi.NewRouter()
	h.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/tarot", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) == "" {
		t.Fatalf("expected non-empty body")
	}
}
