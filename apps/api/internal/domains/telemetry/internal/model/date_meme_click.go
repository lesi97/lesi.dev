package model

type DateMemeClickAction string

const (
	DateMemeClickActionYes DateMemeClickAction = "yes"
	DateMemeClickActionNo  DateMemeClickAction = "no"
)

type DateMemeClickPayload struct {
	Route        string              `json:"route"`
	Action       DateMemeClickAction `json:"action"`
	SecretEnding bool                `json:"secretEnding,omitempty"`
}

type DateMemeClickInput struct {
	Route        string
	Action       DateMemeClickAction
	SecretEnding bool
	IP           string
	UserAgent    string
}
