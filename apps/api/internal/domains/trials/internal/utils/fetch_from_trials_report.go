package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/trials/internal/model"
	"github.com/lesi97/lesi.dev/internal/httpapi"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func FetchFromTrialsReport(ctx context.Context, logger *utils.Logger, url string) (*model.TrialsData, error) {
	now := time.Now()
	freshFor := 5 * time.Minute
	staleFor := 30 * time.Minute
	shouldLog := false
	startedAt := time.Now()
	defer func() {
		if shouldLog {
			logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", url), startedAt, nil)
		}
	}()

	var cachedData *model.TrialsData
	var cacheAge time.Duration

	trialsReportCacheMu.Lock()
	if trialsReportCache != nil {
		cachedData = trialsReportCache.Data
		cacheAge = now.Sub(trialsReportCache.CachedAt)
	}
	trialsReportCacheMu.Unlock()

	if cachedData != nil && cacheAge < freshFor {
		logger.PrintCache("CACHE HIT fetchFromTrialsReport")
		return cachedData, nil
	}

	if !IsTrialsReportAvailable(now) {
		if cachedData != nil && cacheAge < staleFor {
			return cachedData, nil
		}
		return nil, fmt.Errorf("trials report is not available")
	}

	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	shouldLog = true
	headers := map[string]string{
		"Accept": "application/json",
	}
	body, _, err := httpapi.DoRequest(reqCtx, http.DefaultClient, http.MethodGet, url, nil, headers)
	if err != nil {
		if cachedData != nil && cacheAge < staleFor {
			return cachedData, nil
		}
		return nil, err
	}

	result := &model.TrialsData{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		if cachedData != nil && cacheAge < staleFor {
			return cachedData, nil
		}
		return nil, err
	}

	trialsReportCacheMu.Lock()
	trialsReportCache = &TrialsReportCache{
		Data:     result,
		CachedAt: time.Now(),
	}
	trialsReportCacheMu.Unlock()

	return result, nil
}
