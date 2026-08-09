package store

import (
	"testing"

	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
)

func TestBestSpotifyTrackEnrichmentMatchAcceptsExactTrackArtist(t *testing.T) {
	albumName := "Wasteland, Baby!"
	artistName := "Hozier"
	candidate := spotifyEnrichmentCandidate{
		EntityType: model.SpotifyEnrichmentTypeTrack,
		ID:         1,
		Name:       "Talk",
		ArtistName: &artistName,
		AlbumName:  &albumName,
	}

	match := bestSpotifyTrackEnrichmentMatch(candidate, []spotifyTrack{
		{
			ID:   "track-123",
			Name: "Talk",
			ExternalURLs: spotifyExternalURLs{
				Spotify: "https://open.spotify.com/track/track-123",
			},
			Artists: []spotifyArtist{{Name: "Hozier"}},
			Album: spotifyAlbum{
				Name: "Wasteland, Baby!",
				Images: []spotifyImage{
					{URL: "https://i.scdn.co/image/album"},
				},
			},
		},
	})

	if got, want := match.Status, spotifyEnrichmentStatusMatched; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
	if match.Confidence == nil || *match.Confidence < spotifyTrackMatchThreshold {
		t.Fatalf("confidence = %v, want >= %v", match.Confidence, spotifyTrackMatchThreshold)
	}
	if got, want := stringValue(match.SpotifyURL), "https://open.spotify.com/track/track-123"; got != want {
		t.Fatalf("spotify url = %q, want %q", got, want)
	}
}

func TestBestSpotifyTrackEnrichmentMatchRejectsWeakArtistMatch(t *testing.T) {
	artistName := "Hozier"
	candidate := spotifyEnrichmentCandidate{
		EntityType: model.SpotifyEnrichmentTypeTrack,
		ID:         1,
		Name:       "Talk",
		ArtistName: &artistName,
	}

	match := bestSpotifyTrackEnrichmentMatch(candidate, []spotifyTrack{
		{
			ID:   "track-123",
			Name: "Talk",
			ExternalURLs: spotifyExternalURLs{
				Spotify: "https://open.spotify.com/track/track-123",
			},
			Artists: []spotifyArtist{{Name: "Different Artist"}},
		},
	})

	if got, want := match.Status, spotifyEnrichmentStatusLowConfidence; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}

func TestBestSpotifyAlbumEnrichmentMatchAcceptsExactAlbumArtist(t *testing.T) {
	artistName := "Hozier"
	candidate := spotifyEnrichmentCandidate{
		EntityType: model.SpotifyEnrichmentTypeAlbum,
		ID:         1,
		Name:       "Wasteland, Baby!",
		ArtistName: &artistName,
	}

	match := bestSpotifyAlbumEnrichmentMatch(candidate, []spotifyAlbum{
		{
			ID:   "album-123",
			Name: "Wasteland, Baby!",
			ExternalURLs: spotifyExternalURLs{
				Spotify: "https://open.spotify.com/album/album-123",
			},
			Artists: []spotifyArtist{{Name: "Hozier"}},
		},
	})

	if got, want := match.Status, spotifyEnrichmentStatusMatched; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}

func TestBestSpotifyArtistEnrichmentMatchOnlyAcceptsExactArtist(t *testing.T) {
	candidate := spotifyEnrichmentCandidate{
		EntityType: model.SpotifyEnrichmentTypeArtist,
		ID:         1,
		Name:       "Hozier",
	}

	match := bestSpotifyArtistEnrichmentMatch(candidate, []spotifyArtist{
		{
			ID:   "artist-123",
			Name: "Hozier",
			ExternalURLs: spotifyExternalURLs{
				Spotify: "https://open.spotify.com/artist/artist-123",
			},
		},
	})

	if got, want := match.Status, spotifyEnrichmentStatusMatched; got != want {
		t.Fatalf("status = %q, want %q", got, want)
	}
}
