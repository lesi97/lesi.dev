package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
	"github.com/lesi97/lesi.dev/internal/httpapi"
)

const (
	lastfmApplication              = "LastFM"
	lastfmTagSource                = "lastfm"
	lastfmDefaultAPIURL            = "https://ws.audioscrobbler.com/2.0/"
	lastfmTagsRateLimitKey         = "lastfm:rate-limit:tags"
	lastfmDefaultTagLimit          = 8
	lastfmDefaultRateLimitCooldown = 120 * time.Second

	lastfmTagStatusTagged      = "tagged"
	lastfmTagStatusNoTags      = "no_tags"
	lastfmTagStatusError       = "error"
	lastfmTagStatusRateLimited = "rate_limited"
)

type lastfmAPIAuth struct {
	APIKey string
	APIURL string
}

type lastfmTagTarget struct {
	EntityType string
	EntityID   int64
	Artist     string
	Album      string
	Track      string
}

type lastfmTopTagsResponse struct {
	TopTags lastfmTopTags `json:"toptags"`
	Error   int           `json:"error"`
	Message string        `json:"message"`
}

type lastfmInfoTagsResponse struct {
	Track   lastfmInfoEntityTags `json:"track"`
	Album   lastfmInfoEntityTags `json:"album"`
	Artist  lastfmInfoEntityTags `json:"artist"`
	Error   int                  `json:"error"`
	Message string               `json:"message"`
}

type lastfmInfoEntityTags struct {
	TopTags lastfmTopTags `json:"toptags"`
	Tags    lastfmTopTags `json:"tags"`
}

type lastfmTopTags struct {
	Tag lastfmTagList `json:"tag"`
}

type lastfmTag struct {
	Name  string      `json:"name"`
	URL   string      `json:"url"`
	Count lastfmCount `json:"count"`
}

type lastfmTagLookupResult struct {
	Tags  []lastfmTag
	Notes string
}

type lastfmTagList []lastfmTag

type lastfmCount int

type lastfmRateLimitError struct {
	body string
}

type lastfmTagExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func (e *lastfmRateLimitError) Error() string {
	if strings.TrimSpace(e.body) == "" {
		return "last.fm tag lookup rate limited"
	}
	return "last.fm tag lookup rate limited: " + e.body
}

func (l *lastfmTagList) UnmarshalJSON(data []byte) error {
	data = []byte(strings.TrimSpace(string(data)))
	if len(data) == 0 || string(data) == "null" {
		*l = nil
		return nil
	}

	if data[0] == '"' {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		if strings.TrimSpace(value) == "" {
			*l = nil
			return nil
		}
		return fmt.Errorf("unexpected Last.fm tag string %q", value)
	}

	if data[0] == '[' {
		var tags []lastfmTag
		if err := json.Unmarshal(data, &tags); err != nil {
			return err
		}
		*l = tags
		return nil
	}

	var tag lastfmTag
	if err := json.Unmarshal(data, &tag); err != nil {
		return err
	}
	*l = []lastfmTag{tag}
	return nil
}

func (t *lastfmTopTags) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		*t = lastfmTopTags{}
		return nil
	}

	if strings.HasPrefix(raw, `"`) {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		if strings.TrimSpace(value) == "" {
			*t = lastfmTopTags{}
			return nil
		}
		return fmt.Errorf("unexpected Last.fm tags string %q", value)
	}

	type alias lastfmTopTags
	var parsed alias
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}
	*t = lastfmTopTags(parsed)
	return nil
}

func (c *lastfmCount) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		*c = 0
		return nil
	}

	if strings.HasPrefix(raw, `"`) {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		raw = strings.TrimSpace(value)
	}

	if raw == "" {
		*c = 0
		return nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return err
	}
	*c = lastfmCount(value)
	return nil
}

