package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
	"github.com/lesi97/lesi.dev/internal/httpapi"
)

const (
	spotifySearchRateLimitKey = "spotify:rate-limit:search"

	spotifyEnrichmentStatusMatched       = "matched"
	spotifyEnrichmentStatusNoMatch       = "no_match"
	spotifyEnrichmentStatusLowConfidence = "low_confidence"
	spotifyEnrichmentStatusError         = "error"
	spotifyEnrichmentStatusRateLimited   = "rate_limited"

	spotifyTrackMatchThreshold  = 0.90
	spotifyAlbumMatchThreshold  = 0.90
	spotifyArtistMatchThreshold = 0.95
)

type spotifySearchResponse struct {
	Tracks  spotifyTrackSearchPage  `json:"tracks"`
	Albums  spotifyAlbumSearchPage  `json:"albums"`
	Artists spotifyArtistSearchPage `json:"artists"`
}

type spotifyTrackSearchPage struct {
	Items []spotifyTrack `json:"items"`
}

type spotifyAlbumSearchPage struct {
	Items []spotifyAlbum `json:"items"`
}

type spotifyArtistSearchPage struct {
	Items []spotifyArtist `json:"items"`
}

type spotifyEnrichmentCandidate struct {
	EntityType string
	ID         int64
	Name       string
	ArtistName *string
	AlbumName  *string
}

type spotifyEnrichmentMatch struct {
	Status     string
	Confidence *float64
	SpotifyID  *string
	SpotifyURL *string
	ImageURL   *string
	Notes      string
}

type spotifyEnrichmentExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func (s *Store) EnrichSpotifyMetadata(ctx context.Context, input model.SpotifyEnrichmentInput) (*model.SpotifyEnrichmentResult, error) {
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

	result := &model.SpotifyEnrichmentResult{
		EntityType: entityType,
		Status:     "pending",
	}

	if retryAfter := s.spotifyRateLimitRetryAfter(ctx, spotifySearchRateLimitKey); retryAfter != nil {
		return spotifyRateLimitedEnrichmentResult(result, "spotify search cooldown active", *retryAfter), nil
	}
	if retryAfter := s.spotifyRateLimitRetryAfter(ctx, spotifyTokenRateLimitKey); retryAfter != nil {
		return spotifyRateLimitedEnrichmentResult(result, "spotify token cooldown active", *retryAfter), nil
	}

	candidate, ok, err := s.nextSpotifyEnrichmentCandidate(ctx, entityType)
	if err != nil {
		return nil, err
	}
	if !ok {
		result.Status = "complete"
		result.Notes = "no eligible " + entityType + " rows found"
		return result, nil
	}

	result.EntityID = &candidate.ID
	result.EntityName = candidate.Name

	auth, err := s.spotifyAPIAuth(ctx)
	if err != nil {
		if rateLimited := s.spotifyRateLimitedEnrichmentResultFromError(ctx, result, candidate, spotifyTokenRateLimitKey, err); rateLimited != nil {
			return rateLimited, nil
		}
		return nil, err
	}

	search, err := s.fetchSpotifySearch(ctx, auth.APIURL, auth.AccessToken, candidate)
	if errors.Is(err, errSpotifyUnauthorized) {
		accessToken, refreshErr := s.refreshSpotifyAccessToken(ctx, auth.Details)
		if refreshErr != nil {
			if rateLimited := s.spotifyRateLimitedEnrichmentResultFromError(ctx, result, candidate, spotifyTokenRateLimitKey, refreshErr); rateLimited != nil {
				return rateLimited, nil
			}
			return nil, refreshErr
		}
		auth.AccessToken = accessToken
		search, err = s.fetchSpotifySearch(ctx, auth.APIURL, auth.AccessToken, candidate)
	}
	if err != nil {
		if rateLimited := s.spotifyRateLimitedEnrichmentResultFromError(ctx, result, candidate, spotifySearchRateLimitKey, err); rateLimited != nil {
			return rateLimited, nil
		}

		match := spotifyEnrichmentMatch{
			Status: spotifyEnrichmentStatusError,
			Notes:  err.Error(),
		}
		if recordErr := s.recordSpotifyEnrichmentAttempt(ctx, s.DB, candidate, match, err.Error()); recordErr != nil {
			return nil, recordErr
		}
		return spotifyEnrichmentResultFromMatch(result, match), nil
	}

	match := bestSpotifyEnrichmentMatch(candidate, search)
	if match.Status == spotifyEnrichmentStatusMatched {
		if err := s.updateSpotifyEnrichedEntity(ctx, candidate, match); err != nil {
			errorMatch := match
			errorMatch.Status = spotifyEnrichmentStatusError
			errorMatch.Notes = err.Error()
			if recordErr := s.recordSpotifyEnrichmentAttempt(ctx, s.DB, candidate, errorMatch, err.Error()); recordErr != nil {
				return nil, recordErr
			}
			return spotifyEnrichmentResultFromMatch(result, errorMatch), nil
		}
	} else if err := s.recordSpotifyEnrichmentAttempt(ctx, s.DB, candidate, match, ""); err != nil {
		return nil, err
	}

	return spotifyEnrichmentResultFromMatch(result, match), nil
}

