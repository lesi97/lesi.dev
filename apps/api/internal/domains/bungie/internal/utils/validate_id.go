package utils

import (
	"regexp"
	"strings"
)

func ValidateID(id string) bool {
	idRegex := regexp.MustCompile(`^[\w !@#$%^&*()_+={}\[\]:;"'<>,.?/\\|-]+#[0-9]{4}$`)
	invalidRegex := regexp.MustCompile(`\b(drop|alter|delete|insert|update|create|select|truncate|exec|union)\b`)
	if invalidRegex.MatchString(strings.ToLower(id)) {
		return false
	}
	if !idRegex.MatchString(id) {
		return false
	}
	return true
}
