package store

import (
	"testing"
	"time"
)

func TestSpotifyPlayHistoryToScrobbleInputMapsSpotifyData(t *testing.T) {
	artistImageURL := "https://i.scdn.co/image/artist"
	item := spotifyPlayHistory{
		PlayedAt: time.Date(2026, 8, 6, 10, 37, 6, 0, time.UTC),
		Track: spotifyTrack{
			ID:   "track-123",
			Name: "Talk",
			ExternalURLs: spotifyExternalURLs{
				Spotify: "https://open.spotify.com/track/track-123",
			},
			Artists: []spotifyArtist{
				{
					ID:   "artist-123",
					Name: "Hozier",
					ExternalURLs: spotifyExternalURLs{
						Spotify: "https://open.spotify.com/artist/artist-123",
					},
				},
			},
			Album: spotifyAlbum{
				ID:   "album-123",
				Name: "Wasteland, Baby!",
				ExternalURLs: spotifyExternalURLs{
					Spotify: "https://open.spotify.com/album/album-123",
				},
				Images: []spotifyImage{
					{URL: "https://i.scdn.co/image/album-large"},
					{URL: "https://i.scdn.co/image/album-small"},
				},
			},
		},
	}

	input, ok, err := spotifyPlayHistoryToScrobbleInput(
		item,
		[]byte(`{"track":{"id":"track-123"},"played_at":"2026-08-06T10:37:06Z"}`),
		map[string]*string{
			"artist-123": &artistImageURL,
		},
	)
	if err != nil {
		t.Fatalf("spotifyPlayHistoryToScrobbleInput returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected Spotify play history to map to a scrobble")
	}

	if got, want := input.Source, spotifyPollSource; got != want {
		t.Fatalf("source = %q, want %q", got, want)
	}
	if got, want := input.Artist.Name, "Hozier"; got != want {
		t.Fatalf("artist = %q, want %q", got, want)
	}
	if got, want := stringValue(input.Artist.SpotifyID), "artist-123"; got != want {
		t.Fatalf("artist spotify id = %q, want %q", got, want)
	}
	if got, want := stringValue(input.Artist.ImageURL), artistImageURL; got != want {
		t.Fatalf("artist image url = %q, want %q", got, want)
	}
	if input.Album == nil {
		t.Fatal("album is nil")
	}
	if got, want := stringValue(input.Album.ImageURL), "https://i.scdn.co/image/album-large"; got != want {
		t.Fatalf("album image url = %q, want %q", got, want)
	}
	if got, want := stringValue(input.Track.ImageURL), "https://i.scdn.co/image/album-large"; got != want {
		t.Fatalf("track image url = %q, want %q", got, want)
	}
	if got, want := input.ScrobbledAt.Format(time.RFC3339), "2026-08-06T10:37:06Z"; got != want {
		t.Fatalf("scrobbled at = %q, want %q", got, want)
	}
}

func TestSpotifyAccessTokenExpiredSupportsUnixSecondsAndMilliseconds(t *testing.T) {
	futureSeconds := time.Now().Add(5 * time.Minute).Unix()
	pastSeconds := time.Now().Add(-5 * time.Minute).Unix()
	futureMillis := time.Now().Add(5 * time.Minute).UnixMilli()
	pastMillis := time.Now().Add(-5 * time.Minute).UnixMilli()

	if spotifyAccessTokenExpired(&futureSeconds) {
		t.Fatal("future unix seconds expiry should not be expired")
	}
	if !spotifyAccessTokenExpired(&pastSeconds) {
		t.Fatal("past unix seconds expiry should be expired")
	}
	if spotifyAccessTokenExpired(&futureMillis) {
		t.Fatal("future unix milliseconds expiry should not be expired")
	}
	if !spotifyAccessTokenExpired(&pastMillis) {
		t.Fatal("past unix milliseconds expiry should be expired")
	}
}
