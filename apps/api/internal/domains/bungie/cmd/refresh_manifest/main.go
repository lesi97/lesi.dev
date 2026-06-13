package main

import (
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lesi97/lesi.dev/internal/db"
	bungie_utils "github.com/lesi97/lesi.dev/internal/domains/bungie/internal/utils"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func main() {
	if os.Getenv("GO_ENV") != "production" {
		_ = godotenv.Load(".env.local")
	}

	logger := utils.NewColourLogger("brightBlack")

	database, err := db.Connect(logger)
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	clientID := os.Getenv("BUNGIE_CLIENT_ID")
	if clientID == "" {
		logger.Fatal("BUNGIE_CLIENT_ID not found in env")
	}

	httpClient := &http.Client{
		Timeout: 2 * time.Minute,
		Transport: &http.Transport{
			Proxy:             nil,
			ForceAttemptHTTP2: true,
		},
	}

	if err := bungie_utils.GetNewWeapons(database, logger, httpClient, "https://www.bungie.net", clientID); err != nil {
		logger.Fatalf("failed to refresh Bungie manifest: %v", err)
	}

	logger.Println("Bungie manifest refresh complete")
}
