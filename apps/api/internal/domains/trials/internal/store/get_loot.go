package store

import (
	"context"
	"fmt"
	"time"

	trials_utils "github.com/lesi97/lesi.dev/internal/domains/trials/internal/utils"
)

func (s *Store) GetLoot(ctx context.Context) *string {
	trialsData, err := trials_utils.FetchFromTrialsReport(ctx, s.Logger, s.URL, s.Redis)
	if err != nil {
		s.Logger.Printf("ERROR: fetchFromTrialsReport: %v\n", err)
		message := "failed to fetch data from trials report"
		return &message
	}

	layout := "2006-01-02 15:04:05"
	endTime, err := time.Parse(layout, trialsData.EndDate)
	if err != nil {
		s.Logger.Printf("ERROR: parsing EndDate: %v\n", err)
		message := "could not parse trial end time"
		return &message
	}

	if endTime.Before(time.Now()) {
		message := "trials isn't here yet dummy, gift a sub and come back later"
		return &message
	}

	message := fmt.Sprintf(
		"Map: %s | Flawless Loot: %s & Random 3, 5 or 7 win drop | Adept Mode: Random | Chance at Ship, Sparrow & Ghost",
		trialsData.Maps[0].Name,
		trialsData.Rewards.Flawless,
	)

	return &message
}
