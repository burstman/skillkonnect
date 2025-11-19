package app

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func JsonError(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := APIError{
		Error:   code,
		Message: message,
	}

	json.NewEncoder(w).Encode(res)
}
