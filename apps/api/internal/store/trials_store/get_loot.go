package trials_store

import (
	"fmt"
	"time"
)

func (s *TrialsStore) GetLoot() *string {

	trialsData, err := fetchFromTrialsReport()
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

	message := fmt.Sprintf("Map: %s | Flawless Loot: %s & Random 3, 5 or 7 win drop | Adept Mode: Random | Chance at Ship, Sparrow & Ghost", trialsData.Maps[0].Name, trialsData.Rewards.Flawless)

	return &message
}