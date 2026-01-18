package utils

import "strconv"

func GetPlatformEnum(platform string, membershipType int) string {
	str := switchPlatforms(platform)
	if str == "-1" {
		return strconv.Itoa(membershipType)
	}
	return str
}
