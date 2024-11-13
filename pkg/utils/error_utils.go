package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	response := ErrorResponse{
		Code: code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func RespondWithStatusUnavailable(w http.ResponseWriter) {
	// Set header as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	// Create JSON encoder and encode our currentlyPlaying variable
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(map[string]string{"message": "send another request"})
}