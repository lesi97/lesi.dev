package trials_store

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

func (s *TrialsStore) GetPlayerCount() *string {
	trialsCh := make(chan struct {
		data *TrialsData
		err  error
	})
	steamCh := make(chan struct {
		data *SteamData
		err  error
	})

	go func() {
		data, err := s.fetchFromTrialsReport()
		trialsCh <- struct {
			data *TrialsData
			err  error
		}{data, err}
	}()

	go func() {
		data, err := s.fetchFromSteam()
		steamCh <- struct {
			data *SteamData
			err  error
		}{data, err}
	}()

	trialsResult := <-trialsCh
	steamResult := <-steamCh

	if trialsResult.err != nil {
		s.Logger.Printf("ERROR: fetchFromTrialsReport: %v\n", trialsResult.err)
		message := "failed to fetch data from trials report"
		return &message
	}

	if steamResult.err != nil {
		s.Logger.Printf("ERROR: fetchFromSteam: %v\n", steamResult.err)
		message := "failed to fetch data from steam"
		return &message
	}

	steamPlayerCount := humanize.Comma(int64(steamResult.data.Response.PlayerCount))
	trialsPlayerCount := humanize.Comma(int64(trialsResult.data.Platforms.Num0.RecentStats.PlayerCount))

	layout := "2006-01-02 15:04:05"
	updatedAtTime, err := time.Parse(layout, trialsResult.data.Platforms.Num0.RecentStats.UpdatedAt)
	if err != nil {
		s.Logger.Printf("ERROR: parsing updatedAt time: %v\n", err)
		message := "failed to parse time, tell lesi"
		return &message
	}

	if time.Since(updatedAtTime) >= 90*time.Minute {
		message := fmt.Sprintf(
			"There are currently %v players playing Destiny 2 on Steam",
			steamPlayerCount,
		)
		return &message
	}

	elapsedMinutes := int64(time.Since(updatedAtTime).Minutes())
	minuteLabel := "minutes"
	if elapsedMinutes == 1 {
		minuteLabel = "minute"
	}

	message := fmt.Sprintf(
		"There are currently %s players playing on Steam & %s players in Trials of Osiris across all platforms | Trials data last updated: %s %s ago from https://trials.report",
		steamPlayerCount,
		trialsPlayerCount,
		humanize.Comma(elapsedMinutes),
		minuteLabel,
	)

	return &message
}