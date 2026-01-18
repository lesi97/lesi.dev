package model

import "errors"

func (d *PlexWebhookPayload) GetShowName() (string, error) {
	showName := d.Metadata.GrandparentTitle
	if showName == "" {
		showName = d.Metadata.ParentTitle
	}
	if showName == "" {
		showName = d.Metadata.Title
	}
	if showName == "" {
		return "", errors.New("unknown showname")
	}
	return showName, nil
}
