package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleAniAuthInitialRedirectMissingSecret(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := newAuthHandler(logger, &authStoreStub{})

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/anilist/login", nil)
	rec := httptest.NewRecorder()

	h.HandleAniAuthInitialRedirect(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d got %d", http.StatusForbidden, rec.Code)
	}
}
