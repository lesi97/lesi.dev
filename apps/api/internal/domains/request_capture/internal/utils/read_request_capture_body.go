package utils

import "io"

func ReadRequestCaptureBody(body io.Reader) (string, error) {
	rawBody, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(rawBody), nil
}
