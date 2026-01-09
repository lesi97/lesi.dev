package anilist_store

import (
	"context"
	"fmt"
	"strings"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (s *AnilistStore) UpdateAnilist(ctx context.Context, data PlexWebhookPayload) error {
	showName := data.Metadata.GrandparentTitle
	if showName == "" {
		showName = data.Metadata.ParentTitle
	}
	if showName == "" {
		showName = data.Metadata.Title
	}
	if showName == "" {
		return fmt.Errorf("unknown showname")
	}

	plexSeason := data.Metadata.ParentIndex
	plexShowYear := data.Metadata.Year
	plexEpisode := 1
	if strings.ToLower(data.Metadata.Type) == "episode" {
		plexEpisode = data.Metadata.Index
	}

	blockedLabels := map[string]bool{
		"no_anilist": true,
		"no_sync":    true,
		"private":    true,
	}

	episodeLabels, err := s.plexLabels(data.Metadata.RatingKey)
	if hasBlockedLabel(episodeLabels, blockedLabels) {
		return nil
	}
	if err != nil {
		return err
	}

	seasonLabels, err := s.plexLabels(data.Metadata.ParentRatingKey)
	if hasBlockedLabel(seasonLabels, blockedLabels) {
		return nil
	}
	if err != nil {
		return err
	}

	showLabels, err := s.plexLabels(data.Metadata.GrandparentRatingKey)
	if hasBlockedLabel(showLabels, blockedLabels) {
		return nil
	}
	if err != nil {
		return err
	}

	results, err := s.searchTitle(showName)
	if err != nil {
		return err
	}

	best, ok := pickBestAniListMatch(results, showName, plexShowYear)
	if !ok {
		return fmt.Errorf("no anilist results found")
	}

	if best.IsAdult {
		return nil
	}

	absEpisode, err := s.absoluteEpisodeFromPlex(data.Metadata.GrandparentRatingKey, plexSeason, plexEpisode)
	if err != nil {
		return err
	}

	targetMediaID := best.ID
	progress := absEpisode

	chain, chainErr := s.buildSeasonChain(best.ID)
	utils.PrintPrettyJSON(chain)

	if best.Episodes != nil && *best.Episodes > 0 {
		// If the base entry has a known episode count
		// - and our absolute episode is within that range, update the base entry with absolute progress
		// - otherwise, we have gone past the base series and should spill into sequels
		if absEpisode <= *best.Episodes {
			targetMediaID = best.ID
			progress = absEpisode
		} else if chainErr == nil && len(chain) > 0 {
			// Walk the chain subtracting each entry's episode count until we land on the correct sequel
			remaining := absEpisode

			for i := range chain {
				epCount := 0
				if chain[i].Episodes != nil {
					epCount = *chain[i].Episodes
				}

				// If this entry has a known episode count and remaining fits, this is our target
				if epCount > 0 && remaining <= epCount {
					targetMediaID = chain[i].ID
					progress = remaining
					break
				}

				// If episode count is unknown, we cannot safely advance, so we only subtract when known
				if epCount > 0 {
					remaining -= epCount
				}
			}
		}
	} else {
		// Fallback behaviour when the base entry has unknown episode count
		// If the AniList chain exists and Plex season looks like it maps to seasonal entries,
		// prefer seasonal progress (episode within season) against the appropriate chain entry
		// This helps "true seasonal shows" where AniList splits each season into a separate entry
		if chainErr == nil && plexSeason >= 1 && plexSeason <= len(chain) {
			targetMediaID = chain[plexSeason-1].ID
			progress = plexEpisode
		}
	}


	_ = s.updateAniListProgress(targetMediaID, progress)
	return nil
}
