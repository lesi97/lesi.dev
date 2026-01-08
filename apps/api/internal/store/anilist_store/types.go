package anilist_store

type PlexGuid struct {
	ID string `json:"id"`
}

type PlexWebhookPayload struct {
	Event   *string `json:"event,omitempty"`
	Account *struct {
		Title *string `json:"title,omitempty"`
	} `json:"Account,omitempty"`
	Metadata *struct {
		LibrarySectionTitle  *string    `json:"librarySectionTitle,omitempty"`
		Type                 *string    `json:"type,omitempty"`
		Title                *string    `json:"title,omitempty"`
		Year                 *int       `json:"year,omitempty"`
		Index                *int       `json:"index,omitempty"`
		ParentIndex          *int       `json:"parentIndex,omitempty"`
		GrandparentTitle     *string    `json:"grandparentTitle,omitempty"`
		ParentTitle          *string    `json:"parentTitle,omitempty"`
		GrandparentRatingKey *string    `json:"grandparentRatingKey,omitempty"`
		Guid                 []PlexGuid `json:"Guid,omitempty"`
	} `json:"Metadata,omitempty"`
}

type AniListTitle struct {
	Romaji  *string `json:"romaji,omitempty"`
	English *string `json:"english,omitempty"`
	Native  *string `json:"native,omitempty"`
}

type AniListMedia struct {
	ID         int           `json:"id"`
	Title      *AniListTitle `json:"title,omitempty"`
	Format     *string       `json:"format,omitempty"`
	SeasonYear *int          `json:"seasonYear,omitempty"`
	Episodes   *int          `json:"episodes,omitempty"`
	Status     *string       `json:"status,omitempty"`
	Relations  *struct {
		Edges []struct {
			RelationType *string `json:"relationType,omitempty"`
		} `json:"edges,omitempty"`
		Nodes []AniListMedia `json:"nodes,omitempty"`
	} `json:"relations,omitempty"`
}