func isSpotifyEnrichmentType(entityType string) bool {
	switch entityType {
	case model.SpotifyEnrichmentTypeTrack, model.SpotifyEnrichmentTypeAlbum, model.SpotifyEnrichmentTypeArtist:
		return true
	default:
		return false
	}
}

func (s *Store) nextSpotifyEnrichmentCandidate(ctx context.Context, entityType string) (spotifyEnrichmentCandidate, bool, error) {
	switch entityType {
	case model.SpotifyEnrichmentTypeTrack:
		return s.nextSpotifyTrackEnrichmentCandidate(ctx)
	case model.SpotifyEnrichmentTypeAlbum:
		return s.nextSpotifyAlbumEnrichmentCandidate(ctx)
	case model.SpotifyEnrichmentTypeArtist:
		return s.nextSpotifyArtistEnrichmentCandidate(ctx)
	default:
		return spotifyEnrichmentCandidate{}, false, fmt.Errorf("type must be track, album, or artist")
	}
}

func (s *Store) nextSpotifyTrackEnrichmentCandidate(ctx context.Context) (spotifyEnrichmentCandidate, bool, error) {
	candidate := spotifyEnrichmentCandidate{EntityType: model.SpotifyEnrichmentTypeTrack}
	err := s.DB.QueryRow(ctx, `
		SELECT
			t.id,
			t.name,
			a.name,
			al.name
		FROM music.tracks t
		JOIN music.artists a ON a.id = t.artist_id
		LEFT JOIN music.albums al ON al.id = t.album_id
		LEFT JOIN music.spotify_enrichment_attempts e
			ON e.entity_type = 'track'
			AND e.entity_id = t.id
		WHERE (
			t.spotify_id IS NULL
			OR t.url IS NULL
			OR t.url NOT ILIKE 'https://open.spotify.com/track/%'
			OR t.image_url IS NULL
		)
		AND (
			e.id IS NULL
			OR (e.status = 'error' AND e.updated_at < now() - interval '6 hours')
			OR (e.status = 'rate_limited' AND e.updated_at < now() - interval '15 minutes')
		)
		ORDER BY
			coalesce(e.attempts, 0),
			coalesce(e.updated_at, 'epoch'::timestamptz),
			t.id
		LIMIT 1
	`).Scan(&candidate.ID, &candidate.Name, &candidate.ArtistName, &candidate.AlbumName)

	if errors.Is(err, pgx.ErrNoRows) {
		return spotifyEnrichmentCandidate{}, false, nil
	}
	if err != nil {
		return spotifyEnrichmentCandidate{}, false, err
	}

	return candidate, true, nil
}

