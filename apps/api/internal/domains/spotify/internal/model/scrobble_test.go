package model

import (
	"testing"
	"time"
)

func TestScrobbleRequestToInputAcceptsImageUrls(t *testing.T) {
	req := ScrobbleRequest{
		UTS:             UnixTimestamp{Value: 1786012626, Set: true},
		Artist:          " Hozier ",
		ArtistSpotifyID: "artist-123",
		ArtistImageURL:  "https://i.scdn.co/image/artist",
		Album:           "Wasteland, Baby!",
		AlbumSpotifyID:  "album-123",
		AlbumImageURL:   "https://i.scdn.co/image/album",
		Track:           "Talk",
		TrackSpotifyID:  "track-123",
		TrackImageURL:   "https://i.scdn.co/image/track",
		Source:          " spotify-poll ",
	}

	input, err := req.ToInput([]byte(`{"track":"Talk"}`))
	if err != nil {
		t.Fatalf("ToInput returned error: %v", err)
	}

	if got, want := input.ScrobbledAt, time.Unix(1786012626, 0).UTC(); !got.Equal(want) {
		t.Fatalf("scrobbled at = %s, want %s", got, want)
	}
	if got, want := input.Artist.Name, "Hozier"; got != want {
		t.Fatalf("artist = %q, want %q", got, want)
	}
	if got, want := valueOrEmpty(input.Artist.URL), spotifyArtistURL+"artist-123"; got != want {
		t.Fatalf("artist url = %q, want %q", got, want)
	}
	if got, want := valueOrEmpty(input.Artist.ImageURL), "https://i.scdn.co/image/artist"; got != want {
		t.Fatalf("artist image url = %q, want %q", got, want)
	}
	if input.Album == nil {
		t.Fatal("album is nil")
	}
	if got, want := valueOrEmpty(input.Album.URL), spotifyAlbumURL+"album-123"; got != want {
		t.Fatalf("album url = %q, want %q", got, want)
	}
	if got, want := valueOrEmpty(input.Album.ImageURL), "https://i.scdn.co/image/album"; got != want {
		t.Fatalf("album image url = %q, want %q", got, want)
	}
	if got, want := valueOrEmpty(input.Track.URL), spotifyTrackURL+"track-123"; got != want {
		t.Fatalf("track url = %q, want %q", got, want)
	}
	if got, want := valueOrEmpty(input.Track.ImageURL), "https://i.scdn.co/image/track"; got != want {
		t.Fatalf("track image url = %q, want %q", got, want)
	}
	if got, want := input.Source, "spotify-poll"; got != want {
		t.Fatalf("source = %q, want %q", got, want)
	}
}

func TestScrobbleRequestToInputRejectsInvalidImageUrl(t *testing.T) {
	req := ScrobbleRequest{
		UTS:           UnixTimestamp{Value: 1786012626, Set: true},
		Artist:        "Hozier",
		Track:         "Talk",
		TrackImageURL: "ftp://example.com/cover.jpg",
	}

	_, err := req.ToInput(nil)
	if err == nil {
		t.Fatal("ToInput returned nil error")
	}
	if got, want := err.Error(), "track_image_url must be an http or https URL"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
}

func TestScrobbleRequestToInputRequiresAlbumNameForAlbumImage(t *testing.T) {
	req := ScrobbleRequest{
		UTS:           UnixTimestamp{Value: 1786012626, Set: true},
		Artist:        "Hozier",
		AlbumImageURL: "https://i.scdn.co/image/album",
		Track:         "Talk",
	}

	_, err := req.ToInput(nil)
	if err == nil {
		t.Fatal("ToInput returned nil error")
	}
	if got, want := err.Error(), "album is required when album_spotify_id, album_url, or album_image_url is provided"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
