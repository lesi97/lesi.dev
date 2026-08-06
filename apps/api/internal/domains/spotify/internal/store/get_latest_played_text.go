package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lesi97/lesi.dev/internal/cache"
	"github.com/lesi97/lesi.dev/internal/httpapi"
)

const (
	currentlyPlayingCacheKey = "spotify:currently-playing:text"
	currentlyPlayingCacheTTL = 15 * time.Second
)

type spotifyCurrentlyPlayingResponse struct {
	IsPlaying            bool            `json:"is_playing"`
	CurrentlyPlayingType string          `json:"currently_playing_type"`
	Item                 json.RawMessage `json:"item"`
}

type latestPlayedTextResult struct {
	text *string
	err  error
}

func (s *Store) GetLatestPlayedText(ctx context.Context) (*string, error) {
	if cached := s.cachedCurrentlyPlayingText(ctx); cached != nil {
		return cached, nil
	}

	currentCh := make(chan latestPlayedTextResult, 1)
	scrobbledCh := make(chan latestPlayedTextResult, 1)

	go func() {
		text, err := s.currentlyPlayingText(ctx)
		currentCh <- latestPlayedTextResult{text: text, err: err}
	}()

	go func() {
		text, err := s.latestScrobbledText(ctx)
		scrobbledCh <- latestPlayedTextResult{text: text, err: err}
	}()

	var scrobbled latestPlayedTextResult
	scrobbledDone := false

	for {
		select {
		case current := <-currentCh:
			if current.err == nil && current.text != nil {
				s.setCurrentlyPlayingCache(ctx, current.text)
				return current.text, nil
			}
			if current.err == nil {
				s.setCurrentlyPlayingCache(ctx, nil)
			}
			if current.err != nil && s.Logger != nil {
				s.Logger.Printf("Spotify currently playing lookup skipped: %v", current.err)
			}

			if scrobbledDone {
				return scrobbled.text, scrobbled.err
			}

			select {
			case scrobbled = <-scrobbledCh:
				return scrobbled.text, scrobbled.err
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		case scrobbled = <-scrobbledCh:
			scrobbledDone = true
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (s *Store) currentlyPlayingText(ctx context.Context) (*string, error) {
	auth, err := s.spotifyAPIAuth(ctx)
	if err != nil {
		return nil, err
	}

	return s.currentlyPlayingTextWithAuth(ctx, auth)
}

func (s *Store) currentlyPlayingTextWithAuth(ctx context.Context, auth spotifyAPIAuth) (*string, error) {
	current, err := s.fetchSpotifyCurrentlyPlaying(ctx, auth.APIURL, auth.AccessToken)
	if errors.Is(err, errSpotifyUnauthorized) {
		accessToken, refreshErr := s.refreshSpotifyAccessToken(ctx, auth.Details)
		if refreshErr != nil {
			return nil, refreshErr
		}
		current, err = s.fetchSpotifyCurrentlyPlaying(ctx, auth.APIURL, accessToken)
	}
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, nil
	}

	return spotifyCurrentlyPlayingText(*current), nil
}

func (s *Store) refreshCurrentlyPlayingCache(ctx context.Context, auth spotifyAPIAuth) {
	current, err := s.currentlyPlayingTextWithAuth(ctx, auth)
	if err != nil {
		if s.Logger != nil {
			s.Logger.Printf("Spotify currently playing cache refresh skipped: %v", err)
		}
		return
	}

	s.setCurrentlyPlayingCache(ctx, current)
}

func (s *Store) refreshCurrentlyPlayingCacheAsync(ctx context.Context, auth spotifyAPIAuth) {
	cacheCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	go func() {
		defer cancel()
		s.refreshCurrentlyPlayingCache(cacheCtx, auth)
	}()
}

func (s *Store) fetchSpotifyCurrentlyPlaying(ctx context.Context, apiURL string, accessToken string) (*spotifyCurrentlyPlayingResponse, error) {
	endpoint, err := url.Parse(strings.TrimRight(apiURL, "/") + "/me/player/currently-playing")
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("additional_types", "track")
	endpoint.RawQuery = query.Encode()

	body, statusCode, err := httpapi.DoRequest(
		ctx,
		s.httpClient(),
		http.MethodGet,
		endpoint.String(),
		nil,
		map[string]string{
			"Authorization": "Bearer " + accessToken,
			"Accept":        "application/json",
		},
	)
	if err != nil {
		return nil, err
	}
	if statusCode == http.StatusNoContent {
		return nil, nil
	}
	if statusCode == http.StatusUnauthorized {
		return nil, errSpotifyUnauthorized
	}
	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("spotify currently playing failed: status=%d body=%s", statusCode, string(body))
	}
	if len(body) == 0 {
		return nil, nil
	}

	var response spotifyCurrentlyPlayingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func spotifyCurrentlyPlayingText(current spotifyCurrentlyPlayingResponse) *string {
	if !current.IsPlaying {
		return nil
	}
	if current.CurrentlyPlayingType != "" && current.CurrentlyPlayingType != "track" {
		return nil
	}
	if len(current.Item) == 0 || string(current.Item) == "null" {
		return nil
	}

	var track spotifyTrack
	if err := json.Unmarshal(current.Item, &track); err != nil {
		return nil
	}
	if strings.TrimSpace(track.Name) == "" || len(track.Artists) == 0 {
		return nil
	}

	artist := strings.TrimSpace(track.Artists[0].Name)
	if artist == "" {
		return nil
	}

	text := artist + " - " + strings.TrimSpace(track.Name)
	return &text
}

func (s *Store) cachedCurrentlyPlayingText(ctx context.Context) *string {
	if s.Redis == nil {
		return nil
	}

	text, err := s.Redis.Get(ctx, currentlyPlayingCacheKey).Result()
	if cache.IsRedisNil(err) {
		return nil
	}
	if err != nil {
		if s.Logger != nil {
			s.Logger.Printf("Redis currently playing cache read skipped: %v", err)
		}
		return nil
	}
	if strings.TrimSpace(text) == "" {
		return nil
	}

	return &text
}

func (s *Store) setCurrentlyPlayingCache(ctx context.Context, text *string) {
	if s.Redis == nil {
		return
	}
	if text == nil || strings.TrimSpace(*text) == "" {
		_ = s.Redis.Del(ctx, currentlyPlayingCacheKey).Err()
		return
	}

	_ = s.Redis.Set(ctx, currentlyPlayingCacheKey, strings.TrimSpace(*text), currentlyPlayingCacheTTL).Err()
}

func (s *Store) latestScrobbledText(ctx context.Context) (*string, error) {
	var text string
	err := s.DB.QueryRow(ctx, `
		SELECT a.name || ' - ' || t.name
		FROM music.scrobbles s
		JOIN music.artists a ON a.id = s.artist_id
		JOIN music.tracks t ON t.id = s.track_id
		ORDER BY s.scrobbled_at DESC, s.id DESC
		LIMIT 1
	`).Scan(&text)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &text, nil
}
