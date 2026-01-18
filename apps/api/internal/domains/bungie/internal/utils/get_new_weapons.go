package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func GetNewWeapons(database *db.DB, logger *utils.Logger, baseURL string, clientID string) {
	urlPath, err := GetManifestURL(logger, clientID, baseURL)
	if err != nil {
		fmt.Printf("failed to generate manifest URL")
		return
	}

	url := fmt.Sprintf("%s%s", baseURL, *urlPath)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to get manifest")
		return
	}

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
