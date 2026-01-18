package model

type PlexLabelsTag struct {
	Tag     string `json:"tag"`
	TagType string `json:"tagType"`
}

type PlexLabelsMetadata struct {
	Tag       []PlexLabelsTag `json:"Tag"`
	RatingKey string          `json:"ratingKey"`
	Index     int             `json:"index"`
}

type PlexLabelsMediaContainer struct {
	Metadata []PlexLabelsMetadata `json:"Metadata"`
}

type PlexLabels struct {
	MediaContainer PlexLabelsMediaContainer `json:"MediaContainer"`
}
