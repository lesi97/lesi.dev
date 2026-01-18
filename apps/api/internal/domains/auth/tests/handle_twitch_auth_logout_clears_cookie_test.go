package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleTwitchAuthLogoutClearsCookie(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := newAuthHandler(logger, &authStoreStub{})

	req := httptest.NewRequest(http.MethodPost, "/auth/twitch/logout", nil)
	rec := httptest.NewRecorder()

	h.HandleTwitchAuthLogout(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(strings.Join(rec.Header().Values("Set-Cookie"), ","), "lesidev_session=") {
		t.Fatalf("expected session cookie to be cleared")
	}
}