func (s *Store) EnrichLastFMTags(ctx context.Context, input model.LastFMTagEnrichmentInput) (*model.LastFMTagEnrichmentResult, error) {
	entityType := strings.ToLower(strings.TrimSpace(input.EntityType))
	if entityType == "" {
		entityType = model.SpotifyEnrichmentTypeTrack
	}
	if !isSpotifyEnrichmentType(entityType) {
		return nil, fmt.Errorf("type must be track, album, or artist")
	}
	if s.DB == nil {
		return nil, fmt.Errorf("music store database is not configured")
	}

	result := &model.LastFMTagEnrichmentResult{
		EntityType: entityType,
		Status:     "pending",
	}

	if retryAfter := s.lastfmRateLimitRetryAfter(ctx); retryAfter != nil {
		return lastfmTagRateLimitedResult(result, "last.fm tag cooldown active", *retryAfter), nil
	}

	target, ok, err := s.lastFMTagEnrichmentTarget(ctx, input)
	if err != nil {
		return nil, err
	}
	if !ok {
		result.Status = "complete"
		result.Notes = "no eligible " + entityType + " rows found"
		return result, nil
	}

	result.EntityID = &target.EntityID
	result.EntityName = target.LastFMEntityName()

	auth, ok, err := s.lastfmAPIAuth(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		result.Status = "skipped"
		result.Notes = "missing Last.fm API key; set LASTFM_API_KEY or api_keys row LastFM"
		return result, nil
	}

	lookup, err := s.fetchLastFMTopTags(ctx, auth, target)
	if err != nil {
		status := lastfmTagStatusError
		var rateLimitErr *lastfmRateLimitError
		if errors.As(err, &rateLimitErr) {
			status = lastfmTagStatusRateLimited
			retryAfter := s.rememberLastFMRateLimit(ctx)
			_ = s.recordLastFMTagAttempt(ctx, s.DB, target, status, 0, err.Error())
			return lastfmTagRateLimitedResult(result, err.Error(), retryAfter), nil
		}

		if recordErr := s.recordLastFMTagAttempt(ctx, s.DB, target, status, 0, err.Error()); recordErr != nil {
			return nil, recordErr
		}
		result.Status = status
		result.Notes = err.Error()
		return result, nil
	}

	tags := limitLastFMTags(lookup.Tags, lastfmTagLimit())
	tagsFound, err := s.persistLastFMTags(ctx, target, tags)
	if err != nil {
		return nil, err
	}

	result.TagsFound = tagsFound
	result.Tags = lastfmTagNames(tags)
	result.Notes = lookup.Notes
	if tagsFound == 0 {
		result.Status = lastfmTagStatusNoTags
		if result.Notes == "" {
			result.Notes = "last.fm returned no tags"
		}
		return result, nil
	}

	result.Status = lastfmTagStatusTagged
	result.Updated = true
	return result, nil
}

func (s *Store) enrichLastFMTagsForScrobble(ctx context.Context, input model.ScrobbleInput, inserted *model.ScrobbleResult, result *model.SpotifyPollResult) {
	if inserted == nil || result == nil || !lastfmTagPollingEnabled() {
		return
	}
	if !inserted.ArtistCreated && !inserted.AlbumCreated && !inserted.TrackCreated {
		return
	}

	auth, ok, err := s.lastfmAPIAuth(ctx)
	if err != nil {
		result.LastFMTagErrors++
		if s.Logger != nil {
			s.Logger.Printf("Last.fm tag enrichment skipped: %v", err)
		}
		return
	}
	if !ok {
		return
	}

	targets := lastfmTagTargetsForScrobble(input, inserted)
	for _, target := range targets {
		if retryAfter := s.lastfmRateLimitRetryAfter(ctx); retryAfter != nil {
			result.LastFMTagErrors++
			if err := s.recordLastFMTagAttempt(ctx, s.DB, target, lastfmTagStatusRateLimited, 0, "last.fm tag cooldown active"); err != nil && s.Logger != nil {
				s.Logger.Printf("Last.fm tag attempt record skipped: %v", err)
			}
			continue
		}

		result.LastFMTagLookups++
		lookup, err := s.fetchLastFMTopTags(ctx, auth, target)
		if err != nil {
			status := lastfmTagStatusError
			var rateLimitErr *lastfmRateLimitError
			if errors.As(err, &rateLimitErr) {
				status = lastfmTagStatusRateLimited
				s.rememberLastFMRateLimit(ctx)
			}

			result.LastFMTagErrors++
			if recordErr := s.recordLastFMTagAttempt(ctx, s.DB, target, status, 0, err.Error()); recordErr != nil && s.Logger != nil {
				s.Logger.Printf("Last.fm tag attempt record skipped: %v", recordErr)
			}
			if s.Logger != nil {
				s.Logger.Printf("Last.fm tag enrichment skipped for %s %d: %v", target.EntityType, target.EntityID, err)
			}
			continue
		}

		tags := limitLastFMTags(lookup.Tags, lastfmTagLimit())
		tagged, err := s.persistLastFMTags(ctx, target, tags)
		if err != nil {
			result.LastFMTagErrors++
			if s.Logger != nil {
				s.Logger.Printf("Last.fm tag persistence skipped for %s %d: %v", target.EntityType, target.EntityID, err)
			}
			continue
		}
		if tagged > 0 {
			result.LastFMTaggedEntities++
		}
	}
}

