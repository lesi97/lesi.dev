package store

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/steam/internal/model"
	"github.com/lesi97/lesi.dev/internal/httpapi"
)

func (s *Store) ResolveGameIDByNameFromAppList(ctx context.Context, gameName string) (string, string, error) {
	appListURL := fmt.Sprintf("%s/ISteamApps/GetAppList/v2/", s.BaseURL)
	body, statusCode, err := httpapi.DoRequest(ctx, s.HTTPClient, http.MethodGet, appListURL, nil, nil)
	if err != nil {
		return "", "", err
	}

	if statusCode < 200 || statusCode >= 300 {
		return "", "", fmt.Errorf("failed to fetch steam app list")
	}

	var appListData model.SteamAppListData
	err = json.Unmarshal(body, &appListData)
	if err != nil {
		return "", "", err
	}

	normalisedTarget := NormaliseGameName(gameName)
	for _, app := range appListData.AppList.Apps {
		trimmedAppName := strings.TrimSpace(app.Name)
		if trimmedAppName == "" {
			continue
		}

		normalisedAppName := NormaliseGameName(trimmedAppName)
		if normalisedAppName == normalisedTarget {
			return fmt.Sprintf("%d", app.AppID), trimmedAppName, nil
		}
	}

	for _, app := range appListData.AppList.Apps {
		trimmedAppName := strings.TrimSpace(app.Name)
		if trimmedAppName == "" {
			continue
		}

		normalisedAppName := NormaliseGameName(trimmedAppName)
		if strings.Contains(normalisedAppName, normalisedTarget) || strings.Contains(normalisedTarget, normalisedAppName) {
			return fmt.Sprintf("%d", app.AppID), trimmedAppName, nil
		}
	}

	targetTokens := strings.Fields(normalisedTarget)
	bestMatchID := int64(0)
	bestMatchName := ""
	bestScore := 0
	for _, app := range appListData.AppList.Apps {
		trimmedAppName := strings.TrimSpace(app.Name)
		if trimmedAppName == "" {
			continue
		}

		normalisedAppName := NormaliseGameName(trimmedAppName)
		if normalisedAppName == "" {
			continue
		}

		matchedTokenCount := 0
		for _, token := range targetTokens {
			if token != "" && strings.Contains(normalisedAppName, token) {
				matchedTokenCount++
			}
		}

		if matchedTokenCount > bestScore {
			bestScore = matchedTokenCount
			bestMatchID = app.AppID
			bestMatchName = trimmedAppName
		}
	}

	if bestScore > 0 && bestMatchName != "" {
		return fmt.Sprintf("%d", bestMatchID), bestMatchName, nil
	}

	return "", "", nil
}
