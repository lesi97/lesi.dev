package utils

import "fmt"

func GetYungerWeaponData(name string) (*SpecificUserWeaponData, error) {
	switch name {
	case "cloudstrike":
		return &SpecificUserWeaponData{Weapon: "6917529926351690059", HashID: "", CrucibleTracker: "3244015567"}, nil
	case "beloved":
		return &SpecificUserWeaponData{Weapon: "6917529798698579887", HashID: "", CrucibleTracker: "3244015567"}, nil
	default:
		return nil, fmt.Errorf("weapon '%s' not found", name)
	}
}
