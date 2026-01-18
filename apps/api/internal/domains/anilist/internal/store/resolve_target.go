package store

import "github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"

func ResolveTargetProgress(
	best model.AnilistMedia,
	chain []model.AnilistRelationNode,
	plexSeason int,
	plexEpisode int,
	absEpisode int,
) (int, int) {
	targetMediaID := best.ID
	progress := absEpisode

	if best.Episodes != nil && *best.Episodes > 0 {
		if absEpisode <= *best.Episodes {
			targetMediaID = best.ID
			progress = absEpisode
		} else if len(chain) > 0 {
			remaining := absEpisode

			for i := range chain {
				epCount := 0
				if chain[i].Episodes != nil {
					epCount = *chain[i].Episodes
				}

				if epCount > 0 && remaining <= epCount {
					targetMediaID = chain[i].ID
					progress = remaining
					break
				}

				if epCount > 0 {
					remaining -= epCount
				}
			}
		}
	} else {
		// Unknown episode count, map to seasonal chain if Plex season is in range.
		if len(chain) > 0 && plexSeason >= 1 && plexSeason <= len(chain) {
			targetMediaID = chain[plexSeason-1].ID
			progress = plexEpisode
		}
	}

	return targetMediaID, progress
}
