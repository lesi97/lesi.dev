package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/httpapi"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func BungieGET(ctx context.Context, logger *utils.Logger, httpClient *http.Client, clientID string, url string) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	startedAt := time.Now()
	defer func() {
		if ctx.Err() == nil {
			logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", url), startedAt, ctx)
		}
	}()

	headers := map[string]string{
		"X-API-Key": clientID,
	}
	body, statusCode, err := httpapi.DoRequest(ctx, httpClient, http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d: %s", statusCode, string(body))
	}

	isErr, err := isBungieError(body)
	if isErr {
		return nil, err
	}

	return body, nil
}