func (s *Store) nextSpotifyAlbumEnrichmentCandidate(ctx context.Context) (spotifyEnrichmentCandidate, bool, error) {
	candidate := spotifyEnrichmentCandidate{EntityType: model.SpotifyEnrichmentTypeAlbum}
	err := s.DB.QueryRow(ctx, `
		SELECT
			al.id,
			al.name,
			a.name
		FROM music.albums al
		JOIN music.artists a ON a.id = al.artist_id
		LEFT JOIN music.spotify_enrichment_attempts e
			ON e.entity_type = 'album'
			AND e.entity_id = al.id
		WHERE (
			al.spotify_id IS NULL
			OR al.url IS NULL
			OR al.url NOT ILIKE 'https://open.spotify.com/album/%'
			OR al.image_url IS NULL
		)
		AND (
			e.id IS NULL
			OR (e.status = 'error' AND e.updated_at < now() - interval '6 hours')
			OR (e.status = 'rate_limited' AND e.updated_at < now() - interval '15 minutes')
		)
		ORDER BY
			coalesce(e.attempts, 0),
			coalesce(e.updated_at, 'epoch'::timestamptz),
			al.id
		LIMIT 1
	`).Scan(&candidate.ID, &candidate.Name, &candidate.ArtistName)

	if errors.Is(err, pgx.ErrNoRows) {
		return spotifyEnrichmentCandidate{}, false, nil
	}
	if err != nil {
		return spotifyEnrichmentCandidate{}, false, err
	}

	return candidate, true, nil
}

func (s *Store) nextSpotifyArtistEnrichmentCandidate(ctx context.Context) (spotifyEnrichmentCandidate, bool, error) {
	candidate := spotifyEnrichmentCandidate{EntityType: model.SpotifyEnrichmentTypeArtist}
	err := s.DB.QueryRow(ctx, `
		SELECT
			a.id,
			a.name
		FROM music.artists a
		LEFT JOIN music.spotify_enrichment_attempts e
			ON e.entity_type = 'artist'
			AND e.entity_id = a.id
		WHERE (
			a.spotify_id IS NULL
			OR a.url IS NULL
			OR a.url NOT ILIKE 'https://open.spotify.com/artist/%'
			OR a.image_url IS NULL
		)
		AND (
			e.id IS NULL
			OR (e.status = 'error' AND e.updated_at < now() - interval '6 hours')
			OR (e.status = 'rate_limited' AND e.updated_at < now() - interval '15 minutes')
		)
		ORDER BY
			coalesce(e.attempts, 0),
			coalesce(e.updated_at, 'epoch'::timestamptz),
			a.id
		LIMIT 1
	`).Scan(&candidate.ID, &candidate.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		return spotifyEnrichmentCandidate{}, false, nil
	}
	if err != nil {
		return spotifyEnrichmentCandidate{}, false, err
	}

	return candidate, true, nil
}

func (s *Store) fetchSpotifySearch(ctx context.Context, apiURL string, accessToken string, candidate spotifyEnrichmentCandidate) (spotifySearchResponse, error) {
	endpoint, err := url.Parse(strings.TrimRight(apiURL, "/") + "/search")
	if err != nil {
		return spotifySearchResponse{}, err
	}

	query := endpoint.Query()
	query.Set("type", candidate.EntityType)
	query.Set("limit", "10")
	query.Set("q", spotifySearchQuery(candidate))
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
		return spotifySearchResponse{}, err
	}
	if statusCode == http.StatusUnauthorized {
		return spotifySearchResponse{}, errSpotifyUnauthorized
	}
	if statusCode == http.StatusTooManyRequests {
		return spotifySearchResponse{}, &spotifyRateLimitError{
			endpoint: "search",
			body:     string(body),
		}
	}
	if statusCode < 200 || statusCode >= 300 {
		return spotifySearchResponse{}, fmt.Errorf("spotify search failed: status=%d body=%s", statusCode, string(body))
	}

	var response spotifySearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return spotifySearchResponse{}, err
	}

	return response, nil
}

func spotifySearchQuery(candidate spotifyEnrichmentCandidate) string {
	switch candidate.EntityType {
	case model.SpotifyEnrichmentTypeTrack:
		parts := []string{spotifySearchTerm("track", candidate.Name)}
		if artist := stringValue(candidate.ArtistName); artist != "" {
			parts = append(parts, spotifySearchTerm("artist", artist))
		}
		if album := stringValue(candidate.AlbumName); album != "" {
			parts = append(parts, spotifySearchTerm("album", album))
		}
		return strings.Join(parts, " ")
	case model.SpotifyEnrichmentTypeAlbum:
		parts := []string{spotifySearchTerm("album", candidate.Name)}
		if artist := stringValue(candidate.ArtistName); artist != "" {
			parts = append(parts, spotifySearchTerm("artist", artist))
		}
		return strings.Join(parts, " ")
	case model.SpotifyEnrichmentTypeArtist:
		return spotifySearchTerm("artist", candidate.Name)
	default:
		return candidate.Name
	}
}

