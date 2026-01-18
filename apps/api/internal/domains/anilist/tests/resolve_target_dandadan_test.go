package tests

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

func TestHandleUpdateAnilistDandadanSeason2Episode6(t *testing.T) {
	plexSecret := envOrFail(t, "PLEX_WEBHOOK_SECRET")
	plexUsername := envOrFail(t, "PLEX_USERNAME")
	plexToken := envOrFail(t, "PLEX_X_TOKEN")
	anilistClientID := envOrFail(t, "ANILIST_CLIENT_ID")
	anilistClientSecret := envOrFail(t, "ANILIST_CLIENT_SECRET")
	anilistAccessToken := envOrFail(t, "ANILIST_ACCESS_TOKEN")
	anilistRefreshToken := envOrFail(t, "ANILIST_REFRESH_TOKEN")
	anilistBaseURL := envOrFail(t, "ANILIST_BASE_URL")
	anilistAuthURL := envOrFail(t, "ANILIST_AUTH_URL")

	anilistCalls := newDandadanAnilistRecorder()
	anilistServer := httptest.NewServer(http.HandlerFunc(anilistCalls.handle))
	t.Cleanup(anilistServer.Close)

	plexServer := httptest.NewServer(http.HandlerFunc(plexHandler(12, 12)))
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
			Index:                 6,
			ParentIndex:           2,
			LibrarySectionTitle:   "Anime",
			RatingKey:             "rk-episode",
			ParentRatingKey:       "rk-season2",
			GrandparentRatingKey:  "rk-show",
			GrandparentTitle:      "Dandadan",
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
	if preview.TargetMediaID != 201 {
		t.Fatalf("expected target media ID 201, got %d", preview.TargetMediaID)
	}
	if preview.Progress != 6 {
		t.Fatalf("expected progress 6, got %d", preview.Progress)
	}
	if anilistCalls.updateProgressCalls != 0 {
		t.Fatalf("expected no update progress calls, got %d", anilistCalls.updateProgressCalls)
	}
}

type dandadanAnilistRecorder struct {
	updateProgressCalls int
}

func newDandadanAnilistRecorder() *dandadanAnilistRecorder {
	return &dandadanAnilistRecorder{}
}

func (r *dandadanAnilistRecorder) handle(w http.ResponseWriter, req *http.Request) {
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
					"mediaId":  201,
					"progress": 6,
				},
			},
		})
	case strings.Contains(query, "Page(perPage"):
		writeJSON(w, map[string]interface{}{
			"data": map[string]interface{}{
				"Page": map[string]interface{}{
					"media": []map[string]interface{}{
						{
							"id":         200,
							"title":      map[string]string{"romaji": "Dandadan", "english": "Dandadan", "native": "DANDADAN"},
							"format":     "TV",
							"seasonYear": 2024,
							"status":     "RELEASING",
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
	case strings.Contains(query, "Media(id"):
		idValue := payload.Variables["id"]
		id := intFromInterface(idValue)
		edges := []map[string]interface{}{}
		if id == 200 {
			edges = append(edges, map[string]interface{}{
				"relationType": "SEQUEL",
				"node": map[string]interface{}{
					"id":         201,
					"title":      map[string]string{"romaji": "Dandadan S2", "english": "Dandadan Season 2", "native": "DANDADAN"},
					"seasonYear": 2025,
					"episodes":   nil,
					"format":     "TV",
					"isAdult":    false,
				},
			})
		}

		writeJSON(w, map[string]interface{}{
			"data": map[string]interface{}{
				"Media": map[string]interface{}{
					"id":         id,
					"title":      map[string]string{"romaji": "Dandadan", "english": "Dandadan", "native": "DANDADAN"},
					"format":     "TV",
					"seasonYear": 2024,
					"episodes":   nil,
					"isAdult":    false,
					"relations": map[string]interface{}{
						"edges": edges,
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

func intFromInterface(value interface{}) int {
	switch v := value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		parsed, _ := strconv.Atoi(v)
		return parsed
	default:
		return 0
	}
}
