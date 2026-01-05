package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/lesi97/lesi.dev/api/internal/database"
	"github.com/lesi97/lesi.dev/api/internal/utils"
)

type TrialsStore interface {
	GetLoot() *string
	GetPlayerCount() *string
}

type SupabaseTrialsStore struct {
	db     *database.Supabase
	logger *utils.Logger
}

type maps struct {
	Name      string `json:"name"`
	ImagePath string `json:"imagePath"`
}

type rewards struct {
	Flawless string `json:"flawless"`
}

type platforms struct {
	Num0 struct {
		RecentStats struct {
			HourEnding  string `json:"hourEnding"`
			PlayerCount int64    `json:"playerCount"`
			MatchCount  int    `json:"matchCount"`
			UpdatedAt   string `json:"updatedAt"`
		} `json:"recentStats"`
	} `json:"0"`
}

type TrialsData struct {
	WeekNumber     int       `json:"weekNumber"`
	RealWeekNumber int       `json:"realWeekNumber"`
	StartDate      string 	 `json:"startDate"`
	EndDate        string    `json:"endDate"`
	Maps           []maps    `json:"maps"`
	Rewards        rewards   `json:"rewards"`
	Platforms      platforms `json:"platforms"`
}

type SteamData struct {
	Response struct {
		PlayerCount int64 `json:"player_count"`
		Result      int `json:"result"`
	} `json:"response"`
}

func NewSupabaseTrialsStore(db *database.Supabase, logger *utils.Logger) *SupabaseTrialsStore {
	return &SupabaseTrialsStore{
		db: db,
		logger: logger,
	}
}

func (s *SupabaseTrialsStore) GetLoot() *string {

	trialsData, err := fetchFromTrialsReport()
	if err != nil {
		s.logger.Printf("ERROR: fetchFromTrialsReport: %v\n", err)
		message := "failed to fetch data from trials report"
		return &message
	}

	layout := "2006-01-02 15:04:05"
	endTime, err := time.Parse(layout, trialsData.EndDate)
	if err != nil {
		s.logger.Printf("ERROR: parsing EndDate: %v\n", err)
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

func (s *SupabaseTrialsStore) GetPlayerCount() *string {
	trialsCh := make(chan struct {
		data *TrialsData
		err  error
	})
	steamCh := make(chan struct {
		data *SteamData
		err  error
	})

	go func() {
		data, err := fetchFromTrialsReport()
		trialsCh <- struct {
			data *TrialsData
			err  error
		}{data, err}
	}()

	go func() {
		data, err := fetchFromSteam()
		steamCh <- struct {
			data *SteamData
			err  error
		}{data, err}
	}()

	trialsResult := <-trialsCh
	steamResult := <-steamCh

	if trialsResult.err != nil {
		s.logger.Printf("ERROR: fetchFromTrialsReport: %v\n", trialsResult.err)
		message := "failed to fetch data from trials report"
		return &message
	}

	if steamResult.err != nil {
		s.logger.Printf("ERROR: fetchFromSteam: %v\n", steamResult.err)
		message := "failed to fetch data from steam"
		return &message
	}

	layout := "2006-01-02 15:04:05" 
	updatedAtTime, err := time.Parse(layout, trialsResult.data.Platforms.Num0.RecentStats.UpdatedAt)
	if err != nil {
		s.logger.Printf("ERROR: parsing updatedAt time: %v\n", err)
		message := "failed to parse time, tell lesi"
		return &message
	}

	if time.Since(updatedAtTime) >= 90*time.Minute {
		message := fmt.Sprintf(
			"There are currently %v players playing Destiny 2 on Steam",
			steamResult.data.Response.PlayerCount,
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
		humanize.Comma(int64(steamResult.data.Response.PlayerCount)),
		humanize.Comma(int64(trialsResult.data.Platforms.Num0.RecentStats.PlayerCount)),
		humanize.Comma(elapsedMinutes),
		minuteLabel,
	)

	return &message
}


func fetchFromTrialsReport() (*TrialsData, error) {
	const url = "https://api.trialsofthenine.com/weeks/0"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &TrialsData{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}
	return result, nil
}

func fetchFromSteam() (*SteamData, error) {
	const steamURL = "https://api.steampowered.com"
	const destiny2 = "1085660";
	apiKey := os.Getenv("STEAM_CLIENT_ID")
	if apiKey == "" {
		return nil, fmt.Errorf("missing STEAM_CLIENT_ID in environment")
	}

	url := fmt.Sprintf("%s/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=%s&key=%s", steamURL, destiny2, apiKey)

	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	result := &SteamData{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}

	return result, nil
}