func spotifySearchTerm(field string, value string) string {
	value = strings.ReplaceAll(strings.TrimSpace(value), `"`, "")
	return fmt.Sprintf(`%s:"%s"`, field, value)
}

func bestSpotifyEnrichmentMatch(candidate spotifyEnrichmentCandidate, search spotifySearchResponse) spotifyEnrichmentMatch {
	switch candidate.EntityType {
	case model.SpotifyEnrichmentTypeTrack:
		return bestSpotifyTrackEnrichmentMatch(candidate, search.Tracks.Items)
	case model.SpotifyEnrichmentTypeAlbum:
		return bestSpotifyAlbumEnrichmentMatch(candidate, search.Albums.Items)
	case model.SpotifyEnrichmentTypeArtist:
		return bestSpotifyArtistEnrichmentMatch(candidate, search.Artists.Items)
	default:
		return spotifyEnrichmentMatch{
			Status: spotifyEnrichmentStatusError,
			Notes:  "unsupported enrichment type",
		}
	}
}

func bestSpotifyTrackEnrichmentMatch(candidate spotifyEnrichmentCandidate, tracks []spotifyTrack) spotifyEnrichmentMatch {
	if len(tracks) == 0 {
		return spotifyEnrichmentMatch{Status: spotifyEnrichmentStatusNoMatch, Notes: "spotify search returned no tracks"}
	}

	var best spotifyEnrichmentMatch
	for _, track := range tracks {
		spotifyURL := strings.TrimSpace(track.ExternalURLs.Spotify)
		if strings.TrimSpace(track.ID) == "" || spotifyURL == "" {
			continue
		}

		confidence := spotifyTrackConfidence(candidate, track)
		if best.Confidence != nil && confidence <= *best.Confidence {
			continue
		}

		spotifyID := strings.TrimSpace(track.ID)
		imageURL := bestSpotifyImageURL(track.Album.Images)
		best = spotifyEnrichmentMatch{
			Status:     spotifyEnrichmentStatusLowConfidence,
			Confidence: &confidence,
			SpotifyID:  &spotifyID,
			SpotifyURL: &spotifyURL,
			ImageURL:   optionalNonBlankString(imageURL),
			Notes:      "best spotify track match was below confidence threshold",
		}
	}

	if best.Confidence == nil {
		return spotifyEnrichmentMatch{Status: spotifyEnrichmentStatusNoMatch, Notes: "spotify search returned no usable tracks"}
	}
	if *best.Confidence >= spotifyTrackMatchThreshold {
		best.Status = spotifyEnrichmentStatusMatched
		best.Notes = "matched spotify track"
	}
	return best
}

func bestSpotifyAlbumEnrichmentMatch(candidate spotifyEnrichmentCandidate, albums []spotifyAlbum) spotifyEnrichmentMatch {
	if len(albums) == 0 {
		return spotifyEnrichmentMatch{Status: spotifyEnrichmentStatusNoMatch, Notes: "spotify search returned no albums"}
	}

	var best spotifyEnrichmentMatch
	for _, album := range albums {
		spotifyURL := strings.TrimSpace(album.ExternalURLs.Spotify)
		if strings.TrimSpace(album.ID) == "" || spotifyURL == "" {
			continue
		}

		confidence := spotifyAlbumConfidence(candidate, album)
		if best.Confidence != nil && confidence <= *best.Confidence {
			continue
		}

		spotifyID := strings.TrimSpace(album.ID)
		imageURL := bestSpotifyImageURL(album.Images)
		best = spotifyEnrichmentMatch{
			Status:     spotifyEnrichmentStatusLowConfidence,
			Confidence: &confidence,
			SpotifyID:  &spotifyID,
			SpotifyURL: &spotifyURL,
			ImageURL:   optionalNonBlankString(imageURL),
			Notes:      "best spotify album match was below confidence threshold",
		}
	}

	if best.Confidence == nil {
		return spotifyEnrichmentMatch{Status: spotifyEnrichmentStatusNoMatch, Notes: "spotify search returned no usable albums"}
	}
	if *best.Confidence >= spotifyAlbumMatchThreshold {
		best.Status = spotifyEnrichmentStatusMatched
		best.Notes = "matched spotify album"
	}
	return best
}

