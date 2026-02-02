package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func GetNewWeapons(database *db.DB, logger *utils.Logger, httpClient *http.Client, baseURL string, clientID string) {
	urlPath, err := GetManifestURL(logger, httpClient, clientID, baseURL)
	if err != nil {
		fmt.Printf("failed to generate manifest URL")
		return
	}

	url := fmt.Sprintf("%s%s", baseURL, *urlPath)

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("failed to create manifest request")
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("failed to get manifest")
		return
	}
	defer resp.Body.Close()

	var rawData map[string]bungieWeaponRaw
	err = json.NewDecoder(resp.Body).Decode(&rawData)
	if err != nil {
		fmt.Printf("failed to decode data")
		return
	}

	weapons := processWeapons(rawData)
	perks := processPerks(rawData)
	insertWeapons(database, logger, weapons)
	insertPerks(database, logger, perks)
}
