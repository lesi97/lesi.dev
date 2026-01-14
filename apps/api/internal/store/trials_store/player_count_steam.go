package trials_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SteamData struct {
	Response struct {
		PlayerCount int64 `json:"player_count"`
		Result      int `json:"result"`
	} `json:"response"`
}

func (s *TrialsStore) fetchFromSteam() (*SteamData, error) {
	const destiny2 = "1085660"
	url := fmt.Sprintf("%s/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=%s&key=%s", s.steamUrl, destiny2, s.steamClientId)
	defer s.Logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", url), time.Now(), nil)

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