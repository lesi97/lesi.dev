package utils

import "fmt"

func GetTerrorWeaponData(name string) (*SpecificUserWeaponData, error) {
	switch name {
	case "ace":
		return &SpecificUserWeaponData{Weapon: "6917529207684081719", HashID: "38912240", CrucibleTracker: "38912240"}, nil
	case "felwinter":
		return &SpecificUserWeaponData{Weapon: "6917529190261952418", HashID: "1179141605", CrucibleTracker: "3244015567"}, nil
	case "matador":
		return &SpecificUserWeaponData{Weapon: "6917529875871677239", HashID: "2563012876", CrucibleTracker: "3244015567"}, nil
	case "immortal":
		return &SpecificUserWeaponData{Weapon: "6917529880229623656", HashID: "38912240", CrucibleTracker: "38912240"}, nil
	case "thorn":
		return &SpecificUserWeaponData{Weapon: "6917529935035554307", HashID: "3973202132", CrucibleTracker: "38912240"}, nil
	default:
		return nil, fmt.Errorf("weapon '%s' not found", name)
	}
}
