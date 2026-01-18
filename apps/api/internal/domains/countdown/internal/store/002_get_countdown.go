package store

import (
	"context"
	"fmt"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/countdown/internal/model"
	cu "github.com/lesi97/lesi.dev/internal/domains/countdown/internal/utils"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (s *Store) GetCountdownByID(ctx context.Context, uuid string) (*string, error) {
	defer s.Logger.LogExecutionTime("DATABASE CALL: getCountdownById", time.Now(), nil)

	data := &model.FetchData{
		UUID: &uuid,
	}
	err := data.Select(s.DB, &ctx, &uuid)

	if err != nil {
		return nil, err
	}

	result := cu.CountdownToString(data.TargetDate)
	if result == "Passed" {
		message := data.FallbackMessage
		utils.TruncateString(&message, 400)
		return &message, nil
	}

	message := fmt.Sprintf("%s %s", result, data.Message)
	utils.TruncateString(&message, 400)

	return &message, nil
}

