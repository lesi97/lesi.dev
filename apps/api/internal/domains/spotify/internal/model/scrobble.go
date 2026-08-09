package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	spotifyArtistURL = "https://open.spotify.com/artist/"
	spotifyAlbumURL  = "https://open.spotify.com/album/"
	spotifyTrackURL  = "https://open.spotify.com/track/"
)

type UnixTimestamp struct {
	Value int64
	Set   bool
}

func (u *UnixTimestamp) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		return nil
	}

	if strings.HasPrefix(raw, `"`) {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		raw = strings.TrimSpace(str)
		if raw == "" {
			return nil
		}
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fmt.Errorf("uts must be a unix timestamp")
	}

	u.Value = value
	u.Set = true
	return nil
}

type ScrobbleRequest struct {
	UTS             UnixTimestamp `json:"uts"`
	UTCTime         string        `json:"utc_time"`
	ScrobbledAt     string        `json:"scrobbled_at"`
	Artist          string        `json:"artist"`
	ArtistSpotifyID string        `json:"artist_spotify_id"`
	ArtistURL       string        `json:"artist_url"`
	ArtistImageURL  string        `json:"artist_image_url"`
	Album           string        `json:"album"`
	AlbumSpotifyID  string        `json:"album_spotify_id"`
	AlbumURL        string        `json:"album_url"`
	AlbumImageURL   string        `json:"album_image_url"`
	Track           string        `json:"track"`
	TrackSpotifyID  string        `json:"track_spotify_id"`
	TrackURL        string        `json:"track_url"`
	TrackImageURL   string        `json:"track_image_url"`
	Source          string        `json:"source"`
}

type ScrobbleEntityInput struct {
	Name      string
	SpotifyID *string
	URL       *string
	ImageURL  *string
}

type ScrobbleInput struct {
	ScrobbledAt time.Time
	Artist      ScrobbleEntityInput
	Album       *ScrobbleEntityInput
	Track       ScrobbleEntityInput
	Source      string
	RawPayload  string
}

type ScrobbleResult struct {
	ID              int64   `json:"id"`
	ArtistID        int64   `json:"artist_id"`
	AlbumID         *int64  `json:"album_id,omitempty"`
	TrackID         int64   `json:"track_id"`
	ScrobbledAt     string  `json:"scrobbled_at"`
	ArtistURL       *string `json:"artist_url,omitempty"`
	AlbumURL        *string `json:"album_url,omitempty"`
	TrackURL        *string `json:"track_url,omitempty"`
	ArtistImageURL  *string `json:"artist_image_url,omitempty"`
	AlbumImageURL   *string `json:"album_image_url,omitempty"`
	TrackImageURL   *string `json:"track_image_url,omitempty"`
	ArtistSpotifyID *string `json:"artist_spotify_id,omitempty"`
	AlbumSpotifyID  *string `json:"album_spotify_id,omitempty"`
	TrackSpotifyID  *string `json:"track_spotify_id,omitempty"`
}

type SpotifyPollInput struct {
	Limit   int
	AfterMS *int64
}

type SpotifyPollResult struct {
	Fetched                    int     `json:"fetched"`
	Scrobbled                  int     `json:"scrobbled"`
	Skipped                    int     `json:"skipped"`
	AfterMS                    *int64  `json:"after_ms,omitempty"`
	LatestPlayedAt             *string `json:"latest_played_at,omitempty"`
	Source                     string  `json:"source"`
	RateLimited                bool    `json:"rate_limited,omitempty"`
	RateLimitReason            string  `json:"rate_limit_reason,omitempty"`
	RateLimitRetryAfterSeconds *int    `json:"rate_limit_retry_after_seconds,omitempty"`
}

