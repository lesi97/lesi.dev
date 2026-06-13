package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type manifestString string

func (m *manifestString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*m = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*m = manifestString(s)
		return nil
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var number json.Number
	if err := decoder.Decode(&number); err == nil {
		*m = manifestString(number.String())
		return nil
	}

	return fmt.Errorf("expected manifest string or number, got %s", string(data))
}

func (m manifestString) String() string {
	return string(m)
}

type manifestInt int

func (m *manifestInt) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*m = 0
		return nil
	}

	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*m = manifestInt(i)
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if s == "" {
			*m = 0
			return nil
		}

		parsed, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("parse manifest int %q: %w", s, err)
		}

		*m = manifestInt(parsed)
		return nil
	}

	return fmt.Errorf("expected manifest int or string, got %s", string(data))
}

func (m manifestInt) Int() int {
	return int(m)
}

type dbWeapons struct {
	ID                  string
	DisplayName         string
	ItemTypeDisplayName string
	FlavorText          string
	BucketTypeHash      string
	TierTypeHash        string
	TierTypeName        string
	TierType            int
	TalentGridHash      string
}

type dbPerk struct {
	Name        string
	Description string
	ItemType    string
	HashID      string
}

type bungieWeaponRaw struct {
	Hash              manifestString `json:"hash"`
	ItemType          int            `json:"itemType"`
	DisplayProperties struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"displayProperties"`
	ItemTypeDisplayName string `json:"itemTypeDisplayName"`
	FlavorText          string `json:"flavorText"`
	Inventory           struct {
		BucketTypeHash manifestString `json:"bucketTypeHash"`
		TierTypeHash   manifestString `json:"tierTypeHash"`
		TierTypeName   string         `json:"tierTypeName"`
		TierType       manifestInt    `json:"tierType"`
	} `json:"inventory"`
	Plug struct {
		PlugCategoryIdentifier string `json:"plugCategoryIdentifier"`
	} `json:"plug"`
	TalentGrid struct {
		TalentGridHash manifestString `json:"talentGridHash"`
	} `json:"talentGrid"`
}

func (w bungieWeaponRaw) hashID(fallback string) string {
	if hash := w.Hash.String(); hash != "" {
		return hash
	}
	return fallback
}
