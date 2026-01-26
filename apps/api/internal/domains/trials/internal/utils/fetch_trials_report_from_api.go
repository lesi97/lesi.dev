package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/domains/trials/internal/model"
	"github.com/lesi97/lesi.dev/internal/httpapi"
)

func FetchTrialsReportFromAPI(ctx context.Context, url string) (*model.TrialsData, error) {
	headers := map[string]string{
		"Accept": "application/json",
	}
	body, _, err := httpapi.DoRequest(ctx, http.DefaultClient, http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}

	result := &model.TrialsData{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
