package model

type AnilistGraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type AnilistGraphQLError struct {
	Message string `json:"message"`
}

type AnilistGraphQLResponse[T any] struct {
	Data   T                     `json:"data"`
	Errors []AnilistGraphQLError `json:"errors,omitempty"`
}

type AnilistTitle struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

type AnilistRelationEdge struct {
	RelationType string `json:"relationType"`
}

type AnilistRelationNode struct {
	ID         int          `json:"id"`
	Title      AnilistTitle `json:"title"`
	SeasonYear *int         `json:"seasonYear"`
	Episodes   *int         `json:"episodes"`
	Format     string       `json:"format"`
	IsAdult    bool         `json:"isAdult,omitempty"`
}

type AnilistRelations struct {
	Edges []AnilistRelationEdge `json:"edges"`
	Nodes []AnilistRelationNode `json:"nodes"`
}

type AnilistMedia struct {
	ID         int              `json:"id"`
	Title      AnilistTitle     `json:"title"`
	Format     string           `json:"format"`
	SeasonYear *int             `json:"seasonYear"`
	Episodes   *int             `json:"episodes"`
	Status     string           `json:"status"`
	IsAdult    bool             `json:"isAdult"`
	Relations  AnilistRelations `json:"relations"`
}

type AnilistMediaPage struct {
	Media []AnilistMedia `json:"media"`
}

type AnilistPageData struct {
	Page AnilistMediaPage `json:"Page"`
}

type AnilistSearchResponse struct {
	Data   AnilistPageData       `json:"data"`
	Errors []AnilistGraphQLError `json:"errors,omitempty"`
}

type AnilistRelationEdgeNode struct {
	RelationType string              `json:"relationType"`
	Node         AnilistRelationNode `json:"node"`
}

type AnilistMediaRelations struct {
	ID         int          `json:"id"`
	Title      AnilistTitle `json:"title"`
	Format     string       `json:"format"`
	SeasonYear *int         `json:"seasonYear"`
	Episodes   *int         `json:"episodes"`
	IsAdult    bool         `json:"isAdult"`
	Relations  struct {
		Edges []AnilistRelationEdgeNode `json:"edges"`
	} `json:"relations"`
}

type AnilistMediaRelationsData struct {
	Media AnilistMediaRelations `json:"Media"`
}

type AnilistMediaRelationsResponse struct {
	Data   AnilistMediaRelationsData `json:"data"`
	Errors []AnilistGraphQLError     `json:"errors,omitempty"`
}
