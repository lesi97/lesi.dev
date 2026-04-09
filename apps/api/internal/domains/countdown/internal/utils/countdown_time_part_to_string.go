package utils

import "fmt"

func CountdownTimePartToString(value int, unit string) string {
	if value != 1 {
		return fmt.Sprintf("%d %ss", value, unit)
	}

	return fmt.Sprintf("%d %s", value, unit)
}
