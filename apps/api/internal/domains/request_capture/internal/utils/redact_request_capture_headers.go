package utils

import "net/http"

func RedactRequestCaptureHeaders(headers http.Header) http.Header {
	clonedHeaders := headers.Clone()
	clonedHeaders.Del("X-Api-Key")
	return clonedHeaders
}
