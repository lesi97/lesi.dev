package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type Supabase = pgx.Conn

/*
Function that initiates a connection to a postgres db and sets the public Supabase var to be equal to that connection
*/
func Connect(logger *log.Logger) (*Supabase, error) {
	url := os.Getenv("DATABASE_URL")
	supabase, err := pgx.Connect(context.Background(), url)
	if err != nil {
		logger.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	}
	return supabase, nil
}

/*
Function to disconnect from the database once a connection is no longer required
*/
func Disconnect(supabase *Supabase) {
	supabase.Close(context.Background())
}