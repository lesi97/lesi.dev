package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleTwitchFrontendAuthCallbackMissingCode(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := newAuthHandler(logger, &authStoreStub{})

	req := httptest.NewRequest(http.MethodGet, "/auth/twitch/callback?state=state", nil)
	rec := httptest.NewRecorder()

	h.HandleTwitchFrontendAuthCallback(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
}
