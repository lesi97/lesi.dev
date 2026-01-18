package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestCreateCountdownInvalidJSON(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/countdown", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	h.CreateCountdown(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid request body") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestGetCountdownMissingID(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/countdown", nil)
	rec := httptest.NewRecorder()

	h.getCountdown(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid ID") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
