package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lesi97/lesi.dev/internal/domains/trials/internal/model"
	"github.com/lesi97/lesi.dev/internal/utils"
)

func FetchFromTrialsReport(logger *utils.Logger, url string) (*model.TrialsData, error) {
	defer logger.LogExecutionTime(fmt.Sprintf("EXTERNAL API CALL: %v", url), time.Now(), nil)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &model.TrialsData{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(result)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}
	return result, nil
}
