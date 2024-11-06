package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/baywiggins/qIt-backend/internal/api/middlewares"
	"github.com/baywiggins/qIt-backend/pkg/utils"
)

var authURL = make(map[string]string)

func HandleSpotifyAuthRoutes() {
	// Handle our GET auth endpoint which allows frontend to get auth URL
	http.Handle("GET /spotify/auth", middlewares.LoggingMiddleware(http.HandlerFunc(handleSpotifyAuth)))
	// Handle our callback GET endpoint which the user is redirected to once authenticated
	http.Handle("GET /spotify/auth/callback", middlewares.LoggingMiddleware(http.HandlerFunc(handleSpotifyAuthCallBackGet)))

}

// Function to handle auth
func handleSpotifyAuth(w http.ResponseWriter, r *http.Request) {
	var err error;
	// Create return mapping of our authentication URL
	authURL["auth_url"], err = middlewares.GetSpotifyAuthURL()

	if err != nil {
		log.Printf("ERROR: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSpotifyAuth: '%s' \n", err.Error()))
		return
	}
	
	//Set header as JSON
	w.Header().Set("Content-Type", "application/json")
	// Create our JSON encoder, and ensure it does not escape our characters
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.SetEscapeHTML(false)
	// Encode the auth url to a json and return it to the caller
	jsonEncoder.Encode(authURL)
}

func handleSpotifyAuthCallBackGet(w http.ResponseWriter, r *http.Request) {
	// Get the authorization code or error from the request URL
	authCode := r.URL.Query().Get("code")
	errorCode := r.URL.Query().Get("error")
	if authCode == "" && errorCode != ""{
		log.Printf("ERROR in handleSpotiftAuthCallBackGet: '%s' \n", errorCode)
		utils.RespondWithError(w, http.StatusInternalServerError, errorCode)
		return
	} else if authCode == "" {
		log.Printf("ERROR in handleSpotiftAuthCallBackGet: Unknown Error \n")
		utils.RespondWithError(w, http.StatusInternalServerError, "Unknown Error")
		return
	}

	// Grab the access token using the auth code
	err := middlewares.GetAccessTokenFromSpotify(authCode)
	if err != nil {
		log.Printf("ERROR: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}


