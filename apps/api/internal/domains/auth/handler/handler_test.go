package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	auth_store "github.com/lesi97/lesi.dev/internal/domains/auth/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type authStoreStub struct {
	anilistAuthURL                    *string
	anilistAuthErr                    error
	anilistCallbackErr                error
	twitchModAuthURL                  *string
	twitchModAuthErr                  error
	twitchModCallbackErr              error
	twitchFrontendAuthStartResult     *auth_store.TwitchFrontendAuthStartResult
	twitchFrontendAuthStartErr        error
	twitchFrontendCallbackIdentity    *auth_store.TwitchFrontendIdentity
	twitchFrontendCallbackErr         error
	twitchFrontendUpsertUserID        string
	twitchFrontendUpsertUserErr       error
	twitchFrontendCreateSessionToken  string
	twitchFrontendCreateSessionErr    error
	twitchFrontendGetUser             *auth_store.TwitchFrontendUser
	twitchFrontendGetUserErr          error
	twitchFrontendDeleteSessionErr    error
}

func (s *authStoreStub) AnilistAuthUrl() (*string, error) {
	return s.anilistAuthURL, s.anilistAuthErr
}

func (s *authStoreStub) AnilistCallback(code string) error {
	return s.anilistCallbackErr
}

func (s *authStoreStub) TwitchModAuthUrl() (*string, error) {
	return s.twitchModAuthURL, s.twitchModAuthErr
}

func (s *authStoreStub) TwitchModCallback(code string) error {
	return s.twitchModCallbackErr
}

func (s *authStoreStub) TwitchFrontendAuthStart() (*auth_store.TwitchFrontendAuthStartResult, error) {
	return s.twitchFrontendAuthStartResult, s.twitchFrontendAuthStartErr
}

func (s *authStoreStub) TwitchFrontendCallback(code string, state string, expectedState string, pkceVerifier string) (*auth_store.TwitchFrontendIdentity, error) {
	return s.twitchFrontendCallbackIdentity, s.twitchFrontendCallbackErr
}

func (s *authStoreStub) TwitchFrontendUpsertUser(ctx context.Context, identity auth_store.TwitchFrontendIdentity) (string, error) {
	return s.twitchFrontendUpsertUserID, s.twitchFrontendUpsertUserErr
}

func (s *authStoreStub) TwitchFrontendCreateSession(ctx context.Context, userID string, ttl time.Duration) (string, error) {
	return s.twitchFrontendCreateSessionToken, s.twitchFrontendCreateSessionErr
}

func (s *authStoreStub) TwitchFrontendGetUserBySession(ctx context.Context, sessionToken string) (*auth_store.TwitchFrontendUser, error) {
	return s.twitchFrontendGetUser, s.twitchFrontendGetUserErr
}

func (s *authStoreStub) TwitchFrontendDeleteSessionByToken(ctx context.Context, sessionToken string) error {
	return s.twitchFrontendDeleteSessionErr
}

func TestHandleAniAuthInitialRedirectMissingSecret(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
		store:  &authStoreStub{},
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/anilist/login", nil)
	rec := httptest.NewRecorder()

	h.HandleAniAuthInitialRedirect(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d got %d", http.StatusForbidden, rec.Code)
	}
}

func TestHandleAnilistAuthCallbackSuccess(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
		store:  &authStoreStub{},
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/anilist/callback?code=test", nil)
	rec := httptest.NewRecorder()

	h.HandleAnilistAuthCallback(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "updated") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandleTwitchModAuthInitialRedirectMissingSecret(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
		store:  &authStoreStub{},
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/twitch/login", nil)
	rec := httptest.NewRecorder()

	h.HandleTwitchModAuthInitialRedirect(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d got %d", http.StatusForbidden, rec.Code)
	}
}

func TestHandleTwitchModAuthCallbackSuccess(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
		store:  &authStoreStub{},
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/twitch/callback?code=test", nil)
	rec := httptest.NewRecorder()

	h.HandleTwitchModAuthCallback(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "updated") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandleTwitchFrontendAuthStartRedirects(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
		store: &authStoreStub{
			twitchFrontendAuthStartResult: &auth_store.TwitchFrontendAuthStartResult{
				URL:          "https://example.com/oauth",
				State:        "state",
				PKCEVerifier: "pkce",
			},
		},
	}

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

func TestHandleTwitchFrontendAuthCallbackMissingCode(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
		store:  &authStoreStub{},
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/twitch/callback?state=state", nil)
	rec := httptest.NewRecorder()

	h.HandleTwitchFrontendAuthCallback(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandleTwitchAuthMeMissingCookie(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
		store:  &authStoreStub{},
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/twitch/me", nil)
	rec := httptest.NewRecorder()

	h.HandleTwitchAuthMe(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestHandleTwitchAuthLogoutClearsCookie(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
		store:  &authStoreStub{},
	}

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
