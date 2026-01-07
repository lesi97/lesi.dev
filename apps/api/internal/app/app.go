package app

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/lesi97/lesi.dev/internal/api"
	"github.com/lesi97/lesi.dev/internal/database"
	"github.com/lesi97/lesi.dev/internal/store"
	"github.com/lesi97/lesi.dev/internal/store/bungie_store"
	"github.com/lesi97/lesi.dev/internal/store/countdown_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type Application struct {
	Logger				*utils.Logger
	DB 					*database.Supabase
	FrontendHandler		*api.FrontendHandler
	TarotHandler		*api.TarotHandler
	CountdownHandler	*api.CountdownHandler
	TimeapiHandler		*api.TimeapiHandler
	BungieHandler		*api.BungieHandler
	TrialsHandler		*api.TrialsHandler
}

func init() {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load(".env.local")
		if err != nil {
			panic(".env file not loaded")
		}
	}
}

func NewApplication() (*Application, error) {
	logger := utils.NewColourLogger("brightBlack")
	supabase, err := database.Connect(logger)
	if err != nil {
		return nil, err
	}
	
	tarotStore := store.NewTarotStore()
	countdownStore := countdown_store.NewSupabaseCountdownStore(supabase)
	bungieStore := bungie_store.NewSupabaseBungieStore(supabase, logger)
	trialsStore := store.NewSupabaseTrialsStore(supabase, logger)

	frontendHandler := api.NewFrontendHandler(logger)
	tarotHandler := api.NewTarotHandler(logger, tarotStore)
	countdownHandler := api.NewCountdownHandler(logger, countdownStore)
	timeapiHandler := api.NewTimeapiHandler(logger)
	bungieHandler := api.NewBungieHandler(logger, bungieStore)
	trialsHandler := api.NewTrialsHandler(logger, trialsStore)

	app := &Application{
		Logger: logger,
		FrontendHandler: frontendHandler,
		TarotHandler: tarotHandler,
		CountdownHandler: countdownHandler,
		TimeapiHandler: timeapiHandler,
		BungieHandler: bungieHandler,
		TrialsHandler: trialsHandler,
	}

	return app, nil
}