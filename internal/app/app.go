package app

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lesi97/api.lesi.dev/internal/api"
	"github.com/lesi97/api.lesi.dev/internal/database"
	"github.com/lesi97/api.lesi.dev/internal/store"
	"github.com/lesi97/api.lesi.dev/internal/store/bungie_store"
	"github.com/lesi97/api.lesi.dev/internal/utils"
)

type Application struct {
	Logger				*log.Logger
	DB 					*database.Supabase
	TarotHandler		*api.TarotHandler
	CountdownHandler	*api.CountdownHandler
	TimeapiHandler		*api.TimeapiHandler
	BungieHandler		*api.BungieHandler
}

func init() {
	err := godotenv.Overload()
	if err != nil {
		panic(".env file not loaded")
	}

}

func NewApplication() (*Application, error) {
	logger := utils.NewColourLogger("brightMagenta")
	supabase, err := database.Connect(logger)
	if err != nil {
		return nil, err
	}
	
	tarotStore := store.NewTarotStore()
	countdownStore := store.NewSupabaseCountdownStore(supabase)
	bungieStore := bungie_store.NewSupabaseBungieStore(supabase, logger)

	tarotHandler := api.NewTarotHandler(logger, tarotStore)
	countdownHandler := api.NewCountdownHandler(logger, countdownStore)
	timeapiHandler := api.NewTimeapiHandler(logger)
	bungieHandler := api.NewBungieHandler(logger, bungieStore)

	app := &Application{
		Logger: logger,
		TarotHandler: tarotHandler,
		CountdownHandler: countdownHandler,
		TimeapiHandler: timeapiHandler,
		BungieHandler: bungieHandler,
	}

	return app, nil
}