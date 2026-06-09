package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	trials_handler "github.com/lesi97/lesi.dev/internal/domains/trials/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type trialsPlayerCountStoreMock struct{}

func (trialsPlayerCountStoreMock) GetLoot(_ context.Context) *string {
	message := "loot"
	return &message
}

func (trialsPlayerCountStoreMock) GetPlayerCount(_ context.Context) *string {
	message := "player count"
	return &message
}

func TestHandleGetPlayerCountSetsNoStoreHeaders(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := trials_handler.NewHandlerWithStore(logger, trialsPlayerCountStoreMock{})

	req := httptest.NewRequest(http.MethodGet, "/d2/trials/playercount", nil)
	rec := httptest.NewRecorder()

	h.HandleGetPlayerCount(rec, req)

	if rec.Header().Get("Cache-Control") != "no-store, max-age=0" {
		t.Fatalf("expected Cache-Control no-store, max-age=0 got %q", rec.Header().Get("Cache-Control"))
	}
	if rec.Header().Get("Pragma") != "no-cache" {
		t.Fatalf("expected Pragma no-cache got %q", rec.Header().Get("Pragma"))
	}
	if rec.Header().Get("Expires") != "0" {
		t.Fatalf("expected Expires 0 got %q", rec.Header().Get("Expires"))
	}
}
