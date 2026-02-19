package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	steam_handler "github.com/lesi97/lesi.dev/internal/domains/steam/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleGetPlayerCountWithGameID(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	message := "There are currently 10 players on Steam for app 730"
	mock := &storeMock{
		message: &message,
	}

	h := steam_handler.NewHandlerWithStore(logger, mock)

	req := httptest.NewRequest(http.MethodGet, "/steam/playercount?gameId=730", nil)
	rec := httptest.NewRecorder()

	h.HandleGetPlayerCount(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), message) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}

	if mock.receivedID != "730" {
		t.Fatalf("expected gameId to be 730 got %s", mock.receivedID)
	}

	if mock.receivedName != "" {
		t.Fatalf("expected empty gameName got %s", mock.receivedName)
	}
}
