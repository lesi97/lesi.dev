package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lesi97/lesi.dev/internal/domains/request_capture/internal/model"
	domain_utils "github.com/lesi97/lesi.dev/internal/domains/request_capture/internal/utils"
	shared_utils "github.com/lesi97/lesi.dev/internal/utils"
)

func (h *Handler) PostRequestCapture(w http.ResponseWriter, r *http.Request) {
	if err := domain_utils.ValidateRequestCaptureAccess(r); err != nil {
		h.logger.Printf("access denied: %v", err)
		shared_utils.Error(w, http.StatusForbidden, "access denied")
		return
	}

	maxBodyBytes, err := domain_utils.GetRequestCaptureMaxBodyBytes()
	if err != nil {
		h.logger.Printf("Request capture body limit invalid: %v", err)
		shared_utils.Error(w, http.StatusInternalServerError, "body limit misconfigured")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	body, err := domain_utils.ReadRequestCaptureBody(r.Body)
	if err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			shared_utils.Error(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}

		h.logger.Printf("Request capture body read failed: contentType=%q contentLength=%d err=%v", r.Header.Get("Content-Type"), r.ContentLength, err)
		shared_utils.Error(w, http.StatusBadRequest, "could not read request body")
		return
	}

	err = domain_utils.ValidateRequestCaptureBody(body)
	if err != nil {
		h.logger.Printf("Request capture body invalid JSON: contentType=%q contentLength=%d err=%v", r.Header.Get("Content-Type"), r.ContentLength, err)
		shared_utils.Error(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	headers, err := domain_utils.MarshalRequestCaptureHeaders(r.Header)
	if err != nil {
		h.logger.Printf("Request capture headers marshal failed: %v", err)
		shared_utils.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	requestCaptureID, err := h.store.InsertRequestCapture(r.Context(), model.RequestCaptureInput{
		Body:          body,
		ContentLength: r.ContentLength,
		ContentType:   r.Header.Get("Content-Type"),
		Headers:       headers,
		IP:            shared_utils.GetRequestIP(r),
		Path:          r.URL.Path,
		UserAgent:     r.UserAgent(),
	})
	if err != nil {
		h.logger.Printf("ERROR: Request capture POST: %v", err)
		shared_utils.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]any{
		"message": map[string]any{
			"id":                  requestCaptureID,
			"storedBodyBytes":     len(body),
			"maxAllowedBodyBytes": maxBodyBytes,
		},
		"error": nil,
	})
	if err != nil {
		h.logger.Printf("ERROR: Request capture response encode: %v", err)
	}
}
