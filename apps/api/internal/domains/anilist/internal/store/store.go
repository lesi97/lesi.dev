package store

import (
	"context"

	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
	ani_utils "github.com/lesi97/lesi.dev/internal/domains/anilist/internal/utils/anilist"
	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/utils/plex"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Methods interface {
	UpdateAnilist(ctx context.Context, data model.PlexWebhookPayload) error
	UpdatePlexLabels(ctx context.Context, data model.PlexWebhookPayload) error
}

type Store struct {
	DB        *db.DB
	Logger    *utils.Logger
	AniEnv    *model.AnilistEnv
	PlexEnv   *model.PlexEnv
	AniUtils  *ani_utils.Store
	PlexUtils *plex.Store
}

func NewStore(db *db.DB, logger *utils.Logger) (*Store, error) {

	aniEnv := &model.AnilistEnv{}
	err := aniEnv.Validate(db, logger)
	if err != nil {
		return nil, err
	}

	plexEnv := &model.PlexEnv{}
	err = plexEnv.Validate()
	if err != nil {
		return nil, err
	}

	aniUtils := ani_utils.NewStore(logger, aniEnv, db)
	plexUtils := plex.NewStore(logger, plexEnv)

	return &Store{
		DB:        db,
		Logger:    logger,
		AniEnv:    aniEnv,
		PlexEnv:   plexEnv,
		AniUtils:  aniUtils,
		PlexUtils: plexUtils,
	}, nil
}
