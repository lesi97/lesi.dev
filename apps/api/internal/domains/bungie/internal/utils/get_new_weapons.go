package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func GetNewWeapons(database *db.DB, logger *utils.Logger, httpClient *http.Client, baseURL string, clientID string) error {
	urlPath, err := GetManifestURL(logger, httpClient, clientID, baseURL)
	if err != nil {
		return fmt.Errorf("generate manifest URL: %w", err)
	}
	if urlPath == nil || *urlPath == "" {
		return fmt.Errorf("manifest missing DestinyInventoryItemDefinition path")
	}

	url := fmt.Sprintf("%s%s", baseURL, *urlPath)

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create manifest request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("get manifest %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("get manifest %s: unexpected status %d: %s", url, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var rawData map[string]bungieWeaponRaw
	err = json.NewDecoder(resp.Body).Decode(&rawData)
	if err != nil {
		return fmt.Errorf("decode weapon manifest %s (content-type %q): %w", url, resp.Header.Get("Content-Type"), err)
	}

	weapons := processWeapons(rawData)
	perks := processPerks(rawData)
	insertWeapons(database, logger, weapons)
	insertPerks(database, logger, perks)
	return nil
}