func lastfmTagTargetsForScrobble(input model.ScrobbleInput, inserted *model.ScrobbleResult) []lastfmTagTarget {
	targets := []lastfmTagTarget{}

	if inserted.ArtistCreated {
		targets = append(targets, lastfmTagTarget{
			EntityType: model.SpotifyEnrichmentTypeArtist,
			EntityID:   inserted.ArtistID,
			Artist:     input.Artist.Name,
		})
	}

	if inserted.AlbumCreated && inserted.AlbumID != nil && input.Album != nil {
		targets = append(targets, lastfmTagTarget{
			EntityType: model.SpotifyEnrichmentTypeAlbum,
			EntityID:   *inserted.AlbumID,
			Artist:     input.Artist.Name,
			Album:      input.Album.Name,
		})
	}

	if inserted.TrackCreated {
		targets = append(targets, lastfmTagTarget{
			EntityType: model.SpotifyEnrichmentTypeTrack,
			EntityID:   inserted.TrackID,
			Artist:     input.Artist.Name,
			Album:      scrobbleAlbumName(input),
			Track:      input.Track.Name,
		})
	}

	return targets
}

func scrobbleAlbumName(input model.ScrobbleInput) string {
	if input.Album == nil {
		return ""
	}
	return input.Album.Name
}

func (t lastfmTagTarget) LastFMEntityName() string {
	switch t.EntityType {
	case model.SpotifyEnrichmentTypeTrack:
		return strings.TrimSpace(t.Artist + " - " + t.Track)
	case model.SpotifyEnrichmentTypeAlbum:
		return strings.TrimSpace(t.Artist + " - " + t.Album)
	case model.SpotifyEnrichmentTypeArtist:
		return strings.TrimSpace(t.Artist)
	default:
		return ""
	}
}

func (s *Store) lastFMTagEnrichmentTarget(ctx context.Context, input model.LastFMTagEnrichmentInput) (lastfmTagTarget, bool, error) {
	entityType := strings.ToLower(strings.TrimSpace(input.EntityType))
	if input.EntityID != nil {
		return s.lastFMTagEnrichmentTargetByID(ctx, entityType, *input.EntityID)
	}

	return s.nextLastFMTagEnrichmentTarget(ctx, entityType)
}

func (s *Store) lastFMTagEnrichmentTargetByID(ctx context.Context, entityType string, entityID int64) (lastfmTagTarget, bool, error) {
	switch entityType {
	case model.SpotifyEnrichmentTypeTrack:
		return s.lastFMTrackTagEnrichmentTargetByID(ctx, entityID)
	case model.SpotifyEnrichmentTypeAlbum:
		return s.lastFMAlbumTagEnrichmentTargetByID(ctx, entityID)
	case model.SpotifyEnrichmentTypeArtist:
		return s.lastFMArtistTagEnrichmentTargetByID(ctx, entityID)
	default:
		return lastfmTagTarget{}, false, fmt.Errorf("type must be track, album, or artist")
	}
}

func (s *Store) nextLastFMTagEnrichmentTarget(ctx context.Context, entityType string) (lastfmTagTarget, bool, error) {
	switch entityType {
	case model.SpotifyEnrichmentTypeTrack:
		return s.nextLastFMTrackTagEnrichmentTarget(ctx)
	case model.SpotifyEnrichmentTypeAlbum:
		return s.nextLastFMAlbumTagEnrichmentTarget(ctx)
	case model.SpotifyEnrichmentTypeArtist:
		return s.nextLastFMArtistTagEnrichmentTarget(ctx)
	default:
		return lastfmTagTarget{}, false, fmt.Errorf("type must be track, album, or artist")
	}
}

