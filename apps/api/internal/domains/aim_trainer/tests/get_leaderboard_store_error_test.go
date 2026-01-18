package tests

import (
	"context"
	"fmt"
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

type aimTrainerStoreErrorStub struct{}

func (s *aimTrainerStoreErrorStub) GetLeaderboard(ctx context.Context) ([]model.LeaderboardRow, error) {
	return nil, fmt.Errorf("boom")
}

func (s *aimTrainerStoreErrorStub) GetUser(ctx context.Context, username string) (*store.AimTrainerRow, error) {
	return nil, nil
}

func (s *aimTrainerStoreErrorStub) UpsertUpdateUser(ctx context.Context, in model.UpdateInput) (*model.AimTrainerUpdate, error) {
	return nil, nil
}

func TestGetLeaderboardStoreError(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := aim_handler.NewHandlerWithStore(logger, &aimTrainerStoreErrorStub{})

	router := chi.NewRouter()
	h.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/aim-trainer", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d got %d", http.StatusInternalServerError, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "boom") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
