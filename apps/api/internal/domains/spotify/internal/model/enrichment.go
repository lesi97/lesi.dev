package model

const (
	SpotifyEnrichmentTypeTrack  = "track"
	SpotifyEnrichmentTypeAlbum  = "album"
	SpotifyEnrichmentTypeArtist = "artist"
)

type SpotifyEnrichmentInput struct {
	EntityType string
}

type SpotifyEnrichmentResult struct {
	EntityType                 string   `json:"entity_type"`
	EntityID                   *int64   `json:"entity_id,omitempty"`
	EntityName                 string   `json:"entity_name,omitempty"`
	Status                     string   `json:"status"`
	Matched                    bool     `json:"matched"`
	Updated                    bool     `json:"updated"`
	Confidence                 *float64 `json:"confidence,omitempty"`
	SpotifyID                  *string  `json:"spotify_id,omitempty"`
	SpotifyURL                 *string  `json:"spotify_url,omitempty"`
	ImageURL                   *string  `json:"image_url,omitempty"`
	Notes                      string   `json:"notes,omitempty"`
	RateLimited                bool     `json:"rate_limited,omitempty"`
	RateLimitReason            string   `json:"rate_limit_reason,omitempty"`
	RateLimitRetryAfterSeconds *int     `json:"rate_limit_retry_after_seconds,omitempty"`
}
