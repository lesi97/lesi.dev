package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	aim_handler "github.com/lesi97/lesi.dev/internal/domains/aim_trainer/handler"
	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/internal/model"
	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/internal/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type aimTrainerStoreStub struct{}

func (s *aimTrainerStoreStub) GetLeaderboard(ctx context.Context) ([]model.LeaderboardRow, error) {
	return nil, nil
}

func (s *aimTrainerStoreStub) GetUser(ctx context.Context, username string) (*store.AimTrainerRow, error) {
	return nil, nil
}

func (s *aimTrainerStoreStub) UpsertUpdateUser(ctx context.Context, in model.UpdateInput) (*model.AimTrainerUpdate, error) {
	return nil, nil
}

func TestUpdateLeaderboardCORSRejected(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := aim_handler.NewHandlerWithStore(logger, &aimTrainerStoreStub{})

	router := chi.NewRouter()
	h.RegisterRoutes(router)

	t.Setenv("WEB_URL", "https://example.com")

	req := httptest.NewRequest(http.MethodPost, "/aim-trainer", strings.NewReader(`{}`))
	req.Header.Set("Origin", "https://bad.example")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d got %d", http.StatusForbidden, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Rejected due to CORS policy") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
