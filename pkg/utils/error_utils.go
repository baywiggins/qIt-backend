package utils

import (
	"encoding/json"
	"log"
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

// handleSpotifyError centralizes error handling for Spotify API calls.
func HandleSpotifyError(w http.ResponseWriter, err error) {
	log.Printf("Spotify API error: %s\n", err)
	if err.Error() == "invalid access token" {
		RespondWithStatusUnavailable(w)
	} else {
		RespondWithError(w, http.StatusInternalServerError, "Spotify API error: "+err.Error())
	}
}