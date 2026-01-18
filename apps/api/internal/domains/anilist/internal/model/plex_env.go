package model

import (
	"errors"
	"os"
)

type PlexEnv struct {
	BaseUrl  string
	Username string
	XToken   string
}

func (e *PlexEnv) Validate() error {

	baseUrl := os.Getenv("PLEX_API_URL")
	if baseUrl == "" {
		return errors.New("PLEX_API_URL not found in env")
	}

	username := os.Getenv("PLEX_USERNAME")
	if username == "" {
		return errors.New("PLEX_USERNAME not found in env")
	}

	xtoken := os.Getenv("PLEX_X_TOKEN")
	if xtoken == "" {
		return errors.New("PLEX_X_TOKEN not found in env")
	}

	e.BaseUrl = baseUrl
	e.Username = username
	e.XToken = xtoken

	return nil
}
