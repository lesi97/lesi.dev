package model

type DateMemeClickAction string

const (
	DateMemeClickActionYes DateMemeClickAction = "yes"
	DateMemeClickActionNo  DateMemeClickAction = "no"
)

type DateMemeClickPayload struct {
	Route  string              `json:"route"`
	Action DateMemeClickAction `json:"action"`
}

type DateMemeClickInput struct {
	Route     string
	Action    DateMemeClickAction
	IP        string
	UserAgent string
}
