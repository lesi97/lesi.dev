package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	steam_handler "github.com/lesi97/lesi.dev/internal/domains/steam/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleGetPlayerCountMissingParams(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	expectedErr := "you must provide gameId or gameName"
	mock := &storeMock{
		err: errors.New(expectedErr),
	}

	h := steam_handler.NewHandlerWithStore(logger, mock)

	req := httptest.NewRequest(http.MethodGet, "/steam/playercount", nil)
	rec := httptest.NewRecorder()

	h.HandleGetPlayerCount(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), expectedErr) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
