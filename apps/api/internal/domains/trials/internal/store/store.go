package store

import (
	"github.com/lesi97/lesi.dev/internal/db"
	"github.com/lesi97/lesi.dev/internal/utils"
	"github.com/redis/go-redis/v9"
)

type Methods interface {
	GetLoot() *string
	GetPlayerCount() *string
}

type Store struct {
	DB                     *db.DB
	Logger                 *utils.Logger
	Redis                  *redis.Client
	URL                    string
	SteamClientID          string
	SteamURL               string
}
