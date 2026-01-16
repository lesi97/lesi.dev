package app

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/lesi97/lesi.dev/internal/api"
	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store/aim_trainer_store"
	"github.com/lesi97/lesi.dev/internal/store/anilist_store"
	"github.com/lesi97/lesi.dev/internal/store/auth_store"
	"github.com/lesi97/lesi.dev/internal/store/bungie_store"
	"github.com/lesi97/lesi.dev/internal/store/countdown_store"
	"github.com/lesi97/lesi.dev/internal/store/tarot_store"
	"github.com/lesi97/lesi.dev/internal/store/trials_store"
	"github.com/lesi97/lesi.dev/internal/store/twitch_store"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Logger				*utils.Logger
	DB					*database.DB
	Redis				*redis.Client
	TarotHandler		*api.TarotHandler
	CountdownHandler	*api.CountdownHandler
	TimeapiHandler		*api.TimeapiHandler
	BungieHandler		*api.BungieHandler
	TrialsHandler		*api.TrialsHandler
	AnilistHandler 		*api.AnilistHandler
	LocalHandler		*api.LocalHandler
	TwitchHandler		*api.TwitchHandler
	AuthHandler			*api.AuthHandler
	AimTrainerHandler	*api.AimTrainerHandler
}

func init() {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load(".env.local")
		if err != nil {
			panic(".env file not loaded")
		}
	}
}

func NewApplication() (*Application, *chi.Mux, error) {
	logger := utils.NewColourLogger("brightBlack")
	db, err := database.Connect(logger)
	if err != nil {
		return nil, nil, err
	}
	cfg := RedisLoadConfig()
	rdb := RedisNew(cfg)
	
	tarotStore := tarot_store.NewStore(logger)
	trialsStore := trials_store.NewStore(db, logger)
	countdownStore := countdown_store.NewStore(db, logger)
	aimTrainerStore := aim_trainer_store.NewStore(db, logger)

	tarotHandler := api.NewTarotHandler(logger, tarotStore)
	trialsHandler := api.NewTrialsHandler(logger, trialsStore)
	countdownHandler := api.NewCountdownHandler(logger, countdownStore)
	timeapiHandler := api.NewTimeapiHandler(logger)
	localHandler := api.NewLocalHandler(logger, db)
	aimTrainerHandler := api.NewAimTrainerHandler(logger, aimTrainerStore)

	var bungieHandler *api.BungieHandler
	bungieStore, bungieErr := bungie_store.NewStore(db, logger, rdb)	
	if bungieErr != nil {
		logger.Error("Bungie store disabled: " + bungieErr.Error())
	} else {
		bungieHandler = api.NewBungieHandler(logger, bungieStore)
	}

	var anilistHandler *api.AnilistHandler
	anilistStore, anilistErr := anilist_store.NewStore(db, logger)
	if anilistErr != nil {
		logger.Error("AniList store disabled: " + anilistErr.Error())
	} else {
		anilistHandler = api.NewAnilistHandler(logger, anilistStore)
	}

	var twitchHandler *api.TwitchHandler
	twitchStore, twitchErr := twitch_store.NewStore(db, logger, rdb)
	if twitchErr != nil {
		logger.Error("Twitch store disabled: " + twitchErr.Error())
	} else {
		twitchHandler = api.NewTwitchHandler(logger, twitchStore)
	}

	var authHandler *api.AuthHandler
	authStore, authErr := auth_store.NewStore(db, logger)
	if authErr != nil {
		logger.Error("Auth store disabled: " + authErr.Error())
	} else {
		authHandler = api.NewAuthHandler(logger, authStore)
	}

	app := &Application{
		Logger: logger,
		DB: db,
		Redis: rdb,
		TarotHandler: tarotHandler,
		TrialsHandler: trialsHandler,
		BungieHandler: bungieHandler,
		AnilistHandler: anilistHandler,
		CountdownHandler: countdownHandler,
		TimeapiHandler: timeapiHandler,
		LocalHandler: localHandler,
		TwitchHandler: twitchHandler,
		AuthHandler: authHandler,
		AimTrainerHandler: aimTrainerHandler,
	}

	routes := setupRoutes(app)

	return app, routes, nil
}