func bestSpotifyArtistEnrichmentMatch(candidate spotifyEnrichmentCandidate, artists []spotifyArtist) spotifyEnrichmentMatch {
	if len(artists) == 0 {
		return spotifyEnrichmentMatch{Status: spotifyEnrichmentStatusNoMatch, Notes: "spotify search returned no artists"}
	}

	var best spotifyEnrichmentMatch
	for _, artist := range artists {
		spotifyURL := strings.TrimSpace(artist.ExternalURLs.Spotify)
		if strings.TrimSpace(artist.ID) == "" || spotifyURL == "" {
			continue
		}

		confidence := nameMatchScore(candidate.Name, artist.Name)
		if best.Confidence != nil && confidence <= *best.Confidence {
			continue
		}

		spotifyID := strings.TrimSpace(artist.ID)
		imageURL := bestSpotifyImageURL(artist.Images)
		best = spotifyEnrichmentMatch{
			Status:     spotifyEnrichmentStatusLowConfidence,
			Confidence: &confidence,
			SpotifyID:  &spotifyID,
			SpotifyURL: &spotifyURL,
			ImageURL:   optionalNonBlankString(imageURL),
			Notes:      "best spotify artist match was below confidence threshold",
		}
	}

	if best.Confidence == nil {
		return spotifyEnrichmentMatch{Status: spotifyEnrichmentStatusNoMatch, Notes: "spotify search returned no usable artists"}
	}
	if *best.Confidence >= spotifyArtistMatchThreshold {
		best.Status = spotifyEnrichmentStatusMatched
		best.Notes = "matched spotify artist"
	}
	return best
}

func spotifyTrackConfidence(candidate spotifyEnrichmentCandidate, track spotifyTrack) float64 {
	trackScore := nameMatchScore(candidate.Name, track.Name)
	artistScore := nameMatchScore(stringValue(candidate.ArtistName), primarySpotifyArtistName(track.Artists))
	albumName := stringValue(candidate.AlbumName)
	if albumName == "" {
		return 0.65*trackScore + 0.35*artistScore
	}

	albumScore := nameMatchScore(albumName, track.Album.Name)
	return 0.60*trackScore + 0.30*artistScore + 0.10*albumScore
}

func spotifyAlbumConfidence(candidate spotifyEnrichmentCandidate, album spotifyAlbum) float64 {
	albumScore := nameMatchScore(candidate.Name, album.Name)
	artistScore := nameMatchScore(stringValue(candidate.ArtistName), primarySpotifyArtistName(album.Artists))
	return 0.65*albumScore + 0.35*artistScore
}

func primarySpotifyArtistName(artists []spotifyArtist) string {
	if len(artists) == 0 {
		return ""
	}
	return artists[0].Name
}

func nameMatchScore(want string, got string) float64 {
	want = normaliseSpotifyComparableText(want)
	got = normaliseSpotifyComparableText(got)
	if want == "" || got == "" {
		return 0
	}
	if want == got {
		return 1
	}
	if len(want) >= 4 && len(got) >= 4 && (strings.Contains(want, got) || strings.Contains(got, want)) {
		return 0.80
	}
	return 0
}

func normaliseSpotifyComparableText(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastSpace := true

	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastSpace = false
			continue
		}
		if !lastSpace {
			b.WriteByte(' ')
			lastSpace = true
		}
	}

	return strings.TrimSpace(b.String())
}

