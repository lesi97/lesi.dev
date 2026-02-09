package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/trials/internal/model"
	"github.com/lesi97/lesi.dev/internal/httpapi"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func FetchFromSteam(ctx context.Context, logger *utils.Logger, steamURL string, steamClientID string) (*model.SteamData, error) {
	const destiny2 = "1085660"
	url := fmt.Sprintf("%s/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=%s&key=%s", steamURL, destiny2, steamClientID)
	defer logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", url), time.Now(), nil)

	body, _, err := httpapi.DoRequest(ctx, &http.Client{}, http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	result := &model.SteamData{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}

	return result, nil
}
