package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/steam/internal/model"
	"github.com/lesi97/lesi.dev/internal/httpapi"
)

func (s *Store) FetchPlayerCountByGameID(ctx context.Context, gameID string) (int64, error) {
	url := fmt.Sprintf(
		"%s/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=%s&key=%s",
		s.BaseURL,
		gameID,
		s.SteamKey,
	)

	body, statusCode, err := httpapi.DoRequest(ctx, s.HTTPClient, http.MethodGet, url, nil, nil)
	if err != nil {
		safeURL := httpapi.RedactSensitiveQueryValues(url)
		safeError := strings.ReplaceAll(err.Error(), url, safeURL)
		return 0, errors.New(safeError)
	}

	if statusCode < 200 || statusCode >= 300 {
		return 0, errors.New("failed to fetch steam player count")
	}

	var playerData model.SteamCurrentPlayersData
	err = json.Unmarshal(body, &playerData)
	if err != nil {
		return 0, err
	}

	if playerData.Response.Result != 1 {
		return 0, errors.New("steam did not return a valid player count for this game")
	}

	return playerData.Response.PlayerCount, nil
}