func (s *Store) nextLastFMTrackTagEnrichmentTarget(ctx context.Context) (lastfmTagTarget, bool, error) {
	target := lastfmTagTarget{EntityType: model.SpotifyEnrichmentTypeTrack}
	err := s.DB.QueryRow(ctx, `
		SELECT
			t.id,
			a.name,
			coalesce(al.name, ''),
			t.name
		FROM music.tracks t
		JOIN music.artists a ON a.id = t.artist_id
		LEFT JOIN music.albums al ON al.id = t.album_id
		LEFT JOIN music.lastfm_tag_enrichment_attempts e
			ON e.entity_type = 'track'
			AND e.entity_id = t.id
		WHERE NOT EXISTS (
			SELECT 1
			FROM music.track_genres tg
			WHERE tg.track_id = t.id
				AND tg.source = $1
		)
		AND (
			e.id IS NULL
			OR (e.status = 'error' AND e.updated_at < now() - interval '6 hours')
			OR (e.status = 'rate_limited' AND e.updated_at < now() - interval '15 minutes')
		)
		ORDER BY
			coalesce(e.attempts, 0),
			(SELECT count(*) FROM music.scrobbles s WHERE s.track_id = t.id) DESC,
			coalesce(e.updated_at, 'epoch'::timestamptz),
			t.id
		LIMIT 1
	`, lastfmTagSource).Scan(&target.EntityID, &target.Artist, &target.Album, &target.Track)

	if errors.Is(err, pgx.ErrNoRows) {
		return lastfmTagTarget{}, false, nil
	}
	if err != nil {
		return lastfmTagTarget{}, false, err
	}
	return target, true, nil
}

func (s *Store) lastFMTrackTagEnrichmentTargetByID(ctx context.Context, entityID int64) (lastfmTagTarget, bool, error) {
	target := lastfmTagTarget{EntityType: model.SpotifyEnrichmentTypeTrack}
	err := s.DB.QueryRow(ctx, `
		SELECT
			t.id,
			a.name,
			coalesce(al.name, ''),
			t.name
		FROM music.tracks t
		JOIN music.artists a ON a.id = t.artist_id
		LEFT JOIN music.albums al ON al.id = t.album_id
		WHERE t.id = $1
	`, entityID).Scan(&target.EntityID, &target.Artist, &target.Album, &target.Track)

	if errors.Is(err, pgx.ErrNoRows) {
		return lastfmTagTarget{}, false, nil
	}
	if err != nil {
		return lastfmTagTarget{}, false, err
	}
	return target, true, nil
}

func (s *Store) nextLastFMAlbumTagEnrichmentTarget(ctx context.Context) (lastfmTagTarget, bool, error) {
	target := lastfmTagTarget{EntityType: model.SpotifyEnrichmentTypeAlbum}
	err := s.DB.QueryRow(ctx, `
		WITH candidates AS (
			SELECT
				al.id,
				a.name AS artist,
				al.name AS album,
				coalesce(e.attempts, 0) AS attempts,
				coalesce(e.updated_at, 'epoch'::timestamptz) AS last_attempted_at
			FROM music.albums al
			JOIN music.artists a ON a.id = al.artist_id
			LEFT JOIN music.lastfm_tag_enrichment_attempts e
				ON e.entity_type = 'album'
				AND e.entity_id = al.id
			WHERE NOT EXISTS (
				SELECT 1
				FROM music.album_genres alg
				WHERE alg.album_id = al.id
					AND alg.source = $1
			)
			AND (
				e.id IS NULL
				OR (e.status = 'error' AND e.updated_at < now() - interval '6 hours')
				OR (e.status = 'rate_limited' AND e.updated_at < now() - interval '15 minutes')
			)
		),
		album_play_counts AS (
			SELECT
				s.album_id,
				count(*)::bigint AS play_count
			FROM music.scrobbles s
			JOIN candidates c ON c.id = s.album_id
			WHERE s.album_id IS NOT NULL
			GROUP BY s.album_id
		)
		SELECT
			c.id,
			c.artist,
			c.album
		FROM candidates c
		LEFT JOIN album_play_counts apc ON apc.album_id = c.id
		ORDER BY
			c.attempts,
			coalesce(apc.play_count, 0) DESC,
			c.last_attempted_at,
			c.id
		LIMIT 1
	`, lastfmTagSource).Scan(&target.EntityID, &target.Artist, &target.Album)

	if errors.Is(err, pgx.ErrNoRows) {
		return lastfmTagTarget{}, false, nil
	}
	if err != nil {
		return lastfmTagTarget{}, false, err
	}
	return target, true, nil
}

func (s *Store) lastFMAlbumTagEnrichmentTargetByID(ctx context.Context, entityID int64) (lastfmTagTarget, bool, error) {
	target := lastfmTagTarget{EntityType: model.SpotifyEnrichmentTypeAlbum}
	err := s.DB.QueryRow(ctx, `
		SELECT
			al.id,
			a.name,
			al.name
		FROM music.albums al
		JOIN music.artists a ON a.id = al.artist_id
		WHERE al.id = $1
	`, entityID).Scan(&target.EntityID, &target.Artist, &target.Album)

	if errors.Is(err, pgx.ErrNoRows) {
		return lastfmTagTarget{}, false, nil
	}
	if err != nil {
		return lastfmTagTarget{}, false, err
	}
	return target, true, nil
}

