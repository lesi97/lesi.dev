package middleware

import (
	"encoding/json"
	"net/http"
	"time"
)

type profile struct {
	Title       string `json:"title"`
	HeaderImage string `json:"headerImage"`
}

type streamElementsChannelResponse struct {
	Profile 		profile `json:"profile"`
	ID              string `json:"_id"`
	Provider        string `json:"provider"`
	BroadcasterType string `json:"broadcasterType"`
	Suspended       bool   `json:"suspended"`
	ProviderID      string `json:"providerId"`
	Avatar          string `json:"avatar"`
	Username        string `json:"username"`
	Alias           string `json:"alias"`
	DisplayName     string `json:"displayName"`
	Inactive        bool   `json:"inactive"`
	IsPartner       bool   `json:"isPartner"`
}

func FetchStreamElementsChannelDisplayName(channelId string) (string, bool) {
	if channelId == "" {
		return "", false
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.streamelements.com/kappa/v2/channels/"+channelId, nil)
	if err != nil {
		return "", false
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false
	}

	var payload streamElementsChannelResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", false
	}

	if payload.DisplayName == "" {
		return "", false
	}

	return payload.DisplayName, true
}
