package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/baywiggins/qIt-backend/internal/api/middlewares"
	"github.com/baywiggins/qIt-backend/internal/models"
	"github.com/baywiggins/qIt-backend/internal/services"
	"github.com/baywiggins/qIt-backend/pkg/utils"
)

var authURL = make(map[string]string)

func HandleSpotifyAuthRoutes(db *sql.DB) {
	handler := &Handler{DB: db}
	// Handle our GET auth endpoint which allows frontend to get auth URL
	http.Handle("GET /spotify/auth", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(handleSpotifyAuth))))
	// Handle our callback GET endpoint which the user is redirected to once authenticated
	http.Handle("GET /spotify/auth/callback", middlewares.LoggingMiddleware(http.HandlerFunc(handler.handleSpotifyAuthCallBackGet)))
	// Handle testing our spotify authentication
	http.Handle("GET /spotify/auth/test-spotify-auth", middlewares.LoggingMiddleware(http.HandlerFunc(handler.handleTestSpotifyAuth)))
}

// Function to handle auth
func handleSpotifyAuth(w http.ResponseWriter, r *http.Request) {
	var err error;
	// Get uuid from query string
	state := r.URL.Query().Get("state")
	if state == "" {
		log.Printf("ERROR in handleSpotifyAuth: state not provided with request \n")
		utils.RespondWithError(w, http.StatusBadRequest, "Error: state not provided with request")
		return
	}
	// Create return mapping of our authentication URL
	authURL["auth_url"], err = middlewares.GetSpotifyAuthURL(state)

	if err != nil {
		log.Printf("ERROR in handleSpotifyAuth: %s \n", err)
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

func (h *Handler) handleSpotifyAuthCallBackGet(w http.ResponseWriter, r *http.Request) {
	// Handle a page refresh
	
	// Get the authorization code or error from the request URL
	authCode := r.URL.Query().Get("code")
	errorCode := r.URL.Query().Get("error")
	state := r.URL.Query().Get("state")
	if authCode == "" && errorCode != ""{
		log.Printf("ERROR in handleSpotifyAuthCallBackGet: '%s' \n", errorCode)
		utils.RespondWithError(w, http.StatusInternalServerError, errorCode)
		return
	} else if authCode == "" {
		log.Printf("ERROR in handleSpotifyAuthCallBackGet: Unknown Error \n")
		utils.RespondWithError(w, http.StatusInternalServerError, "Unknown Error")
		return
	} else if state == "" {
		log.Printf("ERROR in handleSpotifyAuthCallBackGet: Unknown Error \n")
		utils.RespondWithError(w, http.StatusInternalServerError, "Unknown Error")
		return
	}

	// Grab the access token using the auth code
	err := middlewares.GetAccessTokenFromSpotify(authCode, state, h.DB)
	if err != nil {
		log.Printf("ERROR in handleSpotifyAuthCallBackGet: %s \n", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == "user must authenticate with spotify first" {
			statusCode = http.StatusUnauthorized
		}
		utils.RespondWithError(w, statusCode, err.Error())
		return
	}

	// Return success screen to user
	http.ServeFile(w, r, "./static/index.html")
}

func (h *Handler) handleTestSpotifyAuth(w http.ResponseWriter, r *http.Request) {
	// TODO: IF UNAUTHORIZED AND ENTRY EXISTS IN DB, TRY TO GET REFRESH TOKEN, IF THAT FAILS, SEND SOMETHING SPECIAL TO INDICATE THE USER MUST RE-AUTH WITH SPOTIFY


	// Get the uuid header
	uuid := r.Header.Get("uuid")

	// Get the spotify auth token from the db
	token, err := models.GetAuthTokenByID(h.DB, uuid)

	fmt.Println(uuid)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Error in handleTestSpotifyAuth" 
		if (err.Error() == "sql: no rows in result set") {
			status = http.StatusUnauthorized
			message = "Unauthorized user attempted to access Spotify"
		}
		log.Printf(message + "\n")
		utils.RespondWithError(w, status, message)
		return
	}

	fmt.Println(token)

	sURL, err := url.Parse("https://api.spotify.com/v1/me")
	if err != nil {
		log.Printf("ERROR in handleTestSpotifyAuth: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleTestSpotifyAuth: '%s' \n", err.Error()))
		return
	}

	headers := map[string]string {
		"Authorization": "Bearer "+token,
	}

	// Make spotify API request
	_, err = services.SendSpotifyPlayerRequest(*sURL, http.MethodGet, nil, headers)

	if err != nil {
		log.Printf("ERROR or unauth in handleTestSpotifyAuth: %s \n", err)
		utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error: '%s' \n", err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
