package model

import "strings"

func (d *PlexWebhookPayload) GetEpisodeNumber() int {
	episode := 1
	if strings.ToLower(d.Metadata.Type) == "episode" {
		episode = d.Metadata.Index
	}
	return episode
}