func (s *Store) updateSpotifyEnrichedEntity(ctx context.Context, candidate spotifyEnrichmentCandidate, match spotifyEnrichmentMatch) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	switch candidate.EntityType {
	case model.SpotifyEnrichmentTypeTrack:
		_, err = tx.Exec(ctx, `
			UPDATE music.tracks
			SET
				spotify_id = $1,
				url = $2,
				image_url = coalesce($3::text, image_url),
				updated_at = now()
			WHERE id = $4
		`, optionalString(match.SpotifyID), optionalString(match.SpotifyURL), optionalString(match.ImageURL), candidate.ID)
	case model.SpotifyEnrichmentTypeAlbum:
		_, err = tx.Exec(ctx, `
			UPDATE music.albums
			SET
				spotify_id = $1,
				url = $2,
				image_url = coalesce($3::text, image_url),
				updated_at = now()
			WHERE id = $4
		`, optionalString(match.SpotifyID), optionalString(match.SpotifyURL), optionalString(match.ImageURL), candidate.ID)
	case model.SpotifyEnrichmentTypeArtist:
		_, err = tx.Exec(ctx, `
			UPDATE music.artists
			SET
				spotify_id = $1,
				url = $2,
				image_url = coalesce($3::text, image_url),
				updated_at = now()
			WHERE id = $4
		`, optionalString(match.SpotifyID), optionalString(match.SpotifyURL), optionalString(match.ImageURL), candidate.ID)
	default:
		err = fmt.Errorf("type must be track, album, or artist")
	}
	if err != nil {
		return err
	}

	if err := s.recordSpotifyEnrichmentAttempt(ctx, tx, candidate, match, ""); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) recordSpotifyEnrichmentAttempt(ctx context.Context, executor spotifyEnrichmentExecutor, candidate spotifyEnrichmentCandidate, match spotifyEnrichmentMatch, lastError string) error {
	_, err := executor.Exec(ctx, `
		INSERT INTO music.spotify_enrichment_attempts (
			entity_type,
			entity_id,
			status,
			attempts,
			confidence,
			spotify_id,
			spotify_url,
			image_url,
			notes,
			last_error,
			searched_at
		)
		VALUES ($1, $2, $3, 1, $4::numeric, $5::text, $6::text, $7::text, $8::text, $9::text, now())
		ON CONFLICT (entity_type, entity_id)
		DO UPDATE SET
			status = excluded.status,
			attempts = music.spotify_enrichment_attempts.attempts + 1,
			confidence = excluded.confidence,
			spotify_id = excluded.spotify_id,
			spotify_url = excluded.spotify_url,
			image_url = excluded.image_url,
			notes = excluded.notes,
			last_error = excluded.last_error,
			searched_at = excluded.searched_at,
			updated_at = now()
	`, candidate.EntityType, candidate.ID, match.Status, optionalFloat64(match.Confidence), optionalString(match.SpotifyID), optionalString(match.SpotifyURL), optionalString(match.ImageURL), nullableString(match.Notes), nullableString(lastError))
	return err
}

func spotifyEnrichmentResultFromMatch(result *model.SpotifyEnrichmentResult, match spotifyEnrichmentMatch) *model.SpotifyEnrichmentResult {
	result.Status = match.Status
	result.Matched = match.Status == spotifyEnrichmentStatusMatched
	result.Updated = match.Status == spotifyEnrichmentStatusMatched
	result.Confidence = match.Confidence
	result.SpotifyID = match.SpotifyID
	result.SpotifyURL = match.SpotifyURL
	result.ImageURL = match.ImageURL
	result.Notes = match.Notes
	return result
}

func (s *Store) spotifyRateLimitedEnrichmentResultFromError(ctx context.Context, result *model.SpotifyEnrichmentResult, candidate spotifyEnrichmentCandidate, cacheKey string, err error) *model.SpotifyEnrichmentResult {
	var rateLimitErr *spotifyRateLimitError
	if !errors.As(err, &rateLimitErr) {
		return nil
	}

	retryAfter := s.rememberSpotifyRateLimit(ctx, cacheKey)
	match := spotifyEnrichmentMatch{
		Status: spotifyEnrichmentStatusRateLimited,
		Notes:  rateLimitErr.Error(),
	}
	if recordErr := s.recordSpotifyEnrichmentAttempt(ctx, s.DB, candidate, match, rateLimitErr.Error()); recordErr != nil && s.Logger != nil {
		s.Logger.Printf("Spotify enrichment rate-limit attempt record skipped: %v", recordErr)
	}

	return spotifyRateLimitedEnrichmentResult(result, rateLimitErr.Error(), retryAfter)
}

func spotifyRateLimitedEnrichmentResult(result *model.SpotifyEnrichmentResult, reason string, retryAfterSeconds int) *model.SpotifyEnrichmentResult {
	result.Status = spotifyEnrichmentStatusRateLimited
	result.RateLimited = true
	result.RateLimitReason = reason
	result.RateLimitRetryAfterSeconds = &retryAfterSeconds
	result.Notes = reason
	return result
}

func optionalNonBlankString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func nullableString(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}

func optionalFloat64(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}
