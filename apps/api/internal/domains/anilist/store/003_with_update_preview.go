package store

import "context"

type UpdatePreview struct {
	TargetMediaID int
	Progress      int
}

type contextKey string

const updatePreviewKey contextKey = "anilist_update_preview"

func WithUpdatePreview(ctx context.Context, preview *UpdatePreview) context.Context {
	return context.WithValue(ctx, updatePreviewKey, preview)
}
