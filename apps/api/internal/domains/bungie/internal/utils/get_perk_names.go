package utils

import (
	"fmt"
	"strings"
)

func getPerkNames(perks []databasePerk) string {
	var perkNames []string
	for _, perk := range perks {
		if strings.Contains(perk.ItemType, "Enhanced") {
			perkName := fmt.Sprintf("Enhanced %s", perk.Name)
			perkNames = append(perkNames, perkName)
		} else {
			perkNames = append(perkNames, perk.Name)
		}
	}
	return strings.Join(perkNames, ", ")
}
