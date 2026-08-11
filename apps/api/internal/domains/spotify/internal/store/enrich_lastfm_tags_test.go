package store

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/spotify/internal/model"
)

func TestLastFMTopTagsFromResponseParsesAndDedupesTags(t *testing.T) {
	tags, err := lastfmTopTagsFromResponse([]byte(`{
		"toptags": {
			"tag": [
				{"name":"rock","count":"100","url":"https://www.last.fm/tag/rock"},
				{"name":"Rock","count":50,"url":"https://www.last.fm/tag/Rock"},
				{"name":" alternative ","count":20,"url":"https://www.last.fm/tag/alternative"},
				{"name":" ","count":1}
			]
		}
	}`))
	if err != nil {
		t.Fatalf("lastfmTopTagsFromResponse returned error: %v", err)
	}

	if got, want := len(tags), 2; got != want {
		t.Fatalf("tag count = %d, want %d", got, want)
	}
	if got, want := tags[0].Name, "rock"; got != want {
		t.Fatalf("first tag = %q, want %q", got, want)
	}
	if got, want := int(tags[0].Count), 100; got != want {
		t.Fatalf("first tag count = %d, want %d", got, want)
	}
	if got, want := tags[1].Name, "alternative"; got != want {
		t.Fatalf("second tag = %q, want %q", got, want)
	}
}

func TestLastFMTopTagsFromResponseTreatsMissingResourceAsNoTags(t *testing.T) {
	tags, err := lastfmTopTagsFromResponse([]byte(`{"error":6,"message":"Track not found"}`))
	if err != nil {
		t.Fatalf("lastfmTopTagsFromResponse returned error: %v", err)
	}
	if tags != nil {
		t.Fatalf("tags = %v, want nil", tags)
	}
}

func TestLastFMTopTagsFromResponseReturnsRateLimitError(t *testing.T) {
	_, err := lastfmTopTagsFromResponse([]byte(`{"error":29,"message":"Rate limit exceeded"}`))
	var rateLimitErr *lastfmRateLimitError
	if !errors.As(err, &rateLimitErr) {
		t.Fatalf("error = %v, want lastfmRateLimitError", err)
	}
}

