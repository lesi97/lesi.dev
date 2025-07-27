package app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lesi97/api.lesi.dev/internal/api"
	"github.com/lesi97/api.lesi.dev/internal/database"
	"github.com/lesi97/api.lesi.dev/internal/store"
)

type Application struct {
	Logger				*log.Logger
	DB 					*database.Supabase
	TarotHandler		*api.TarotHandler
	CountdownHandler	*api.CountdownHandler
	TimeapiHandler		*api.TimeapiHandler
}

func init() {
	err := godotenv.Overload()
	if err != nil {
		panic(".env file not loaded")
	}
	if os.Getenv("GO_ENV") != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	supabase, err := database.Connect(logger)
	if err != nil {
		return nil, err
	}
	
	tarotStore := store.NewTarotStore()
	countdownStore := store.NewSupabaseCountdownStore(supabase)

	tarotHandler := api.NewTarotHandler(logger, tarotStore)
	countdownHandler := api.NewCountdownHandler(logger, countdownStore)
	timeapiHandler := api.NewTimeapiHandler(logger)

	app := &Application{
		Logger: logger,
		TarotHandler: tarotHandler,
		CountdownHandler: countdownHandler,
		TimeapiHandler: timeapiHandler,
	}

	return app, nil
}