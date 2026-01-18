package store

import (
	"context"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/model"
	ani_utils "github.com/lesi97/lesi.dev/internal/domains/anilist/utils/anilist"
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

	allowed, err := s.PlexUtils.ValidateLabels(ctx, &data)
	if err != nil {
		return err
	}
	if !allowed {
		return nil
	}

	results, err := s.AniUtils.SearchTitle(ctx, showName)
	if err != nil {
		return err
	}

	best, ok := ani_utils.PickBestMatch(results, showName, plexShowYear)
	if !ok {
		return fmt.Errorf("no anilist results found")
	}

	if best.IsAdult {
		return nil
	}

	absEpisode, err := s.PlexUtils.AbsoluteEpisodeFromPlex(ctx, data.Metadata.GrandparentRatingKey, plexSeason, plexEpisode)
	if err != nil {
		return err
	}

	targetMediaID := best.ID
	progress := absEpisode

	chain, chainErr := s.AniUtils.BuildSeasonChain(ctx, best.ID)
	if chainErr != nil {
		chain = nil
	}

	targetMediaID, progress = ResolveTargetProgress(*best, chain, plexSeason, plexEpisode, absEpisode)

	if absEpisode%100 == 0 {
		s.Logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  "",
			Title:    fmt.Sprintf("Episde %v of %v watched!", absEpisode, showName),
			Username: "Anilist Milestone",
		})
	}

	if preview, ok := ctx.Value(updatePreviewKey).(*UpdatePreview); ok && preview != nil {
		preview.TargetMediaID = targetMediaID
		preview.Progress = progress
		return nil
	}

	_ = s.AniUtils.UpdateProgress(ctx, targetMediaID, progress)
	return nil
}
