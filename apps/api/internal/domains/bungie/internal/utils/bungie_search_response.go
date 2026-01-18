package utils

type bungieSearchResponse struct {
	IconPath                    string `json:"iconPath"`
	CrossSaveOverride           int    `json:"crossSaveOverride"`
	ApplicableMembershipTypes   []int  `json:"applicableMembershipTypes"`
	IsPublic                    bool   `json:"isPublic"`
	MembershipType              int    `json:"membershipType"`
	MembershipID                string `json:"membershipId"`
	DisplayName                 string `json:"displayName"`
	BungieGlobalDisplayName     string `json:"bungieGlobalDisplayName"`
	BungieGlobalDisplayNameCode int    `json:"bungieGlobalDisplayNameCode"`
}
