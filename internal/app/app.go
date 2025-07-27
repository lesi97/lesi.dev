package app

import (
	"log"
	"os"

	"github.com/lesi97/api.lesi.dev/internal/api"
	"github.com/lesi97/api.lesi.dev/internal/database"
	"github.com/lesi97/api.lesi.dev/internal/store"
)

type Application struct {
	Logger			*log.Logger
	TarotHandler	*api.TarotHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	supabase, err := database.Connect(logger)
	if err != nil {
		return nil, err
	}
	defer database.Disconnect(supabase)
	

	tarotStore := store.NewTarotStore()
	// countdownStore := store.NewSupabaseCountdownStore()

	tarotHandler := api.NewTarotHandler(logger, tarotStore)
	// countdownHandler := api.NewTarotHandler(logger, tarotStore)

	app := &Application{
		Logger: logger,
		TarotHandler: tarotHandler,

	}

	return app, nil
}