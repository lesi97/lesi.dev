package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ReadIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, errors.New("invalid ID parameter")
	}

	id, err := strconv.ParseInt(idParam, 10, 64) // Base 10, 64 bit int
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("invalid ID parameter")
	}

	return id, nil
}