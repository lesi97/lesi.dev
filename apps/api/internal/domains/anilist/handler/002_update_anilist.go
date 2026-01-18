package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/lesi97/lesi.dev/internal/domains/anilist/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

type contextKey string
const contextKeyPlexPayload contextKey = "anilist"

func (h *Handler) HandleUpdateAnilist(w http.ResponseWriter, r *http.Request) {

	status, plexUsername, err := h.httpUtils.ValidateRequest(r)
	if err != nil {
		if os.Getenv("GO_ENV") != "development" {
			utils.Error(w, status, err.Error())
			return
		}
	}

	const discord_username = "Anilist Updater"
	ctx := context.WithValue(r.Context(), contextKeyPlexPayload, discord_username)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		message := fmt.Sprintf("%v", err)
		h.logger.SendDiscordNotification(utils.SendDiscordNotificationArgs{
			Content:  message,
			Username: discord_username,
			Title:    "Failed to parse Anilist form data",
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

	var plexData model.PlexWebhookPayload

	err = json.Unmarshal([]byte(payload), &plexData)
	if err != nil {
		var decoded string
		err2 := json.Unmarshal([]byte(payload), &decoded)
		if err2 == nil {
			err = json.Unmarshal([]byte(decoded), &plexData)
		}

		if err != nil {
			normalised := h.httpUtils.NormalisePayloadString(payload)
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
			Content:  message,
			Username: discord_username,
			Title:    "Failed to update Anilist",
		})
		h.logger.Error(message)
		utils.Error(w, http.StatusInternalServerError, "failed to update anilist")
		return
	}

	utils.Success(w, http.StatusNoContent, nil)

}