func (s *Store) nextLastFMArtistTagEnrichmentTarget(ctx context.Context) (lastfmTagTarget, bool, error) {
	target := lastfmTagTarget{EntityType: model.SpotifyEnrichmentTypeArtist}
	err := s.DB.QueryRow(ctx, `
		SELECT
			a.id,
			a.name
		FROM music.artists a
		LEFT JOIN music.lastfm_tag_enrichment_attempts e
			ON e.entity_type = 'artist'
			AND e.entity_id = a.id
		WHERE NOT EXISTS (
			SELECT 1
			FROM music.artist_genres ag
			WHERE ag.artist_id = a.id
				AND ag.source = $1
		)
		AND (
			e.id IS NULL
			OR (e.status = 'error' AND e.updated_at < now() - interval '6 hours')
			OR (e.status = 'rate_limited' AND e.updated_at < now() - interval '15 minutes')
		)
		ORDER BY
			coalesce(e.attempts, 0),
			(SELECT count(*) FROM music.scrobbles s WHERE s.artist_id = a.id) DESC,
			coalesce(e.updated_at, 'epoch'::timestamptz),
			a.id
		LIMIT 1
	`, lastfmTagSource).Scan(&target.EntityID, &target.Artist)

	if errors.Is(err, pgx.ErrNoRows) {
		return lastfmTagTarget{}, false, nil
	}
	if err != nil {
		return lastfmTagTarget{}, false, err
	}
	return target, true, nil
}

func (s *Store) lastFMArtistTagEnrichmentTargetByID(ctx context.Context, entityID int64) (lastfmTagTarget, bool, error) {
	target := lastfmTagTarget{EntityType: model.SpotifyEnrichmentTypeArtist}
	err := s.DB.QueryRow(ctx, `
		SELECT
			a.id,
			a.name
		FROM music.artists a
		WHERE a.id = $1
	`, entityID).Scan(&target.EntityID, &target.Artist)

	if errors.Is(err, pgx.ErrNoRows) {
		return lastfmTagTarget{}, false, nil
	}
	if err != nil {
		return lastfmTagTarget{}, false, err
	}
	return target, true, nil
}

func (s *Store) lastfmAPIAuth(ctx context.Context) (lastfmAPIAuth, bool, error) {
	apiKey := strings.TrimSpace(os.Getenv("LASTFM_API_KEY"))
	apiURL := strings.TrimSpace(os.Getenv("LASTFM_API_URL"))

	if s.DB != nil {
		apiDetails, err := s.DB.FetchApiDetails(ctx, lastfmApplication, s.Logger)
		if err != nil {
			return lastfmAPIAuth{}, false, err
		}
		if apiDetails != nil {
			if key := stringValue(apiDetails.ClientID); key != "" {
				apiKey = key
			}
			if baseURL := stringValue(apiDetails.BaseURL); baseURL != "" {
				apiURL = baseURL
			}
		}
	}

	if apiKey == "" {
		return lastfmAPIAuth{}, false, nil
	}
	if apiURL == "" {
		apiURL = lastfmDefaultAPIURL
	}

	return lastfmAPIAuth{
		APIKey: apiKey,
		APIURL: normaliseLastFMAPIURL(apiURL),
	}, true, nil
}

func normaliseLastFMAPIURL(value string) string {
	value = strings.TrimRight(strings.TrimSpace(value), "/")
	if value == "" {
		return lastfmDefaultAPIURL
	}
	if strings.HasSuffix(value, "/2.0") {
		return value + "/"
	}
	return value + "/2.0/"
}

