package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/lesi97/lesi.dev/internal/store/anilist_store"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type AnilistHandler struct {
	logger         *utils.Logger
	store 			anilist_store.AnilistStoreInterface
}

const anilistContextKey = "anilist"

func NewAnilistHandler(logger *utils.Logger, store anilist_store.AnilistStoreInterface)  *AnilistHandler {
	return &AnilistHandler{
		logger: logger,
		store: store,
	}
}

type contextKey string

const contextKeyPlexPayload contextKey = "discord_username"

func (h *AnilistHandler) HandleUpdateAnilist(w http.ResponseWriter, r *http.Request) {

	secret := r.URL.Query().Get("secret")
	if secret == "" || secret != os.Getenv("PLEX_WEBHOOK_SECRET") {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	ua := r.Header.Get("User-Agent")
	if !strings.Contains(ua, "PlexMediaServer") {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	plexUsername := os.Getenv("PLEX_USERNAME")
	if plexUsername == "" {
		utils.Error(w, http.StatusInternalServerError, "No Plex username found in env")
		return
	}

	const discord_username = "Anilist Updater"
	ctx := context.WithValue(r.Context(), contextKeyPlexPayload, discord_username)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		message := fmt.Sprintf("%v", err)
		h.logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: discord_username,
			Title: "Failed to parse Anilist form data",
		})
		h.logger.Error(message)
		utils.Error(w, http.StatusBadRequest, err)
		return
	}

	payload := r.FormValue("payload")
	if payload == "" {
		utils.Error(w, http.StatusBadRequest, "missing payload")
		return
	}

	var plexData anilist_store.PlexWebhookPayload

	err = json.Unmarshal([]byte(payload), &plexData)
	if err != nil {
		var decoded string
		err2 := json.Unmarshal([]byte(payload), &decoded)
		if err2 == nil {
			err = json.Unmarshal([]byte(decoded), &plexData)
		}

		if err != nil {
			normalised := normalisePayloadString(payload)
			err3 := json.Unmarshal([]byte(normalised), &plexData)
			if err3 != nil {
				message := fmt.Sprintf("payload parse failed \nerr=%v \nerr2=%v \nerr3=%v \nprefix=%q", err, err2, err3, payload[:min(20, len(payload))])
				h.logger.Error(message)
				utils.Error(w, http.StatusBadRequest, "invalid payload json")
				return
			}
		}
	}


	// Only scrobble for intended user for Anime libraries as this is for Anilist, an anime platform
	mediaTypeRegex := regexp.MustCompile(`(?i)anime`)
	if strings.ToLower(plexData.Event) != "media.scrobble" ||
		!mediaTypeRegex.MatchString(strings.ToLower(plexData.Metadata.LibrarySectionTitle)) ||
		!strings.EqualFold(plexData.Account.Title, plexUsername) {
			utils.Success(w, http.StatusNoContent, nil)
			return
	}

	err = h.store.UpdateAnilist(ctx, plexData)
	if err != nil {
		message := fmt.Sprintf("%v", err)
		h.logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content: message,
			Username: discord_username,
			Title: "Failed to update Anilist",
		})
		h.logger.Error(message)
		utils.Error(w, http.StatusInternalServerError, "failed to update anilist")
		return
	}

	utils.Success(w, http.StatusNoContent, nil)

}

func normalisePayloadString(payload string) string {
	s := strings.TrimSpace(payload)

	if strings.HasPrefix(s, `"{`) || strings.HasPrefix(s, `"[\{`) {
		return s
	}

	if strings.HasPrefix(s, `\{`) {
		s = strings.TrimPrefix(s, `\`)
		s = strings.ReplaceAll(s, `\"`, `"`)
		return s
	}

	if strings.HasPrefix(s, `\[`){
		s = strings.TrimPrefix(s, `\`)
		s = strings.ReplaceAll(s, `\"`, `"`)
		return s
	}

	return s
}
