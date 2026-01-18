package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestGetRandomTarotReturnsMessage(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := NewHandler(logger)

	req := httptest.NewRequest(http.MethodGet, "/tarot/random", nil)
	rec := httptest.NewRecorder()

	h.getRandomTarot(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) == "" {
		t.Fatalf("expected non-empty body")
	}
}

func TestGetAllCardsReturnsJSON(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := NewHandler(logger)

	req := httptest.NewRequest(http.MethodGet, "/tarot/all", nil)
	rec := httptest.NewRecorder()

	h.getAllCards(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "\"cards\"") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