func (s *Store) fetchLastFMTopTags(ctx context.Context, auth lastfmAPIAuth, target lastfmTagTarget) (lastfmTagLookupResult, error) {
	tags, err := s.fetchLastFMTopTagsForTarget(ctx, auth, target)
	if err != nil || len(tags) > 0 {
		return lastfmTagLookupResult{Tags: tags}, err
	}

	tags, err = s.fetchLastFMInfoTagsForTarget(ctx, auth, target)
	if err != nil || len(tags) > 0 {
		return lastfmTagLookupResult{Tags: tags, Notes: "used Last.fm getInfo tags fallback"}, err
	}

	switch target.EntityType {
	case model.SpotifyEnrichmentTypeTrack:
		if strings.TrimSpace(target.Album) != "" {
			albumTarget := lastfmTagTarget{
				EntityType: model.SpotifyEnrichmentTypeAlbum,
				EntityID:   target.EntityID,
				Artist:     target.Artist,
				Album:      target.Album,
			}
			tags, err = s.fetchLastFMTopTagsForTarget(ctx, auth, albumTarget)
			if err != nil || len(tags) > 0 {
				return lastfmTagLookupResult{Tags: tags, Notes: "used Last.fm album tags fallback"}, err
			}
			tags, err = s.fetchLastFMInfoTagsForTarget(ctx, auth, albumTarget)
			if err != nil || len(tags) > 0 {
				return lastfmTagLookupResult{Tags: tags, Notes: "used Last.fm album info tags fallback"}, err
			}
		}

		artistTarget := lastfmTagTarget{
			EntityType: model.SpotifyEnrichmentTypeArtist,
			EntityID:   target.EntityID,
			Artist:     target.Artist,
		}
		tags, err = s.fetchLastFMTopTagsForTarget(ctx, auth, artistTarget)
		if err != nil || len(tags) > 0 {
			return lastfmTagLookupResult{Tags: tags, Notes: "used Last.fm artist tags fallback"}, err
		}
		tags, err = s.fetchLastFMInfoTagsForTarget(ctx, auth, artistTarget)
		if err != nil || len(tags) > 0 {
			return lastfmTagLookupResult{Tags: tags, Notes: "used Last.fm artist info tags fallback"}, err
		}
	case model.SpotifyEnrichmentTypeAlbum:
		artistTarget := lastfmTagTarget{
			EntityType: model.SpotifyEnrichmentTypeArtist,
			EntityID:   target.EntityID,
			Artist:     target.Artist,
		}
		tags, err = s.fetchLastFMTopTagsForTarget(ctx, auth, artistTarget)
		if err != nil || len(tags) > 0 {
			return lastfmTagLookupResult{Tags: tags, Notes: "used Last.fm artist tags fallback"}, err
		}
		tags, err = s.fetchLastFMInfoTagsForTarget(ctx, auth, artistTarget)
		if err != nil || len(tags) > 0 {
			return lastfmTagLookupResult{Tags: tags, Notes: "used Last.fm artist info tags fallback"}, err
		}
	}

	return lastfmTagLookupResult{}, nil
}

func (s *Store) fetchLastFMTopTagsForTarget(ctx context.Context, auth lastfmAPIAuth, target lastfmTagTarget) ([]lastfmTag, error) {
	method, err := lastfmTopTagsMethod(target.EntityType)
	if err != nil {
		return nil, err
	}

	body, err := s.fetchLastFMTagResponse(ctx, auth, target, method)
	if err != nil {
		return nil, err
	}

	return lastfmTopTagsFromResponse(body)
}

func (s *Store) fetchLastFMInfoTagsForTarget(ctx context.Context, auth lastfmAPIAuth, target lastfmTagTarget) ([]lastfmTag, error) {
	method, err := lastfmInfoMethod(target.EntityType)
	if err != nil {
		return nil, err
	}

	body, err := s.fetchLastFMTagResponse(ctx, auth, target, method)
	if err != nil {
		return nil, err
	}

	return lastfmInfoTagsFromResponse(body, target.EntityType)
}

func (s *Store) fetchLastFMTagResponse(ctx context.Context, auth lastfmAPIAuth, target lastfmTagTarget, method string) ([]byte, error) {
	endpoint, err := url.Parse(auth.APIURL)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("api_key", auth.APIKey)
	query.Set("format", "json")
	query.Set("autocorrect", "1")
	query.Set("method", method)

	switch target.EntityType {
	case model.SpotifyEnrichmentTypeArtist:
		query.Set("artist", target.Artist)
	case model.SpotifyEnrichmentTypeAlbum:
		query.Set("artist", target.Artist)
		query.Set("album", target.Album)
	case model.SpotifyEnrichmentTypeTrack:
		query.Set("artist", target.Artist)
		query.Set("track", target.Track)
	default:
		return nil, fmt.Errorf("unsupported Last.fm tag entity type %q", target.EntityType)
	}
	endpoint.RawQuery = query.Encode()

	body, statusCode, err := httpapi.DoRequest(
		ctx,
		s.httpClient(),
		http.MethodGet,
		endpoint.String(),
		nil,
		map[string]string{"Accept": "application/json"},
	)
	if err != nil {
		return nil, err
	}
	if statusCode == http.StatusTooManyRequests {
		return nil, &lastfmRateLimitError{body: string(body)}
	}
	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("last.fm tag lookup failed: status=%d body=%s", statusCode, string(body))
	}

	return body, nil
}

