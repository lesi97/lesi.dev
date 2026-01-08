package trials_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type SteamData struct {
	Response struct {
		PlayerCount int64 `json:"player_count"`
		Result      int `json:"result"`
	} `json:"response"`
}

func fetchFromSteam() (*SteamData, error) {
	const steamURL = "https://api.steampowered.com"
	const destiny2 = "1085660"

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