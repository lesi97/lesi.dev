package httpapi

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	requestmetrics "github.com/lesi97/lesi.dev/internal/request_metrics"
)

func DoRequest(
	ctx context.Context,
	client *http.Client,
	method string,
	url string,
	body io.Reader,
	headers map[string]string,
) ([]byte, int, error) {
	start := time.Now()
	safeURL := RedactSensitiveQueryValues(url)

	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		safeError := strings.ReplaceAll(err.Error(), url, safeURL)
		return nil, 0, fmt.Errorf("failed to create %s request for %s: %s", method, safeURL, safeError)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		requestmetrics.AddFetchCallsDuration(ctx, time.Since(start), err)
		safeError := strings.ReplaceAll(err.Error(), url, safeURL)
		return nil, 0, fmt.Errorf("%s request failed for %s: %s", method, safeURL, safeError)
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			requestmetrics.AddFetchCallsDuration(ctx, time.Since(start), err)
			return nil, resp.StatusCode, err
		}
		defer gz.Close()
		reader = gz
	}

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		requestmetrics.AddFetchCallsDuration(ctx, time.Since(start), err)
		return nil, resp.StatusCode, err
	}

	requestmetrics.AddFetchCallsDuration(ctx, time.Since(start), nil)
	return bodyBytes, resp.StatusCode, nil
}
