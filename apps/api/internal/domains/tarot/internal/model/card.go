package model

type Card struct {
	Name       string `json:"name"`
	NameShort  string `json:"name_short"`
	Value      string `json:"value"`
	ValueInt   int    `json:"value_int"`
	Suit       string `json:"suit"`
	CardType   string `json:"card_type"`
	MeaningUp  string `json:"meaning_up"`
	MeaningRev string `json:"meaning_rev"`
	Desc       string `json:"desc"`
}

type CardsResponse struct {
	Cards *[]Card `json:"cards"`
}
