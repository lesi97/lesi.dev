package store

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
	ani_utils "github.com/lesi97/lesi.dev/internal/domains/anilist/internal/utils/anilist"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func (s *Store) UpdateAnilist(ctx context.Context, data model.PlexWebhookPayload) error {

	showName, err := data.GetShowName()
	if err != nil {
		return err
	}

	plexSeason := data.Metadata.ParentIndex
	plexShowYear := data.Metadata.Year
	plexEpisode := data.GetEpisodeNumber()
	isSpecial := strings.EqualFold(strings.TrimSpace(data.Metadata.ParentTitle), "specials")

	allowed, err := s.PlexUtils.ValidateLabels(ctx, &data)
	if err != nil {
		return err
	}
	if !allowed {
		return nil
	}

	searchTitle := showName
	if isSpecial && strings.TrimSpace(data.Metadata.Title) != "" {
		searchTitle = fmt.Sprintf("%s %s", showName, data.Metadata.Title)
	}

	results, err := s.AniUtils.SearchTitle(ctx, searchTitle)
	if err != nil {
		return err
	}

	best := &model.AnilistMedia{}
	ok := false
	if isSpecial {
		best, ok = ani_utils.PickBestSpecialMatch(results, showName, data.Metadata.Title, plexShowYear)
	} else {
		best, ok = ani_utils.PickBestMatch(results, showName, plexShowYear)
	}
	if !ok {
		return fmt.Errorf("no anilist results found")
	}

	if best.IsAdult {
		return nil
	}

	targetMediaID := best.ID
	progress := plexEpisode

	if isSpecial {
		if best.Episodes != nil && *best.Episodes > 0 && progress > *best.Episodes {
			progress = *best.Episodes
		}
		if progress < 1 {
			progress = 1
		}
	} else {
		absEpisode, err := s.PlexUtils.AbsoluteEpisodeFromPlex(ctx, data.Metadata.GrandparentRatingKey, plexSeason, plexEpisode)
		if err != nil {
			return err
		}

		chain, chainErr := s.AniUtils.BuildSeasonChain(ctx, best.ID)
		if chainErr != nil {
			chain = nil
		}

		targetMediaID, progress = ResolveTargetProgress(*best, chain, plexSeason, plexEpisode, absEpisode)

		if absEpisode%100 == 0 {
			s.Logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
				Content:  "",
				Title:    fmt.Sprintf("Episode %v of %v watched!", absEpisode, showName),
				Username: "Anilist Milestone",
			})
		}
	}

	if preview, ok := ctx.Value(updatePreviewKey).(*UpdatePreview); ok && preview != nil {
		preview.TargetMediaID = targetMediaID
		preview.Progress = progress
		return nil
	}

	if os.Getenv("GO_ENV") == "development" {
		s.Logger.Printf("SHOW NAME: %v\n", showName)
		s.Logger.Printf("MEDIA ID: %v\n", targetMediaID)
		s.Logger.Printf("PROGRESS: %v\n", progress)
		return nil
	} 
	_ = s.AniUtils.UpdateProgress(ctx, targetMediaID, progress)
	return nil
}
