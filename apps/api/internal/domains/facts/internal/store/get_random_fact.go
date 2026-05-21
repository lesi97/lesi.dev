package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const cheeseFactType = "cheese"

func (s *Store) GetRandomFact(ctx context.Context) (*string, error) {
	query := `
		SELECT fact
		FROM public.facts
		ORDER BY random()
		LIMIT 1
	`

	return s.randomFact(ctx, query)
}

func (s *Store) GetRandomCheeseFact(ctx context.Context) (*string, error) {
	query := `
		SELECT fact
		FROM public.facts
		WHERE lower(trim(fact_type)) = $1
		ORDER BY random()
		LIMIT 1
	`

	return s.randomFact(ctx, query, cheeseFactType)
}

func (s *Store) randomFact(ctx context.Context, query string, args ...any) (*string, error) {
	var fact string
	err := s.DB.QueryRow(ctx, query, args...).Scan(&fact)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("no facts found")
	}
	if err != nil {
		return nil, err
	}

	return &fact, nil
}
