package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	auth_store "github.com/lesi97/lesi.dev/internal/domains/auth/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleTwitchFrontendAuthStartRedirects(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := newAuthHandler(logger, &authStoreStub{
		twitchFrontendAuthStartResult: &auth_store.TwitchFrontendAuthStartResult{
			URL:          "https://example.com/oauth",
			State:        "state",
			PKCEVerifier: "pkce",
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/auth/twitch", nil)
	rec := httptest.NewRecorder()

	h.HandleTwitchFrontendAuthStart(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status %d got %d", http.StatusFound, rec.Code)
	}
	if rec.Header().Get("Location") != "https://example.com/oauth" {
		t.Fatalf("unexpected location: %s", rec.Header().Get("Location"))
	}
	cookies := rec.Header().Values("Set-Cookie")
	if len(cookies) < 2 {
		t.Fatalf("expected cookies to be set")
	}
}
