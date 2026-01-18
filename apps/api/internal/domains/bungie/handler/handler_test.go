package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleGetPlayTimeMissingID(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
	}

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

func TestHandleGetPrimaryMissingID(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/d2//primary", nil)
	rec := httptest.NewRecorder()

	h.HandleGetPrimary(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleGetSecondaryMissingID(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/d2//secondary", nil)
	rec := httptest.NewRecorder()

	h.HandleGetSecondary(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleGetHeavyMissingID(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/d2//heavy", nil)
	rec := httptest.NewRecorder()

	h.HandleGetHeavy(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleGetTerrorKillCountMissingName(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/d2/terror/weapons", nil)
	rec := httptest.NewRecorder()

	h.HandleGetTerrorKillCount(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleGetYungerKillCountMissingName(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/d2/yunger/weapons", nil)
	rec := httptest.NewRecorder()

	h.HandleGetYungerKillCount(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleResetReturnsOK(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/d2/reset", nil)
	rec := httptest.NewRecorder()

	h.HandleReset(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
}
