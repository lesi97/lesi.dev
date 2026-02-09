package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type SendDiscordNotificationArgs struct {
	Content  string
	Username string
	Title 	 string
}

type DiscordNotificationData struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

func (l *Logger) SendDiscordNotification(data SendDiscordNotificationArgs) {
	
	if os.Getenv("GO_ENV") == "development" {
		return
	}

	discordURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if discordURL == "" {
		l.Error("DISCORD_WEBHOOK_URL is not defined in environment variables")
		return
	}

	var content string = "";
	if data.Content != "" {
		content = fmt.Sprintf("```%s```",data.Content)
	}

	
	var contentMessage string
	userID := os.Getenv("DISCORD_USER_ID")
	if userID == "" {
		l.Error("DISCORD_USER_ID is not defined in environment variables")
		contentMessage = fmt.Sprintf("%v %s", data.Title, content)
	} else {
		contentMessage = fmt.Sprintf("<@%s> %v %s`", userID, data.Title, content)
	}


	payload := DiscordNotificationData{
		Content:  contentMessage,
		Username: data.Username,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		l.Error(fmt.Sprintf("Error marshalling Discord payload: %s", err.Error()))
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodPost, discordURL, bytes.NewReader(bodyBytes))
	if err != nil {
		l.Error(fmt.Sprintf("Error creating Discord request: %s", err.Error()))
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		l.Error(fmt.Sprintf("Error sending Discord notification: %s", err.Error()))
		return
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		l.Error(fmt.Sprintf("Discord webhook failed: %s %s", resp.Status, string(respBody)))
		return
	}
}
