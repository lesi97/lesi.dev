package store

import (
	"context"
	"fmt"

	"github.com/lesi97/lesi.dev/internal/domains/telemetry/internal/model"
)

func (s *Store) UpsertDateMemeClick(ctx context.Context, input model.DateMemeClickInput) error {
	var yesIncrement int
	var noIncrement int

	switch input.Action {
	case model.DateMemeClickActionYes:
		yesIncrement = 1
	case model.DateMemeClickActionNo:
		noIncrement = 1
	default:
		return fmt.Errorf("unsupported date meme click action: %s", input.Action)
	}

	query := `
		INSERT INTO logs.date_meme_clicks (
			route,
			ip,
			click_date,
			yes_clicks,
			no_clicks,
			user_agent,
			first_clicked_at,
			last_clicked_at
		)
		VALUES ($1, $2, (now() at time zone 'utc')::date, $3, $4, $5, now(), now())
		ON CONFLICT (route, ip, click_date)
		DO UPDATE SET
			yes_clicks = logs.date_meme_clicks.yes_clicks + EXCLUDED.yes_clicks,
			no_clicks = logs.date_meme_clicks.no_clicks + EXCLUDED.no_clicks,
			user_agent = EXCLUDED.user_agent,
			last_clicked_at = now()
	`

	_, err := s.DB.Exec(ctx, query, input.Route, input.IP, yesIncrement, noIncrement, input.UserAgent)
	return err
}
