package app

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/lesi97/lesi.dev/internal/api"
	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store/anilist_store"
	"github.com/lesi97/lesi.dev/internal/store/bungie_store"
	"github.com/lesi97/lesi.dev/internal/store/countdown_store"
	"github.com/lesi97/lesi.dev/internal/store/tarot_store"
	"github.com/lesi97/lesi.dev/internal/store/trials_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Application struct {
	Logger				*utils.Logger
	DB					*database.DB
	TarotHandler		*api.TarotHandler
	CountdownHandler	*api.CountdownHandler
	TimeapiHandler		*api.TimeapiHandler
	BungieHandler		*api.BungieHandler
	TrialsHandler		*api.TrialsHandler
	AnilistHandler 		*api.AnilistHandler
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
	
	tarotStore := tarot_store.NewStore(logger)
	trialsStore := trials_store.NewStore(db, logger)
	bungieStore := bungie_store.NewStore(db, logger)
	countdownStore := countdown_store.NewStore(db, logger)

	tarotHandler := api.NewTarotHandler(logger, tarotStore)
	trialsHandler := api.NewTrialsHandler(logger, trialsStore)
	bungieHandler := api.NewBungieHandler(logger, bungieStore)
	countdownHandler := api.NewCountdownHandler(logger, countdownStore)
	timeapiHandler := api.NewTimeapiHandler(logger)

	var anilistHandler *api.AnilistHandler
	anilistStore, anilistErr := anilist_store.NewStore(db, logger)
	if anilistErr != nil {
		logger.Error("AniList store disabled: " + anilistErr.Error())
	} else {
		anilistHandler = api.NewAnilistHandler(logger, anilistStore)
	}

	app := &Application{
		Logger: logger,
		DB: db,
		TarotHandler: tarotHandler,
		TrialsHandler: trialsHandler,
		BungieHandler: bungieHandler,
		AnilistHandler: anilistHandler,
		CountdownHandler: countdownHandler,
		TimeapiHandler: timeapiHandler,
	}

	routes := setupRoutes(app)

	return app, routes, nil
}