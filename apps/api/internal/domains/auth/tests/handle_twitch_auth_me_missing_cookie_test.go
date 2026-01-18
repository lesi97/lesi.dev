package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleTwitchAuthMeMissingCookie(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := newAuthHandler(logger, &authStoreStub{})

	req := httptest.NewRequest(http.MethodGet, "/auth/twitch/me", nil)
	rec := httptest.NewRecorder()

	h.HandleTwitchAuthMe(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d got %d", http.StatusUnauthorized, rec.Code)
	}
}
