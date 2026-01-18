package utils

import "time"

type plugObjective struct {
	ObjectiveHash   int64 `json:"objectiveHash"`
	Progress        int   `json:"progress"`
	CompletionValue int   `json:"completionValue"`
	Complete        bool  `json:"complete"`
	Visible         bool  `json:"visible"`
}

type plug struct {
	PlugItemHash      int64           `json:"plugItemHash"`
	CanInsert         bool            `json:"canInsert"`
	Enabled           bool            `json:"enabled"`
	PlugObjectives    []plugObjective `json:"plugObjectives,omitempty"`
	InsertFailIndexes []int           `json:"insertFailIndexes,omitempty"`
	EnableFailIndexes []int           `json:"enableFailIndexes,omitempty"`
}

type plugs struct {
	Plugs map[string][]plug `json:"plugs"`
}

type plugSets struct {
	Data    plugs `json:"data"`
	Privacy int   `json:"privacy"`
}

type levelProgression struct {
	Level               int `json:"level"`
	ProgressToNextLevel int `json:"progressToNextLevel"`
	NextLevelAt         int `json:"nextLevelAt"`
}

type character struct {
	MembershipID         string            `json:"membershipId"`
	MembershipType       int               `json:"membershipType"`
	CharacterID          string            `json:"characterId"`
	ClassHash            int64             `json:"classHash"`
	RaceHash             int64             `json:"raceHash"`
	GenderHash           int64             `json:"genderHash"`
	ClassType            int               `json:"classType"`
	EmblemPath           string            `json:"emblemPath"`
	EmblemBackgroundPath string            `json:"emblemBackgroundPath"`
	EmblemHash           int64             `json:"emblemHash"`
	Light                int               `json:"light"`
	LevelProgression     levelProgression  `json:"levelProgression"`
	BaseCharacterLevel   int               `json:"baseCharacterLevel"`
	PercentToNextLevel   float64           `json:"percentToNextLevel"`
	DateLastPlayed       time.Time         `json:"dateLastPlayed"`
	MinutesPlayedTotal   string            `json:"minutesPlayedTotal"`
	Stats                map[string]int    `json:"stats"`
}

type characters struct {
	Data    map[string]character `json:"data"`
	Privacy int                  `json:"privacy"`
}

type equipmentItem struct {
	ItemHash                   int    `json:"itemHash"`
	ItemInstanceID             string `json:"itemInstanceId"`
	Quantity                   int    `json:"quantity"`
	BindStatus                 int    `json:"bindStatus"`
	Location                   int    `json:"location"`
	BucketHash                 int    `json:"bucketHash"`
	TransferStatus             int    `json:"transferStatus"`
	Lockable                   bool   `json:"lockable"`
	State                      int    `json:"state"`
	DismantlePermission        int    `json:"dismantlePermission"`
	IsWrapper                  bool   `json:"isWrapper"`
	TooltipNotificationIndexes []int  `json:"tooltipNotificationIndexes"`
	VersionNumber              int    `json:"versionNumber"`
}

type individualSocket struct {
	PlugHash  int64 `json:"plugHash"`
	IsEnabled bool  `json:"isEnabled"`
	IsVisible bool  `json:"isVisible"`
}

type itemSockets struct {
	Sockets []individualSocket `json:"sockets"`
}

type sockets struct {
	Data    map[string]itemSockets `json:"data"`
	Privacy int                    `json:"privacy"`
}

type plugObjectiveMap struct {
	ObjectivesPerPlug map[string][]plugObjective `json:"objectivesPerPlug"`
}

type plugObjectives struct {
	Data    map[string]plugObjectiveMap `json:"data"`
	Privacy int                         `json:"privacy"`
}

type items struct {
	Items []equipmentItem `json:"items"`
}

type characterEquipment struct {
	Data    map[string]items `json:"data"`
	Privacy int              `json:"privacy"`
}

type characterPlugSets struct {
	Data    map[string]plugs `json:"data"`
	Privacy int              `json:"privacy"`
}

type perk struct {
	PerkHash int64  `json:"perkHash"`
	IconPath string `json:"iconPath"`
	IsActive bool   `json:"isActive"`
	Visible  bool   `json:"visible"`
}

type perksForItem struct {
	Perks []perk `json:"perks"`
}

type itemPerks struct {
	Data    map[string]perksForItem `json:"data"`
	Privacy int                     `json:"privacy"`
}

type itemComponents struct {
	Sockets        sockets        `json:"sockets"`
	PlugObjectives plugObjectives `json:"plugObjectives"`
	Perks          itemPerks      `json:"perks"`
}

type bungieProfileResponse struct {
	ResponseMintedTimestamp            time.Time          `json:"responseMintedTimestamp"`
	SecondaryComponentsMintedTimestamp time.Time          `json:"secondaryComponentsMintedTimestamp"`
	ProfilePlugSets                    plugSets           `json:"profilePlugSets"`
	Characters                         characters         `json:"characters"`
	CharacterEquipment                 characterEquipment `json:"characterEquipment"`
	CharacterPlugSets                  characterPlugSets  `json:"characterPlugSets"`
	ItemComponents                     itemComponents     `json:"itemComponents"`
}

type BungieProfile struct {
	Response        bungieProfileResponse `json:"Response"`
	ErrorCode       int                   `json:"ErrorCode"`
	ThrottleSeconds int                   `json:"ThrottleSeconds"`
	ErrorStatus     string                `json:"ErrorStatus"`
	Message         string                `json:"Message"`
	MessageData     struct{}              `json:"MessageData"`
}