func lastfmTopTagsMethod(entityType string) (string, error) {
	switch entityType {
	case model.SpotifyEnrichmentTypeArtist:
		return "artist.getTopTags", nil
	case model.SpotifyEnrichmentTypeAlbum:
		return "album.getTopTags", nil
	case model.SpotifyEnrichmentTypeTrack:
		return "track.getTopTags", nil
	default:
		return "", fmt.Errorf("unsupported Last.fm tag entity type %q", entityType)
	}
}

func lastfmInfoMethod(entityType string) (string, error) {
	switch entityType {
	case model.SpotifyEnrichmentTypeArtist:
		return "artist.getInfo", nil
	case model.SpotifyEnrichmentTypeAlbum:
		return "album.getInfo", nil
	case model.SpotifyEnrichmentTypeTrack:
		return "track.getInfo", nil
	default:
		return "", fmt.Errorf("unsupported Last.fm tag entity type %q", entityType)
	}
}

func lastfmTopTagsFromResponse(body []byte) ([]lastfmTag, error) {
	var response lastfmTopTagsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Error != 0 {
		message := strings.TrimSpace(response.Message)
		if message == "" {
			message = fmt.Sprintf("last.fm error %d", response.Error)
		}

		switch response.Error {
		case 6, 7:
			return nil, nil
		case 29:
			return nil, &lastfmRateLimitError{body: message}
		default:
			return nil, fmt.Errorf("last.fm tag lookup failed: error=%d message=%s", response.Error, message)
		}
	}

	return normaliseLastFMTags(response.TopTags.Tag), nil
}

func lastfmInfoTagsFromResponse(body []byte, entityType string) ([]lastfmTag, error) {
	var response lastfmInfoTagsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Error != 0 {
		message := strings.TrimSpace(response.Message)
		if message == "" {
			message = fmt.Sprintf("last.fm error %d", response.Error)
		}

		switch response.Error {
		case 6, 7:
			return nil, nil
		case 29:
			return nil, &lastfmRateLimitError{body: message}
		default:
			return nil, fmt.Errorf("last.fm tag lookup failed: error=%d message=%s", response.Error, message)
		}
	}

	var entity lastfmInfoEntityTags
	switch entityType {
	case model.SpotifyEnrichmentTypeTrack:
		entity = response.Track
	case model.SpotifyEnrichmentTypeAlbum:
		entity = response.Album
	case model.SpotifyEnrichmentTypeArtist:
		entity = response.Artist
	default:
		return nil, fmt.Errorf("unsupported Last.fm tag entity type %q", entityType)
	}

	if tags := normaliseLastFMTags(entity.TopTags.Tag); len(tags) > 0 {
		return tags, nil
	}
	return normaliseLastFMTags(entity.Tags.Tag), nil
}

func normaliseLastFMTags(input []lastfmTag) []lastfmTag {
	tags := make([]lastfmTag, 0, len(input))
	seen := map[string]struct{}{}
	for _, tag := range input {
		tag.Name = strings.TrimSpace(tag.Name)
		tag.URL = strings.TrimSpace(tag.URL)
		if tag.Name == "" {
			continue
		}

		key := strings.ToLower(tag.Name)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		tags = append(tags, tag)
	}

	return tags
}

func limitLastFMTags(tags []lastfmTag, limit int) []lastfmTag {
	if limit < 1 || len(tags) <= limit {
		return tags
	}
	return tags[:limit]
}

func (s *Store) persistLastFMTags(ctx context.Context, target lastfmTagTarget, tags []lastfmTag) (int, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	for _, tag := range tags {
		genreID, err := upsertLastFMGenre(ctx, tx, tag)
		if err != nil {
			return 0, err
		}
		if err := insertLastFMGenreLink(ctx, tx, target, genreID, int(tag.Count)); err != nil {
			return 0, err
		}
	}

	status := lastfmTagStatusTagged
	if len(tags) == 0 {
		status = lastfmTagStatusNoTags
	}
	if err := s.recordLastFMTagAttempt(ctx, tx, target, status, len(tags), ""); err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return len(tags), nil
}

func lastfmTagNames(tags []lastfmTag) []string {
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		name := strings.TrimSpace(tag.Name)
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}

func upsertLastFMGenre(ctx context.Context, tx pgx.Tx, tag lastfmTag) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `
		INSERT INTO music.genres (name, url)
		VALUES ($1, $2::text)
		ON CONFLICT (lower(name))
		DO UPDATE SET
			url = coalesce(excluded.url, music.genres.url),
			updated_at = now()
		RETURNING id
	`, tag.Name, nullableString(tag.URL)).Scan(&id)
	return id, err
}

