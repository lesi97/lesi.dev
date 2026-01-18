package utils

import "net/http"

type Methods interface {
	NormalisePayloadString(payload string) string
	ValidateRequest(r *http.Request) (int, string, error)
}

type Store struct{}

func NewStore() *Store {
	return &Store{}

}
