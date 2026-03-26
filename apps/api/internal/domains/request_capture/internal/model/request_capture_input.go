package model

type RequestCaptureInput struct {
	Body          string
	ContentLength int64
	ContentType   string
	Headers       string
	IP            string
	Path          string
	UserAgent     string
}