func insertLastFMGenreLink(ctx context.Context, tx pgx.Tx, target lastfmTagTarget, genreID int64, weight int) error {
	switch target.EntityType {
	case model.SpotifyEnrichmentTypeArtist:
		_, err := tx.Exec(ctx, `
			INSERT INTO music.artist_genres (artist_id, genre_id, source, weight)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (artist_id, genre_id, source)
			DO UPDATE SET
				weight = excluded.weight,
				updated_at = now()
		`, target.EntityID, genreID, lastfmTagSource, weight)
		return err
	case model.SpotifyEnrichmentTypeAlbum:
		_, err := tx.Exec(ctx, `
			INSERT INTO music.album_genres (album_id, genre_id, source, weight)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (album_id, genre_id, source)
			DO UPDATE SET
				weight = excluded.weight,
				updated_at = now()
		`, target.EntityID, genreID, lastfmTagSource, weight)
		return err
	case model.SpotifyEnrichmentTypeTrack:
		_, err := tx.Exec(ctx, `
			INSERT INTO music.track_genres (track_id, genre_id, source, weight)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (track_id, genre_id, source)
			DO UPDATE SET
				weight = excluded.weight,
				updated_at = now()
		`, target.EntityID, genreID, lastfmTagSource, weight)
		return err
	default:
		return fmt.Errorf("unsupported Last.fm tag entity type %q", target.EntityType)
	}
}

func (s *Store) recordLastFMTagAttempt(ctx context.Context, executor lastfmTagExecutor, target lastfmTagTarget, status string, tagsFound int, lastError string) error {
	if executor == nil {
		return nil
	}

	_, err := executor.Exec(ctx, `
		INSERT INTO music.lastfm_tag_enrichment_attempts (
			entity_type,
			entity_id,
			status,
			attempts,
			tags_found,
			last_error,
			searched_at
		)
		VALUES ($1, $2, $3, 1, $4, $5::text, now())
		ON CONFLICT (entity_type, entity_id)
		DO UPDATE SET
			status = excluded.status,
			attempts = music.lastfm_tag_enrichment_attempts.attempts + 1,
			tags_found = excluded.tags_found,
			last_error = excluded.last_error,
			searched_at = excluded.searched_at,
			updated_at = now()
	`, target.EntityType, target.EntityID, status, tagsFound, nullableString(lastError))
	return err
}

func (s *Store) lastfmRateLimitRetryAfter(ctx context.Context) *int {
	if s.Redis == nil {
		return nil
	}

	ttl, err := s.Redis.TTL(ctx, lastfmTagsRateLimitKey).Result()
	if err != nil || ttl <= 0 {
		return nil
	}

	seconds := int((ttl + time.Second - time.Nanosecond) / time.Second)
	return &seconds
}

func (s *Store) rememberLastFMRateLimit(ctx context.Context) int {
	cooldown := lastfmRateLimitCooldown()
	if s.Redis == nil {
		return int(cooldown / time.Second)
	}

	if err := s.Redis.Set(ctx, lastfmTagsRateLimitKey, "1", cooldown).Err(); err != nil && s.Logger != nil {
		s.Logger.Printf("Redis Last.fm rate-limit cooldown write skipped: %v", err)
	}

	return int(cooldown / time.Second)
}

func lastfmTagRateLimitedResult(result *model.LastFMTagEnrichmentResult, reason string, retryAfterSeconds int) *model.LastFMTagEnrichmentResult {
	result.Status = lastfmTagStatusRateLimited
	result.RateLimited = true
	result.RateLimitReason = reason
	result.RateLimitRetryAfterSeconds = &retryAfterSeconds
	result.Notes = reason
	return result
}

func lastfmRateLimitCooldown() time.Duration {
	raw := strings.TrimSpace(os.Getenv("LASTFM_RATE_LIMIT_COOLDOWN_SECONDS"))
	if raw == "" {
		return lastfmDefaultRateLimitCooldown
	}

	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds < 1 {
		return lastfmDefaultRateLimitCooldown
	}

	return time.Duration(seconds) * time.Second
}

func lastfmTagLimit() int {
	raw := strings.TrimSpace(os.Getenv("LASTFM_TAG_LIMIT"))
	if raw == "" {
		return lastfmDefaultTagLimit
	}

	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 {
		return lastfmDefaultTagLimit
	}
	if limit > 50 {
		return 50
	}
	return limit
}

func lastfmTagPollingEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LASTFM_TAGS_ON_POLL"))) {
	case "", "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
