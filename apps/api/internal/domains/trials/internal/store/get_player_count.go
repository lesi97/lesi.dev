package store

import (
	"context"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/lesi97/lesi.dev/internal/domains/trials/internal/model"
	trials_utils "github.com/lesi97/lesi.dev/internal/domains/trials/internal/utils"
	core_utils "github.com/lesi97/lesi.dev/internal/utils"
)

func (s *Store) GetPlayerCount(ctx context.Context) *string {
	trialsAvailable := trials_utils.IsTrialsReportAvailable(time.Now())

	trialsCh := make(chan struct {
		data *model.TrialsData
		err  error
	}, 1)

	steamCh := make(chan struct {
		data *model.SteamData
		err  error
	}, 1)

	if trialsAvailable {
		go func() {
			data, err := trials_utils.FetchFromTrialsReport(context.Background(), s.Logger, s.URL, s.Redis)
			trialsCh <- struct {
				data *model.TrialsData
				err  error
			}{data, err}
		}()
	}

	go func() {
		data, err := trials_utils.FetchFromSteam(ctx, s.Logger, s.SteamURL, s.SteamClientID)
		steamCh <- struct {
			data *model.SteamData
			err  error
		}{data, err}
	}()

	var steamResult struct {
		data *model.SteamData
		err  error
	}
	steamAvailable := false

	steamResult = <-steamCh
	if steamResult.err != nil {
		err := fmt.Sprintf("ERROR: fetchFromSteam: %v\n", steamResult.err)
		s.Logger.Print(err)
		s.Logger.SendDiscordNotification(core_utils.SendDiscordNotificationArgs{
			Content:  err,
			Username: "Trials Store",
			Title:    "Trials Store Error",
		})
	} else {
		steamAvailable = true
	}

	var trialsResult struct {
		data *model.TrialsData
		err  error
	}
	trialsFresh := false
	var updatedAtTime time.Time

	if trialsAvailable {
		select {
		case trialsResult = <-trialsCh:
		case <-time.After(2 * time.Second):
			trialsResult.err = fmt.Errorf("trials report request timed out")
		}

		if trialsResult.err == nil {
			updatedAtTime, _ = trials_utils.ParseTrialsUpdatedAt(trialsResult.data)
			trialsFresh = trials_utils.IsTrialsFresh(updatedAtTime)
		}
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

	if !steamAvailable {
		elapsedMinutes := int64(time.Since(updatedAtTime).Minutes())
		minuteLabel := "minutes"
		if elapsedMinutes == 1 {
			minuteLabel = "minute"
		}
		trialsPlayerCount := humanize.Comma(int64(trialsResult.data.Platforms.Num0.RecentStats.PlayerCount))
		message := fmt.Sprintf(
			"There are currently %s players in trials across all platforms | Trials data last updated: %s %s ago from https://trials.report",
			trialsPlayerCount,
			humanize.Comma(elapsedMinutes),
			minuteLabel,
		)
		return &message
	}

	elapsedMinutes := int64(time.Since(updatedAtTime).Minutes())
	minuteLabel := "minutes"
	if elapsedMinutes == 1 {
		minuteLabel = "minute"
	}
	trialsPlayerCount := humanize.Comma(int64(trialsResult.data.Platforms.Num0.RecentStats.PlayerCount))
	steamPlayerCount := humanize.Comma(int64(steamResult.data.Response.PlayerCount))
	message := fmt.Sprintf(
		"There are currently %s players playing on Steam & %s players in Trials across all platforms | Trials data last updated: %s %s ago from https://trials.report",
		steamPlayerCount,
		trialsPlayerCount,
		humanize.Comma(elapsedMinutes),
		minuteLabel,
	)

	return &message
}
