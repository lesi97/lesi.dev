package bootstrap

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/lesi97/lesi.dev/internal/db"
	aim_trainer_handler "github.com/lesi97/lesi.dev/internal/domains/aim_trainer/handler"
	anilist_handler "github.com/lesi97/lesi.dev/internal/domains/anilist/handler"
	auth_handler "github.com/lesi97/lesi.dev/internal/domains/auth/handler"
	bungie_handler "github.com/lesi97/lesi.dev/internal/domains/bungie/handler"
	countdown_handler "github.com/lesi97/lesi.dev/internal/domains/countdown/handler"
	local_handler "github.com/lesi97/lesi.dev/internal/domains/local/handler"
	tarot_handler "github.com/lesi97/lesi.dev/internal/domains/tarot/handler"
	time_handler "github.com/lesi97/lesi.dev/internal/domains/time/handler"
	trials_handler "github.com/lesi97/lesi.dev/internal/domains/trials/handler"
	twitch_handler "github.com/lesi97/lesi.dev/internal/domains/twitch/handler"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Logger            *utils.Logger
	DB                *db.DB
	Redis             *redis.Client
	tarot             *tarot_handler.Handler
	CountdownHandler  *countdown_handler.Handler
	time              *time_handler.Handler
	BungieHandler     *bungie_handler.Handler
	TrialsHandler     *trials_handler.Handler
	AnilistHandler    *anilist_handler.Handler
	LocalHandler      *local_handler.Handler
	TwitchHandler     *twitch_handler.Handler
	AuthHandler       *auth_handler.Handler
	AimTrainerHandler *aim_trainer_handler.Handler
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
	db, err := db.Connect(logger)
	if err != nil {
		return nil, nil, err
	}
	cfg := RedisLoadConfig()
	rdb := RedisNew(cfg)

	tarot := tarot_handler.NewHandler(logger)
	time := time_handler.NewHandler(logger)

	trialsHandler := trials_handler.NewHandler(logger, db)
	countdownHandler := countdown_handler.NewHandler(logger, db)
	localHandler := local_handler.NewHandler(logger, db)
	aimTrainerHandler := aim_trainer_handler.NewHandler(logger, db)

	var anilistHandler *anilist_handler.Handler
	anilistHandler, anilistErr := anilist_handler.NewHandler(logger, db)
	if anilistErr != nil {
		logger.Error("AniList store disabled: " + anilistErr.Error())
		anilistHandler = nil
	}

	var bungieHandler *bungie_handler.Handler
	bungieHandler, bungieErr := bungie_handler.NewHandler(logger, db, rdb)
	if bungieErr != nil {
		logger.Error("Bungie store disabled: " + bungieErr.Error())
		bungieHandler = nil
	}

	var twitchHandler *twitch_handler.Handler
	twitchHandler, twitchErr := twitch_handler.NewHandler(logger, db, rdb)
	if twitchErr != nil {
		logger.Error("Twitch store disabled: " + twitchErr.Error())
		twitchHandler = nil
	}

	var authHandler *auth_handler.Handler
	authHandler, authErr := auth_handler.NewHandler(logger, db)
	if authErr != nil {
		logger.Error("Auth store disabled: " + authErr.Error())
		authHandler = nil
	}

	app := &Application{
		Logger:            logger,
		DB:                db,
		Redis:             rdb,
		tarot:             tarot,
		time:              time,
		TrialsHandler:     trialsHandler,
		BungieHandler:     bungieHandler,
		AnilistHandler:    anilistHandler,
		CountdownHandler:  countdownHandler,
		LocalHandler:      localHandler,
		TwitchHandler:     twitchHandler,
		AuthHandler:       authHandler,
		AimTrainerHandler: aimTrainerHandler,
	}

	routes := setupRoutes(app)

	return app, routes, nil
}
