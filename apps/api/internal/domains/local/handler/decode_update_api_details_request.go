package handler

import (
	"encoding/json"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/domains/local/model"
)

func DecodeUpdateApiDetailsRequest(r *http.Request) (model.UpdateApiDetailsRequest, error) {
	var req model.UpdateApiDetailsRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		return model.UpdateApiDetailsRequest{}, err
	}

	return req, nil
}
