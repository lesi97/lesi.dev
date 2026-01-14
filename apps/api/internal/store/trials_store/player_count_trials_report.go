package trials_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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

func (s *TrialsStore) fetchFromTrialsReport() (*TrialsData, error) {
	defer s.Logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", s.url), time.Now(), nil)

	req, err := http.NewRequest("GET", s.url, nil)
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