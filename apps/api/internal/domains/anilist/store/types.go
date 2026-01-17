package anilist_store

type PlexGuid struct {
	ID string `json:"id"`
}

type PlexWebhookPayload struct {
	Event   string `json:"event"`
	User    bool   `json:"user"`
	Owner   bool   `json:"owner"`
	Account struct {
		ID    int    `json:"id"`
		Thumb string `json:"thumb"`
		Title string `json:"title"`
	} `json:"Account"`
	Server struct {
		Title string `json:"title"`
		UUID  string `json:"uuid"`
	} `json:"Server"`
	Player struct {
		Local         bool   `json:"local"`
		PublicAddress string `json:"publicAddress"`
		Title         string `json:"title"`
		UUID          string `json:"uuid"`
	} `json:"Player"`
	Metadata PlexMetadata `json:"Metadata"`
}

type PlexMetadata struct {
	LibrarySectionType    string  `json:"librarySectionType"`
	RatingKey             string  `json:"ratingKey"`
	Key                   string  `json:"key"`
	ParentRatingKey       string  `json:"parentRatingKey"`
	GrandparentRatingKey  string  `json:"grandparentRatingKey"`
	PlexGuid              string  `json:"-"`
	ParentGUID            string  `json:"parentGuid"`
	GrandparentGUID       string  `json:"grandparentGuid"`
	GrandparentSlug       string  `json:"grandparentSlug"`
	Type                  string  `json:"type"`
	Title                 string  `json:"title"`
	GrandparentKey        string  `json:"grandparentKey"`
	ParentKey             string  `json:"parentKey"`
	LibrarySectionTitle   string  `json:"librarySectionTitle"`
	LibrarySectionID      int     `json:"librarySectionID"`
	LibrarySectionKey     string  `json:"librarySectionKey"`
	GrandparentTitle      string  `json:"grandparentTitle"`
	ParentTitle           string  `json:"parentTitle"`
	OriginalTitle         string  `json:"originalTitle"`
	ContentRating         string  `json:"contentRating"`
	Summary               string  `json:"summary"`
	Index                 int     `json:"index"`
	ParentIndex           int     `json:"parentIndex"`
	AudienceRating        float64 `json:"audienceRating"`
	ViewCount             int     `json:"viewCount"`
	LastViewedAt          int     `json:"lastViewedAt"`
	Year                  int     `json:"year"`
	Thumb                 string  `json:"thumb"`
	Art                   string  `json:"art"`
	ParentThumb           string  `json:"parentThumb"`
	GrandparentThumb      string  `json:"grandparentThumb"`
	GrandparentArt        string  `json:"grandparentArt"`
	GrandparentTheme      string  `json:"grandparentTheme"`
	Duration              int     `json:"duration"`
	OriginallyAvailableAt string  `json:"originallyAvailableAt"`
	AddedAt               int     `json:"addedAt"`
	UpdatedAt             int     `json:"updatedAt"`
	AudienceRatingImage   string  `json:"audienceRatingImage"`
	ChapterSource         string  `json:"chapterSource"`

	Image []struct {
		Alt  string `json:"alt"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"Image"`

	UltraBlurColors struct {
		TopLeft     string `json:"topLeft"`
		TopRight    string `json:"topRight"`
		BottomRight string `json:"bottomRight"`
		BottomLeft  string `json:"bottomLeft"`
	} `json:"UltraBlurColors"`

	GUID []struct {
		ID string `json:"id"`
	} `json:"-"`

	Rating []struct {
		Image string  `json:"image"`
		Value float64 `json:"value"`
		Type  string  `json:"type"`
	} `json:"Rating"`

	Director []struct {
		ID     int    `json:"id"`
		Filter string `json:"filter"`
		Tag    string `json:"tag"`
		TagKey string `json:"tagKey"`
	} `json:"Director"`

	Writer []struct {
		ID     int    `json:"id"`
		Filter string `json:"filter"`
		Tag    string `json:"tag"`
		TagKey string `json:"tagKey"`
		Thumb  string `json:"thumb"`
	} `json:"Writer"`

	Role []struct {
		ID     int    `json:"id"`
		Filter string `json:"filter"`
		Tag    string `json:"tag"`
		TagKey string `json:"tagKey"`
		Role   string `json:"role"`
		Thumb  string `json:"thumb"`
	} `json:"Role"`
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