package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	time_handler "github.com/lesi97/lesi.dev/internal/domains/time/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestGetDateTimeReturnsDateAndTime(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h, err := time_handler.NewHandler(logger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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