func TestLastFMInfoTagsFromResponseTreatsEmptyStringTagsAsNoTags(t *testing.T) {
	tags, err := lastfmInfoTagsFromResponse([]byte(`{"album":{"tags":""}}`), model.SpotifyEnrichmentTypeAlbum)
	if err != nil {
		t.Fatalf("lastfmInfoTagsFromResponse returned error: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("tags = %v, want none", tags)
	}
}

func TestLastFMInfoTagsFromResponseParsesTopTags(t *testing.T) {
	tags, err := lastfmInfoTagsFromResponse([]byte(`{
		"track": {
			"toptags": {
				"tag": {"name":"hip-hop","count":"42","url":"https://www.last.fm/tag/hip-hop"}
			}
		}
	}`), model.SpotifyEnrichmentTypeTrack)
	if err != nil {
		t.Fatalf("lastfmInfoTagsFromResponse returned error: %v", err)
	}
	if got, want := len(tags), 1; got != want {
		t.Fatalf("tag count = %d, want %d", got, want)
	}
	if got, want := tags[0].Name, "hip-hop"; got != want {
		t.Fatalf("tag name = %q, want %q", got, want)
	}
	if got, want := int(tags[0].Count), 42; got != want {
		t.Fatalf("tag count value = %d, want %d", got, want)
	}
}

func TestLastFMTagTargetsForScrobbleUsesCreatedEntities(t *testing.T) {
	albumID := int64(20)
	targets := lastfmTagTargetsForScrobble(
		model.ScrobbleInput{
			Artist: model.ScrobbleEntityInput{Name: "Hozier"},
			Album:  &model.ScrobbleEntityInput{Name: "Wasteland, Baby!"},
			Track:  model.ScrobbleEntityInput{Name: "Talk"},
		},
		&model.ScrobbleResult{
			ArtistID:      10,
			AlbumID:       &albumID,
			TrackID:       30,
			ArtistCreated: true,
			AlbumCreated:  true,
			TrackCreated:  true,
		},
	)

	if got, want := len(targets), 3; got != want {
		t.Fatalf("target count = %d, want %d", got, want)
	}
	if got, want := targets[0].EntityType, model.SpotifyEnrichmentTypeArtist; got != want {
		t.Fatalf("first entity type = %q, want %q", got, want)
	}
	if got, want := targets[1].EntityType, model.SpotifyEnrichmentTypeAlbum; got != want {
		t.Fatalf("second entity type = %q, want %q", got, want)
	}
	if got, want := targets[2].EntityType, model.SpotifyEnrichmentTypeTrack; got != want {
		t.Fatalf("third entity type = %q, want %q", got, want)
	}
	if got, want := targets[2].Album, "Wasteland, Baby!"; got != want {
		t.Fatalf("track fallback album = %q, want %q", got, want)
	}
}

func TestLastFMTagTargetsForScrobbleSkipsExistingEntities(t *testing.T) {
	targets := lastfmTagTargetsForScrobble(
		model.ScrobbleInput{
			Artist: model.ScrobbleEntityInput{Name: "Hozier"},
			Track:  model.ScrobbleEntityInput{Name: "Talk"},
		},
		&model.ScrobbleResult{
			ArtistID: 10,
			TrackID:  30,
		},
	)

	if len(targets) != 0 {
		t.Fatalf("target count = %d, want 0", len(targets))
	}
}

func TestLastFMTagSettings(t *testing.T) {
	t.Setenv("LASTFM_TAGS_ON_POLL", "")
	if !lastfmTagPollingEnabled() {
		t.Fatal("Last.fm tag polling should default to enabled")
	}

	t.Setenv("LASTFM_TAGS_ON_POLL", "false")
	if lastfmTagPollingEnabled() {
		t.Fatal("Last.fm tag polling should be disabled")
	}

	t.Setenv("LASTFM_TAG_LIMIT", "100")
	if got, want := lastfmTagLimit(), 50; got != want {
		t.Fatalf("lastfmTagLimit() = %d, want %d", got, want)
	}

	t.Setenv("LASTFM_RATE_LIMIT_COOLDOWN_SECONDS", "30")
	if got, want := lastfmRateLimitCooldown(), 30*time.Second; got != want {
		t.Fatalf("lastfmRateLimitCooldown() = %s, want %s", got, want)
	}
}

func TestFetchLastFMTopTagsFallsBackFromTrackToAlbumTags(t *testing.T) {
	methods := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.URL.Query().Get("method")
		methods = append(methods, method)

		if got, want := r.URL.Query().Get("artist"), "SpaceMan Zack"; got != want {
			t.Fatalf("artist = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		switch method {
		case "track.getTopTags", "track.getInfo":
			_, _ = w.Write([]byte(`{"toptags":{"tag":[]}}`))
		case "album.getTopTags":
			if got, want := r.URL.Query().Get("album"), "Drugs in the Rain"; got != want {
				t.Fatalf("album = %q, want %q", got, want)
			}
			_, _ = w.Write([]byte(`{"toptags":{"tag":[{"name":"emo rap","count":"10"}]}}`))
		default:
			t.Fatalf("unexpected method %q", method)
		}
	}))
	defer server.Close()

	store := &Store{HTTPClient: server.Client()}
	lookup, err := store.fetchLastFMTopTags(
		context.Background(),
		lastfmAPIAuth{APIKey: "key", APIURL: server.URL},
		lastfmTagTarget{
			EntityType: model.SpotifyEnrichmentTypeTrack,
			Artist:     "SpaceMan Zack",
			Album:      "Drugs in the Rain",
			Track:      "Drugs in the Rain",
		},
	)
	if err != nil {
		t.Fatalf("fetchLastFMTopTags returned error: %v", err)
	}

	if got, want := len(lookup.Tags), 1; got != want {
		t.Fatalf("tag count = %d, want %d", got, want)
	}
	if got, want := lookup.Tags[0].Name, "emo rap"; got != want {
		t.Fatalf("tag name = %q, want %q", got, want)
	}
	if got, want := lookup.Notes, "used Last.fm album tags fallback"; got != want {
		t.Fatalf("notes = %q, want %q", got, want)
	}

	wantMethods := []string{"track.getTopTags", "track.getInfo", "album.getTopTags"}
	if len(methods) != len(wantMethods) {
		t.Fatalf("methods = %v, want %v", methods, wantMethods)
	}
	for i := range wantMethods {
		if methods[i] != wantMethods[i] {
			t.Fatalf("methods = %v, want %v", methods, wantMethods)
		}
	}
}
