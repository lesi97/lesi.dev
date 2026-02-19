package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

func (s *Store) GetPlayerCount(ctx context.Context, gameID string, gameName string) (*string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	trimmedGameID := strings.TrimSpace(gameID)
	trimmedGameName := strings.TrimSpace(gameName)

	if trimmedGameID == "" && trimmedGameName == "" {
		return nil, errors.New("you must provide gameId or gameName")
	}

	if trimmedGameID != "" {
		_, err := strconv.ParseInt(trimmedGameID, 10, 64)
		if err != nil {
			return nil, errors.New("invalid gameId")
		}
	}

	resolvedGameID := trimmedGameID
	resolvedGameName := trimmedGameName

	if resolvedGameID == "" {
		resolvedGameIDFromName, resolvedGameNameFromName, err := s.ResolveGameIDByName(ctx, resolvedGameName)
		if err != nil {
			return nil, err
		}
		resolvedGameID = resolvedGameIDFromName
		resolvedGameName = resolvedGameNameFromName
	}

	if resolvedGameName == "" {
		resolvedGameNameFromID, err := s.ResolveGameNameByID(ctx, resolvedGameID)
		if err == nil {
			resolvedGameName = resolvedGameNameFromID
		}
	}

	playerCount, err := s.FetchPlayerCountByGameID(ctx, resolvedGameID)
	if err != nil {
		return nil, err
	}

	if resolvedGameName == "" {
		resolvedGameName = "app " + resolvedGameID
	}

	message := fmt.Sprintf(
		"There are currently %s players on Steam for %s",
		humanize.Comma(playerCount),
		resolvedGameName,
	)

	return &message, nil
}