func (r ScrobbleRequest) ToInput(rawPayload []byte) (ScrobbleInput, error) {
	artist, err := normaliseRequiredEntity("artist", r.Artist, r.ArtistSpotifyID, r.ArtistURL, r.ArtistImageURL, spotifyArtistURL)
	if err != nil {
		return ScrobbleInput{}, err
	}

	track, err := normaliseRequiredEntity("track", r.Track, r.TrackSpotifyID, r.TrackURL, r.TrackImageURL, spotifyTrackURL)
	if err != nil {
		return ScrobbleInput{}, err
	}

	album, err := normaliseOptionalEntity("album", r.Album, r.AlbumSpotifyID, r.AlbumURL, r.AlbumImageURL, spotifyAlbumURL)
	if err != nil {
		return ScrobbleInput{}, err
	}

	scrobbledAt, err := r.normaliseScrobbledAt()
	if err != nil {
		return ScrobbleInput{}, err
	}

	source := strings.TrimSpace(r.Source)
	if source == "" {
		source = "api"
	}
	if len(source) > 100 {
		return ScrobbleInput{}, errors.New("source must be 100 characters or fewer")
	}

	raw := "{}"
	if len(rawPayload) > 0 {
		raw = string(rawPayload)
	}

	return ScrobbleInput{
		ScrobbledAt: scrobbledAt,
		Artist:      artist,
		Album:       album,
		Track:       track,
		Source:      source,
		RawPayload:  raw,
	}, nil
}

func (r ScrobbleRequest) normaliseScrobbledAt() (time.Time, error) {
	if r.UTS.Set {
		return time.Unix(r.UTS.Value, 0).UTC(), nil
	}

	if strings.TrimSpace(r.ScrobbledAt) != "" {
		scrobbledAt, err := time.Parse(time.RFC3339, strings.TrimSpace(r.ScrobbledAt))
		if err != nil {
			return time.Time{}, errors.New("scrobbled_at must be RFC3339")
		}
		return scrobbledAt.UTC(), nil
	}

	if strings.TrimSpace(r.UTCTime) != "" {
		scrobbledAt, err := time.ParseInLocation("02 Jan 2006, 15:04", strings.TrimSpace(r.UTCTime), time.UTC)
		if err != nil {
			return time.Time{}, errors.New("utc_time must use the Last.fm export format, for example 06 Aug 2026, 10:37")
		}
		return scrobbledAt.UTC(), nil
	}

	return time.Time{}, errors.New("uts, scrobbled_at, or utc_time is required")
}

func normaliseRequiredEntity(field string, name string, spotifyID string, explicitURL string, explicitImageURL string, spotifyURLPrefix string) (ScrobbleEntityInput, error) {
	entity, err := normaliseEntity(field, name, spotifyID, explicitURL, explicitImageURL, spotifyURLPrefix)
	if err != nil {
		return ScrobbleEntityInput{}, err
	}
	if entity == nil {
		return ScrobbleEntityInput{}, fmt.Errorf("%s is required", field)
	}
	return *entity, nil
}

func normaliseOptionalEntity(field string, name string, spotifyID string, explicitURL string, explicitImageURL string, spotifyURLPrefix string) (*ScrobbleEntityInput, error) {
	entity, err := normaliseEntity(field, name, spotifyID, explicitURL, explicitImageURL, spotifyURLPrefix)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func normaliseEntity(field string, name string, spotifyID string, explicitURL string, explicitImageURL string, spotifyURLPrefix string) (*ScrobbleEntityInput, error) {
	trimmedName := strings.TrimSpace(name)
	trimmedSpotifyID := strings.TrimSpace(spotifyID)
	trimmedURL := strings.TrimSpace(explicitURL)
	trimmedImageURL := strings.TrimSpace(explicitImageURL)

	if trimmedName == "" {
		if trimmedSpotifyID != "" || trimmedURL != "" || trimmedImageURL != "" {
			return nil, fmt.Errorf("%s is required when %s_spotify_id, %s_url, or %s_image_url is provided", field, field, field, field)
		}
		return nil, nil
	}

	var spotifyIDValue *string
	if trimmedSpotifyID != "" {
		spotifyIDValue = &trimmedSpotifyID
	}

	urlValue, err := normaliseURL(trimmedURL)
	if err != nil {
		return nil, fmt.Errorf("%s_url must be an http or https URL", field)
	}

	imageURLValue, err := normaliseURL(trimmedImageURL)
	if err != nil {
		return nil, fmt.Errorf("%s_image_url must be an http or https URL", field)
	}

	if urlValue == nil && spotifyIDValue != nil {
		generated := spotifyURLPrefix + *spotifyIDValue
		urlValue = &generated
	}

	return &ScrobbleEntityInput{
		Name:      trimmedName,
		SpotifyID: spotifyIDValue,
		URL:       urlValue,
		ImageURL:  imageURLValue,
	}, nil
}

func normaliseURL(value string) (*string, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, errors.New("unsupported scheme")
	}
	if parsed.Host == "" {
		return nil, errors.New("missing host")
	}

	return &value, nil
}
