package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/model"
	"github.com/lesi97/lesi.dev/internal/domains/aim_trainer/store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type aimTrainerStoreStub struct {
	getLeaderboardRows []model.LeaderboardRow
	getLeaderboardErr  error
}

func (s *aimTrainerStoreStub) GetLeaderboard(ctx context.Context) ([]model.LeaderboardRow, error) {
	return s.getLeaderboardRows, s.getLeaderboardErr
}

func (s *aimTrainerStoreStub) GetUser(ctx context.Context, username string) (*store.AimTrainerRow, error) {
	return nil, nil
}

func (s *aimTrainerStoreStub) UpsertUpdateUser(ctx context.Context, in model.UpdateInput) (*model.AimTrainerUpdate, error) {
	return nil, nil
}

func TestGetLeaderboardStoreError(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
		store: &aimTrainerStoreStub{
			getLeaderboardErr: fmt.Errorf("boom"),
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/aim-trainer", nil)
	rec := httptest.NewRecorder()

	h.getLeaderboard(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d got %d", http.StatusInternalServerError, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "boom") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestUpdateLeaderboardCORSRejected(t *testing.T) {
	logger := utils.NewColourLogger("brightBlack")
	h := &Handler{
		logger: logger,
		store:  &aimTrainerStoreStub{},
	}

	t.Setenv("WEB_URL", "https://example.com")

	req := httptest.NewRequest(http.MethodPost, "/aim-trainer", strings.NewReader(`{}`))
	req.Header.Set("Origin", "https://bad.example")
	rec := httptest.NewRecorder()

	h.updateLeaderboard(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d got %d", http.StatusForbidden, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Rejected due to CORS policy") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
