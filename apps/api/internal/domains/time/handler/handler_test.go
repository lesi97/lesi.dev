package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestGetDateTimeReturnsDateAndTime(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := NewHandler(logger)

	req := httptest.NewRequest(http.MethodGet, "/time", nil)
	rec := httptest.NewRecorder()

	h.GetDateTime(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "\"date\"") || !strings.Contains(body, "\"time\"") {
		t.Fatalf("unexpected body: %s", body)
	}
}
