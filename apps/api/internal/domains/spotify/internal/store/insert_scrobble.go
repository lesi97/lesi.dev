package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
)

func (s *Store) InsertScrobble(ctx context.Context, input model.ScrobbleInput) (*model.ScrobbleResult, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	artistID, artistCreated, err := upsertArtist(ctx, tx, input.Artist)
	if err != nil {
		return nil, err
	}

	var albumID *int64
	var albumCreated bool
	if input.Album != nil {
		id, created, err := upsertAlbum(ctx, tx, artistID, *input.Album)
		if err != nil {
			return nil, err
		}
		albumID = &id
		albumCreated = created
	}

	trackID, trackCreated, err := upsertTrack(ctx, tx, artistID, albumID, input.Track)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO music.scrobbles (
			scrobbled_at,
			artist_id,
			album_id,
			track_id,
			source,
			raw_payload
		)
		VALUES ($1, $2, $3::bigint, $4, $5, $6::jsonb)
		ON CONFLICT (scrobbled_at, artist_id, track_id)
		DO UPDATE SET
			album_id = excluded.album_id,
			source = excluded.source,
			raw_payload = excluded.raw_payload,
			updated_at = now()
		RETURNING id, (xmax = 0) AS created
	`

	var scrobbleID int64
	var scrobbleCreated bool
	err = tx.QueryRow(
		ctx,
		query,
		input.ScrobbledAt.UTC(),
		artistID,
		optionalInt64(albumID),
		trackID,
		input.Source,
		input.RawPayload,
	).Scan(&scrobbleID, &scrobbleCreated)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	result := &model.ScrobbleResult{
		ID:              scrobbleID,
		ArtistID:        artistID,
		AlbumID:         albumID,
		TrackID:         trackID,
		ScrobbledAt:     input.ScrobbledAt.UTC().Format(time.RFC3339),
		ArtistCreated:   artistCreated,
		AlbumCreated:    albumCreated,
		TrackCreated:    trackCreated,
		ScrobbleCreated: scrobbleCreated,
		ArtistURL:       input.Artist.URL,
		TrackURL:        input.Track.URL,
		ArtistImageURL:  input.Artist.ImageURL,
		TrackImageURL:   input.Track.ImageURL,
		ArtistSpotifyID: input.Artist.SpotifyID,
		TrackSpotifyID:  input.Track.SpotifyID,
	}
	if input.Album != nil {
		result.AlbumURL = input.Album.URL
		result.AlbumImageURL = input.Album.ImageURL
		result.AlbumSpotifyID = input.Album.SpotifyID
	}

	return result, nil
}

func upsertArtist(ctx context.Context, tx pgx.Tx, entity model.ScrobbleEntityInput) (int64, bool, error) {
	spotifyID := optionalString(entity.SpotifyID)
	url := optionalString(entity.URL)
	imageURL := optionalString(entity.ImageURL)

	id, ok, err := selectID(ctx, tx, `
		SELECT id
		FROM music.artists
		WHERE ($2::text IS NOT NULL AND spotify_id = $2::text)
			OR (
				($2::text IS NULL OR spotify_id IS NULL)
				AND lower(name) = lower($1)
			)
		ORDER BY
			CASE
				WHEN $2::text IS NOT NULL AND spotify_id = $2::text THEN 0
				ELSE 1
			END
		LIMIT 1
	`, entity.Name, spotifyID)
	if err != nil {
		return 0, false, err
	}

	if ok {
		_, err = tx.Exec(ctx, `
			UPDATE music.artists
			SET
				name = $1,
				spotify_id = coalesce(spotify_id, $2::text),
				url = coalesce($3::text, url),
				image_url = coalesce($4::text, image_url),
				updated_at = now()
			WHERE id = $5
		`, entity.Name, spotifyID, url, imageURL, id)
		return id, false, err
	}

	id, err = insertID(ctx, tx, `
		INSERT INTO music.artists (name, spotify_id, url, image_url)
		VALUES ($1, $2::text, $3::text, $4::text)
		RETURNING id
	`, entity.Name, spotifyID, url, imageURL)
	return id, true, err
}

func upsertAlbum(ctx context.Context, tx pgx.Tx, artistID int64, entity model.ScrobbleEntityInput) (int64, bool, error) {
	spotifyID := optionalString(entity.SpotifyID)
	url := optionalString(entity.URL)
	imageURL := optionalString(entity.ImageURL)

	id, ok, err := selectID(ctx, tx, `
		SELECT id
		FROM music.albums
		WHERE ($3::text IS NOT NULL AND spotify_id = $3::text)
			OR (
				($3::text IS NULL OR spotify_id IS NULL)
				AND artist_id = $1
				AND lower(name) = lower($2)
			)
		ORDER BY
			CASE
				WHEN $3::text IS NOT NULL AND spotify_id = $3::text THEN 0
				ELSE 1
			END
		LIMIT 1
	`, artistID, entity.Name, spotifyID)
	if err != nil {
		return 0, false, err
	}

	if ok {
		_, err = tx.Exec(ctx, `
			UPDATE music.albums
			SET
				artist_id = $1,
				name = $2,
				spotify_id = coalesce(spotify_id, $3::text),
				url = coalesce($4::text, url),
				image_url = coalesce($5::text, image_url),
				updated_at = now()
			WHERE id = $6
		`, artistID, entity.Name, spotifyID, url, imageURL, id)
		return id, false, err
	}

	id, err = insertID(ctx, tx, `
		INSERT INTO music.albums (artist_id, name, spotify_id, url, image_url)
		VALUES ($1, $2, $3::text, $4::text, $5::text)
		RETURNING id
	`, artistID, entity.Name, spotifyID, url, imageURL)
	return id, true, err
}

func upsertTrack(ctx context.Context, tx pgx.Tx, artistID int64, albumID *int64, entity model.ScrobbleEntityInput) (int64, bool, error) {
	spotifyID := optionalString(entity.SpotifyID)
	url := optionalString(entity.URL)
	imageURL := optionalString(entity.ImageURL)
	album := optionalInt64(albumID)

	id, ok, err := selectID(ctx, tx, `
		SELECT id
		FROM music.tracks
		WHERE ($4::text IS NOT NULL AND spotify_id = $4::text)
			OR (
				($4::text IS NULL OR spotify_id IS NULL)
				AND artist_id = $1
				AND coalesce(album_id, 0) = coalesce($2::bigint, 0)
				AND lower(name) = lower($3)
			)
		ORDER BY
			CASE
				WHEN $4::text IS NOT NULL AND spotify_id = $4::text THEN 0
				ELSE 1
			END
		LIMIT 1
	`, artistID, album, entity.Name, spotifyID)
	if err != nil {
		return 0, false, err
	}

	if ok {
		_, err = tx.Exec(ctx, `
			UPDATE music.tracks
			SET
				artist_id = $1,
				album_id = coalesce($2::bigint, album_id),
				name = $3,
				spotify_id = coalesce(spotify_id, $4::text),
				url = coalesce($5::text, url),
				image_url = coalesce($6::text, image_url),
				updated_at = now()
			WHERE id = $7
		`, artistID, album, entity.Name, spotifyID, url, imageURL, id)
		return id, false, err
	}

	id, err = insertID(ctx, tx, `
		INSERT INTO music.tracks (artist_id, album_id, name, spotify_id, url, image_url)
		VALUES ($1, $2::bigint, $3, $4::text, $5::text, $6::text)
		RETURNING id
	`, artistID, album, entity.Name, spotifyID, url, imageURL)
	return id, true, err
}

func selectID(ctx context.Context, tx pgx.Tx, query string, args ...any) (int64, bool, error) {
	var id int64
	err := tx.QueryRow(ctx, query, args...).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return id, true, nil
}

func insertID(ctx context.Context, tx pgx.Tx, query string, args ...any) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func optionalString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func optionalInt64(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}
