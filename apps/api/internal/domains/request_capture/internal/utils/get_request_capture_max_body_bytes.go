package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetRequestCaptureMaxBodyBytes() (int64, error) {
	rawValue := strings.TrimSpace(os.Getenv("REQUEST_CAPTURE_MAX_BODY_BYTES"))
	if rawValue == "" {
		return 25 << 20, nil
	}

	value, err := strconv.ParseInt(rawValue, 10, 64)
	if err != nil {
		return 0, err
	}

	if value <= 0 {
		return 0, strconv.ErrSyntax
	}

	return value, nil
}
