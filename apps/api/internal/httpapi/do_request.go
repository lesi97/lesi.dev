package httpapi

import (
	"context"
	"io"
	"net/http"
)

func DoRequest(
	ctx context.Context,
	client *http.Client,
	method string,
	url string,
	body io.Reader,
	headers map[string]string,
) ([]byte, int, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, 0, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return bodyBytes, resp.StatusCode, nil
}
