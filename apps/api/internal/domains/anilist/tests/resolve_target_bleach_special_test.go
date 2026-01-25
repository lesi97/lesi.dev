package tests

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	anilist_handler "github.com/lesi97/lesi.dev/internal/domains/anilist/handler"
	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/store"
	ani_utils "github.com/lesi97/lesi.dev/internal/domains/anilist/internal/utils/anilist"
	http_utils "github.com/lesi97/lesi.dev/internal/domains/anilist/internal/utils/http"
	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/utils/plex"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleUpdateAnilistBleachSpecialEpisode2(t *testing.T) {
	plexSecret := envOrFail(t, "PLEX_WEBHOOK_SECRET")
	plexUsername := envOrFail(t, "PLEX_USERNAME")
	plexToken := envOrFail(t, "PLEX_X_TOKEN")
	anilistClientID := envOrFail(t, "ANILIST_CLIENT_ID")
	anilistClientSecret := envOrFail(t, "ANILIST_CLIENT_SECRET")
	anilistAccessToken := envOrFail(t, "ANILIST_ACCESS_TOKEN")
	anilistRefreshToken := envOrFail(t, "ANILIST_REFRESH_TOKEN")
	anilistBaseURL := envOrFail(t, "ANILIST_BASE_URL")
	anilistAuthURL := envOrFail(t, "ANILIST_AUTH_URL")

	anilistCalls := newBleachSpecialAnilistRecorder()
	anilistServer := httptest.NewServer(http.HandlerFunc(anilistCalls.handle))
	t.Cleanup(anilistServer.Close)

	plexServer := httptest.NewServer(http.HandlerFunc(plexHandler(1, 1)))
	t.Cleanup(plexServer.Close)

	logger := &utils.Logger{Logger: log.New(io.Discard, "", 0)}
	aniEnv := &model.AnilistEnv{
		DiscordUsername: "Anilist",
		BaseUrl:         anilistBaseURL,
		GraphqlUrl:      anilistServer.URL,
		AuthUrl:         anilistAuthURL,
		ClientId:        anilistClientID,
		ClientSecret:    anilistClientSecret,
		AccessToken:     anilistAccessToken,
		RefreshToken:    anilistRefreshToken,
		RefreshExpires:  time.Now().Add(24 * time.Hour).UnixMilli(),
	}
	plexEnv := &model.PlexEnv{
		BaseUrl:  plexServer.URL,
		Username: plexUsername,
		XToken:   plexToken,
	}

	aniUtils := ani_utils.NewStore(logger, aniEnv, nil)
	plexUtils := plex.NewStore(logger, plexEnv)
	anilistStore := &store.Store{
		Logger:    logger,
		AniEnv:    aniEnv,
		PlexEnv:   plexEnv,
		AniUtils:  aniUtils,
		PlexUtils: plexUtils,
	}

	handler := anilist_handler.NewHandlerWithDeps(logger, anilistStore, http_utils.NewStore())

	payload := model.PlexWebhookPayload{
		Event: "media.scrobble",
		Account: struct {
			ID    int    `json:"id"`
			Thumb string `json:"thumb"`
			Title string `json:"title"`
		}{
			ID:    1,
			Title: plexUsername,
		},
		Metadata: model.PlexMetadata{
			Type:                  "episode",
			Index:                 2,
			ParentIndex:           0,
			LibrarySectionTitle:   "Anime",
			RatingKey:             "rk-episode",
			ParentRatingKey:       "rk-season0",
			GrandparentRatingKey:  "rk-show",
			GrandparentTitle:      "Bleach",
			ParentTitle:           "Specials",
			Title:                 "The Sealed Sword Frenzy",
			LibrarySectionType:    "show",
			GrandparentTheme:      "",
			GrandparentSlug:       "",
			GrandparentGUID:       "",
			ParentGUID:            "",
			GrandparentKey:        "",
			ParentKey:             "",
			GrandparentArt:        "",
			GrandparentThumb:      "",
			ParentThumb:           "",
			Art:                   "",
			Thumb:                 "",
			OriginallyAvailableAt: "",
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}
	t.Logf("data in payload=%s", string(payloadJSON))

	body, contentType := buildMultipartPayload(t, payload)
	req := httptest.NewRequest(http.MethodPost, "/v1/anilist?secret="+url.QueryEscape(plexSecret), body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "PlexMediaServer")

	preview := &store.UpdatePreview{}
	req = req.WithContext(store.WithUpdatePreview(req.Context(), preview))

	recorder := httptest.NewRecorder()
	handler.HandleUpdateAnilist(recorder, req)

	t.Logf("data out status=%d targetMediaID=%d progress=%d", recorder.Code, preview.TargetMediaID, preview.Progress)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204 status, got %d", recorder.Code)
	}
	if preview.TargetMediaID != 101 {
		t.Fatalf("expected target media ID 101, got %d", preview.TargetMediaID)
	}
	if preview.Progress != 1 {
		t.Fatalf("expected progress 1, got %d", preview.Progress)
	}
	if anilistCalls.updateProgressCalls != 0 {
		t.Fatalf("expected no update progress calls, got %d", anilistCalls.updateProgressCalls)
	}
}

type bleachSpecialAnilistRecorder struct {
	updateProgressCalls int
}

func newBleachSpecialAnilistRecorder() *bleachSpecialAnilistRecorder {
	return &bleachSpecialAnilistRecorder{}
}

func (r *bleachSpecialAnilistRecorder) handle(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var payload struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}
	_ = json.NewDecoder(req.Body).Decode(&payload)
	query := payload.Query

	switch {
	case strings.Contains(query, "SaveMediaListEntry"):
		r.updateProgressCalls++
		writeJSON(w, map[string]interface{}{
			"data": map[string]interface{}{
				"SaveMediaListEntry": map[string]interface{}{
					"id":       1,
					"mediaId":  101,
					"progress": 1,
				},
			},
		})
	case strings.Contains(query, "Page(perPage"):
		writeJSON(w, map[string]interface{}{
			"data": map[string]interface{}{
				"Page": map[string]interface{}{
					"media": []map[string]interface{}{
						{
							"id":         100,
							"title":      map[string]string{"romaji": "Bleach", "english": "Bleach", "native": "BLEACH"},
							"format":     "TV",
							"seasonYear": 2004,
							"episodes":   366,
							"status":     "FINISHED",
							"isAdult":    false,
							"relations": map[string]interface{}{
								"edges": []interface{}{},
								"nodes": []interface{}{},
							},
						},
						{
							"id":         101,
							"title":      map[string]string{"romaji": "Bleach: The Sealed Sword Frenzy", "english": "Bleach: The Sealed Sword Frenzy", "native": "BLEACH"},
							"format":     "OVA",
							"seasonYear": 2006,
							"episodes":   1,
							"status":     "FINISHED",
							"isAdult":    false,
							"relations": map[string]interface{}{
								"edges": []interface{}{},
								"nodes": []interface{}{},
							},
						},
					},
				},
			},
		})
	default:
		writeJSON(w, map[string]interface{}{
			"data": map[string]interface{}{},
		})
	}
}
