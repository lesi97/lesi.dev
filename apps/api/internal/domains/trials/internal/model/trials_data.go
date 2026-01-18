package model

type TrialsData struct {
	WeekNumber     int            `json:"weekNumber"`
	RealWeekNumber int            `json:"realWeekNumber"`
	StartDate      string         `json:"startDate"`
	EndDate        string         `json:"endDate"`
	Maps           []TrialsMap    `json:"maps"`
	Rewards        TrialsRewards  `json:"rewards"`
	Platforms      TrialsPlatforms `json:"platforms"`
}
