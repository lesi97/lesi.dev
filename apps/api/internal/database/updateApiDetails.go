package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lesi97/lesi.dev/internal/utils"
)

func (db *DB) UpdateApiDetails(
	ctx context.Context,
	application string,
	clientID *string,
	clientSecret *string,
	accessToken *string,
	refreshToken *string,
	refreshTokenExpiry *int64,
	baseURL *string,
	redirectURL *string,
	logger *utils.Logger,
) error {
	if !isAllowedApplication(application) {
		if logger != nil {
			logger.Errorf("Invalid application: " + application)
		}
		return errors.New("invalid application")
	}

	cols := []string{"name"}
	args := []any{application}

	if accessToken != nil {
		enc, err := Encrypt(*accessToken)
		if err != nil {
			return err
		}
		cols = append(cols, "access_token")
		args = append(args, enc)
	}

	if refreshToken != nil {
		enc, err := Encrypt(*refreshToken)
		if err != nil {
			return err
		}
		cols = append(cols, "refresh_token")
		args = append(args, enc)
	}

	if refreshTokenExpiry != nil {
		cols = append(cols, "refresh_token_expiry")
		args = append(args, *refreshTokenExpiry)
	}

	if clientID != nil {
		enc, err := Encrypt(*clientID)
		if err != nil {
			return err
		}
		cols = append(cols, "client_id")
		args = append(args, enc)
	}

	if clientSecret != nil {
		enc, err := Encrypt(*clientSecret)
		if err != nil {
			return err
		}
		cols = append(cols, "client_secret")
		args = append(args, enc)
	}

	if baseURL != nil {
		cols = append(cols, "base_url")
		args = append(args, *baseURL)
	}

	if redirectURL != nil {
		cols = append(cols, "redirect_url")
		args = append(args, *redirectURL)
	}

	placeholders := make([]string, 0, len(cols))
	for i := range cols {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}

	if len(cols) == 1 {
		q := fmt.Sprintf(
			"insert into api_keys (%s) values (%s) on conflict (name) do nothing",
			strings.Join(cols, ", "),
			strings.Join(placeholders, ", "),
		)
		_, err := db.Exec(ctx, q, args...)
		if err != nil && logger != nil {
			logger.Errorf("Error upserting record: " + err.Error())
		}
		return err
	}

	setParts := make([]string, 0, len(cols)-1)
	for i := 1; i < len(cols); i++ {
		col := cols[i]
		setParts = append(setParts, fmt.Sprintf("%s = excluded.%s", col, col))
	}

	query := fmt.Sprintf(
		"insert into api_keys (%s) values (%s) on conflict (name) do update set %s",
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(setParts, ", "),
	)

	_, err := db.Exec(ctx, query, args...)
	if err != nil && logger != nil {
		logger.Errorf("Error upserting record: " + err.Error())
	}

	return err
}
