package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleAnilistAuthCallbackSuccess(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := newAuthHandler(logger, &authStoreStub{})

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
