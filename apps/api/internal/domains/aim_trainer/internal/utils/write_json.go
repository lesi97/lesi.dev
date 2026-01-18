package utils

import (
	"encoding/json"
	"net/http"
)

// Legacy support for aim-trainer game which expects a specific resposne
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
