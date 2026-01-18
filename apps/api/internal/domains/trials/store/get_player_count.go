package store

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/lesi97/lesi.dev/internal/domains/trials/model"
	trials_utils "github.com/lesi97/lesi.dev/internal/domains/trials/utils"
	core_utils "github.com/lesi97/lesi.dev/internal/utils"
)

func (s *Store) GetPlayerCount() *string {
	trialsCh := make(chan struct {
		data *model.TrialsData
		err  error
	}, 1)

	steamCh := make(chan struct {
		data *model.SteamData
		err  error
	}, 1)

	go func() {
		data, err := trials_utils.FetchFromTrialsReport(s.Logger, s.URL)
		trialsCh <- struct {
			data *model.TrialsData
			err  error
		}{data, err}
	}()

	if s.SteamClientIDAvailable {
		go func() {
			data, err := trials_utils.FetchFromSteam(s.Logger, s.SteamURL, s.SteamClientID)
			steamCh <- struct {
				data *model.SteamData
				err  error
			}{data, err}
		}()
	}

	trialsResult := <-trialsCh
	if trialsResult.err != nil {
		s.Logger.Printf("ERROR: fetchFromTrialsReport: %v\n", trialsResult.err)
		message := "failed to fetch data from trials report"
		return &message
	}

	updatedAtTime, _ := trials_utils.ParseTrialsUpdatedAt(trialsResult.data)
	trialsFresh := trials_utils.IsTrialsFresh(updatedAtTime)

	var steamResult struct {
		data *model.SteamData
		err  error
	}
	steamAvailable := false

	if s.SteamClientIDAvailable {
		steamResult = <-steamCh
		if steamResult.err != nil {
			err := fmt.Sprintf("ERROR: fetchFromSteam: %v\n", steamResult.err)
			s.Logger.Print(err)
			s.Logger.SendDiscordNotification(core_utils.SendDiscordNotificationArgs{
				Content:  err,
				Username: "Trials Store",
				Title:    "Trials Store Error",
			})
			message := "failed to fetch data from steam"
			return &message
		}
		steamAvailable = true
	}

	if !steamAvailable && !trialsFresh {
		message := "no data available"
		return &message
	}

	if steamAvailable && !trialsFresh {
		steamPlayerCount := humanize.Comma(int64(steamResult.data.Response.PlayerCount))
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

	trialsPlayerCount := humanize.Comma(int64(trialsResult.data.Platforms.Num0.RecentStats.PlayerCount))

	if !steamAvailable {
		message := fmt.Sprintf(
			"There are currently %s players in Trials of Osiris across all platforms | Trials data last updated: %s %s ago from https://trials.report",
			trialsPlayerCount,
			humanize.Comma(elapsedMinutes),
			minuteLabel,
		)
		return &message
	}

	steamPlayerCount := humanize.Comma(int64(steamResult.data.Response.PlayerCount))
	message := fmt.Sprintf(
		"There are currently %s players playing on Steam & %s players in Trials of Osiris across all platforms | Trials data last updated: %s %s ago from https://trials.report",
		steamPlayerCount,
		trialsPlayerCount,
		humanize.Comma(elapsedMinutes),
		minuteLabel,
	)

	return &message
}
