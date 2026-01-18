package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	anilist_handler "github.com/lesi97/lesi.dev/internal/domains/anilist/handler"
	"github.com/lesi97/lesi.dev/internal/domains/anilist/model"
	"github.com/lesi97/lesi.dev/internal/domains/anilist/store"
	ani_utils "github.com/lesi97/lesi.dev/internal/domains/anilist/utils/anilist"
	http_utils "github.com/lesi97/lesi.dev/internal/domains/anilist/utils/http"
	"github.com/lesi97/lesi.dev/internal/domains/anilist/utils/plex"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func TestHandleUpdateAnilistBleachSeason2Episode10(t *testing.T) {
	plexSecret := envOrFail(t, "PLEX_WEBHOOK_SECRET")
	plexUsername := envOrFail(t, "PLEX_USERNAME")
	plexToken := envOrFail(t, "PLEX_X_TOKEN")
	anilistClientID := envOrFail(t, "ANILIST_CLIENT_ID")
	anilistClientSecret := envOrFail(t, "ANILIST_CLIENT_SECRET")
	anilistAccessToken := envOrFail(t, "ANILIST_ACCESS_TOKEN")
	anilistRefreshToken := envOrFail(t, "ANILIST_REFRESH_TOKEN")
	anilistBaseURL := envOrFail(t, "ANILIST_BASE_URL")
	anilistAuthURL := envOrFail(t, "ANILIST_AUTH_URL")

	anilistCalls := newAnilistRecorder()
	anilistServer := httptest.NewServer(http.HandlerFunc(anilistCalls.handle))
	t.Cleanup(anilistServer.Close)

	plexServer := httptest.NewServer(http.HandlerFunc(plexHandler(20, 12)))
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
			Index:                 10,
			ParentIndex:           2,
			LibrarySectionTitle:   "Anime",
			RatingKey:             "rk-episode",
			ParentRatingKey:       "rk-season2",
			GrandparentRatingKey:  "rk-show",
			GrandparentTitle:      "Bleach",
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
	if preview.TargetMediaID != 100 {
		t.Fatalf("expected target media ID 100, got %d", preview.TargetMediaID)
	}
	if preview.Progress != 30 {
		t.Fatalf("expected progress 30, got %d", preview.Progress)
	}
	if anilistCalls.updateProgressCalls != 0 {
		t.Fatalf("expected no update progress calls, got %d", anilistCalls.updateProgressCalls)
	}
}

type anilistRecorder struct {
	updateProgressCalls int
}

func newAnilistRecorder() *anilistRecorder {
	return &anilistRecorder{}
}

func (r *anilistRecorder) handle(w http.ResponseWriter, req *http.Request) {
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
					"mediaId":  100,
					"progress": 30,
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
					},
				},
			},
		})
	case strings.Contains(query, "Media(id"):
		writeJSON(w, map[string]interface{}{
			"data": map[string]interface{}{
				"Media": map[string]interface{}{
					"id":         100,
					"title":      map[string]string{"romaji": "Bleach", "english": "Bleach", "native": "BLEACH"},
					"format":     "TV",
					"seasonYear": 2004,
					"episodes":   366,
					"isAdult":    false,
					"relations": map[string]interface{}{
						"edges": []interface{}{},
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

func plexHandler(season1Count int, season2Count int) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if strings.HasSuffix(path, "/children") {
			if strings.Contains(path, "/rk-show/children") {
				writeJSON(w, map[string]interface{}{
					"MediaContainer": map[string]interface{}{
						"Metadata": []map[string]interface{}{
							{"ratingKey": "s1", "index": 1},
							{"ratingKey": "s2", "index": 2},
						},
					},
				})
				return
			}
			if strings.Contains(path, "/s1/children") {
				writeJSON(w, seasonChildren(season1Count))
				return
			}
			if strings.Contains(path, "/s2/children") {
				writeJSON(w, seasonChildren(season2Count))
				return
			}
		}

		if strings.Contains(path, "/library/metadata/") {
			writeJSON(w, map[string]interface{}{
				"MediaContainer": map[string]interface{}{
					"Metadata": []map[string]interface{}{
						{"Tag": []interface{}{}},
					},
				},
			})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func seasonChildren(count int) map[string]interface{} {
	metadata := make([]map[string]interface{}, 0, count)
	for i := 0; i < count; i++ {
		metadata = append(metadata, map[string]interface{}{
			"ratingKey": "ep-" + strconv.Itoa(i+1),
			"index":     i + 1,
		})
	}

	return map[string]interface{}{
		"MediaContainer": map[string]interface{}{
			"Metadata": metadata,
		},
	}
}

func buildMultipartPayload(t *testing.T, payload model.PlexWebhookPayload) (io.Reader, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	err = writer.WriteField("payload", string(raw))
	if err != nil {
		t.Fatalf("failed to write payload field: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	return body, writer.FormDataContentType()
}

func writeJSON(w http.ResponseWriter, payload map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}
