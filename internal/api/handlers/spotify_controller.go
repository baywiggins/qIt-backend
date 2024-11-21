package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/baywiggins/qIt-backend/internal/api/middlewares"
	"github.com/baywiggins/qIt-backend/internal/config"
	"github.com/baywiggins/qIt-backend/internal/models"
	"github.com/baywiggins/qIt-backend/internal/services"
	"github.com/baywiggins/qIt-backend/pkg/utils"
)


func HandleSpotifyControllerRoutes(db *sql.DB) {
	// Get currently playing track
	http.Handle("GET /spotify/currently-playing", middlewares.LoggingMiddleware(http.HandlerFunc(handleCurrentlyPlaying)))
	// Get current queue
	http.Handle("GET /spotify/queue", middlewares.LoggingMiddleware(http.HandlerFunc(handleGetQueue)))
	// Get search by track
	http.Handle("GET /spotify/search/track", middlewares.LoggingMiddleware(http.HandlerFunc(handleSearchByTrack)))
	// Get search by URL (next/previous page of results)
	http.Handle("GET /spotify/search/url", middlewares.LoggingMiddleware(http.HandlerFunc(handleSearchByURL)))
	// Play
	http.Handle("GET /spotify/play", middlewares.LoggingMiddleware(http.HandlerFunc(handlePlay)))
	// Pause
	http.Handle("GET /spotify/pause", middlewares.LoggingMiddleware(http.HandlerFunc(handlePause)))
	// Add to queue
	http.Handle("GET /spotify/add-to-queue", middlewares.LoggingMiddleware(http.HandlerFunc(handleAddToQueue)))
}

func handleCurrentlyPlaying(w http.ResponseWriter, r *http.Request) {
	var err error;
	// Parse our request URL
	sURL, err := url.Parse(config.SpotifyPlayerURL + "/currently-playing")
	if err != nil {
		log.Printf("ERROR in handleCurrentlyPlaying: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error: '%s' \n", err.Error()))
		return
	}
	// Grab our access token from middleware cache
	accessToken, err := middlewares.GetAccessToken()
	if err != nil {
		log.Printf("ERROR in handleCurrentlyPlaying: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCurrentlyPlaying: '%s' \n", err.Error()))
		return
	}
	// Add our headers in a map
	headers := map[string]string {
		"Authorization": "Bearer "+accessToken,
	}
	// Call our spotify API function to get the response body
	body, err := services.SendSpotifyPlayerRequest(*sURL, http.MethodGet, nil, headers)
	if err != nil {
		if err.Error() == "invalid access token" {
			utils.RespondWithStatusUnavailable(w)
			return
		}
		log.Printf("ERROR in handleCurrentlyPlaying: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCurrentlyPlaying: '%s' \n", err.Error()))
		return
	}

	currentlyPlaying, err := services.UnmarshalJSON[models.CurrentlyPlaying](body)
	if err != nil {
		log.Printf("ERROR in handleCurrentlyPlaying: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCurrentlyPlaying: '%s' \n", err.Error()))
		return
	}
	// Set header as JSON
	w.Header().Set("Content-Type", "application/json")
	// Create JSON encoder and encode our currentlyPlaying variable
	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(currentlyPlaying)
	if err != nil {
		log.Printf("ERROR in handleCurrentlyPlaying: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCurrentlyPlaying: '%s' \n", err.Error()))
		return
	}
}

func handleGetQueue(w http.ResponseWriter, r *http.Request) {
	var err error;
	// Parse our request URL
	sURL, err := url.Parse(config.SpotifyPlayerURL + "/queue")
	if err != nil {
		log.Printf("ERROR in handleGetQueue: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleGetQueue: '%s' \n", err.Error()))
		return
	}
	// Grab our access token from middleware cache
	accessToken, err := middlewares.GetAccessToken()
	if err != nil {
		log.Printf("ERROR in handleGetQueue: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleGetQueue: '%s' \n", err.Error()))
		return
	}

	// Add our headers in a map
	headers := map[string]string {
		"Authorization": "Bearer "+accessToken,
	}
	// Call our spotify API function to get the response body
	body, err := services.SendSpotifyPlayerRequest(*sURL, http.MethodGet, nil, headers)
	if err != nil {
		if err.Error() == "invalid access token" {
			utils.RespondWithStatusUnavailable(w)
			return
		}
		log.Printf("ERROR in handleGetQueue: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByTrack: '%s' \n", err.Error()))
		return
	}
	// Unmarshal our data to return to caller
	currentQueue, err := services.UnmarshalJSON[models.CurrentQueue](body)
	if err != nil {
		log.Printf("ERROR in handleGetQueue: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleGetQueue: '%s' \n", err.Error()))
		return
	}
	// Set header as JSON
	w.Header().Set("Content-Type", "application/json")
	// Create JSON encoder and encode our currentlyPlaying variable
	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(currentQueue)
	if err != nil {
		log.Printf("ERROR in handleGetQueue: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleGetQueue: '%s' \n", err.Error()))
		return
	}
}

func handleSearchByTrack(w http.ResponseWriter, r *http.Request) {
	var err error;
	// Get the track from the GET request query param
	track := r.URL.Query().Get("track")
	if track == "" {
		log.Println("ERROR: 'track' was null", )
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByTrack: '%s' \n", fmt.Errorf("'track' was null")))
		return
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		log.Println("ERROR: 'limit' was null", )
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByTrack: '%s' \n", fmt.Errorf("'limit' was null")))
		return
	}
	// Parse our request URL
	sURL, err := url.Parse(config.SpotifySearchURL)
	if err != nil {
		log.Printf("ERROR in handleSearchByTrack: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByTrack: '%s' \n", err.Error()))
		return
	}
	// Grab our access token from middleware cache
	accessToken, err := middlewares.GetAccessToken()
	if err != nil {
		log.Printf("ERROR in handleSearchByTrack: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByTrack: '%s' \n", err.Error()))
		return
	}
	// Add our query params in a map
	queryParams := map[string]string {
		"q": "track:" + track,
		"type": "track",
		"limit": limit,
	}
	// Add our headers in a map
	headers := map[string]string {
		"Authorization": "Bearer "+accessToken,
	}
	// Call our spotify API function to get the response body
	body, err := services.SendSpotifyPlayerRequest(*sURL, http.MethodGet, queryParams, headers)
	if err != nil {
		if err.Error() == "invalid access token" {
			utils.RespondWithStatusUnavailable(w)
			return
		}
		log.Printf("ERROR in handleSearchByTrack: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByTrack: '%s' \n", err.Error()))
		return
	}

	searchByTrack, err := services.UnmarshalJSON[models.SearchByTrack](body)
	if err != nil {
		log.Printf("ERROR in handleSearchByTrack: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByTrack: '%s' \n", err.Error()))
		return
	}

	// Set header as JSON
	w.Header().Set("Content-Type", "application/json")
	// Create JSON encoder and encode our currentlyPlaying variable
	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(searchByTrack)
	if err != nil {
		log.Printf("ERROR in handleSearchByTrack: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByTrack: '%s' \n", err.Error()))
		return
	}
}

func handleSearchByURL(w http.ResponseWriter, r *http.Request) {
	// Get the URL from the GET request query param
	urlQueryParam := r.URL.Query().Get("url")
	if urlQueryParam == "" {
		log.Println("ERROR: 'urlQueryParam' was null", )
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByURL: '%s' \n", fmt.Errorf("'urlQueryParam' was null")))
		return
	}
	// Parse request URL
	sURL, err := url.Parse(urlQueryParam)
	if err != nil {
		log.Printf("ERROR in handleSearchByURL: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByURL: '%s' \n", err.Error()))
		return
	}
	// Grab our access token from middleware cache
	accessToken, err := middlewares.GetAccessToken()
	if err != nil {
		log.Printf("ERROR in handleSearchByURL: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByURL: '%s' \n", err.Error()))
		return
	}
	// Add our headers in a map
	headers := map[string]string {
		"Authorization": "Bearer "+accessToken,
	}
	// Call our spotify API function to get the response body
	body, err := services.SendSpotifyPlayerRequest(*sURL, http.MethodGet, nil, headers)
	if err != nil {
		if err.Error() == "invalid access token" {
			utils.RespondWithStatusUnavailable(w)
			return
		}
		log.Printf("ERROR in handleSearchByURL: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByURL: '%s' \n", err.Error()))
		return
	}
	searchByURL, err := services.UnmarshalJSON[models.SearchByTrack](body)
	if err != nil {
		log.Printf("ERROR in handleSearchByURL: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByURL: '%s' \n", err.Error()))
		return
	}
	// Set header as JSON
	w.Header().Set("Content-Type", "application/json")
	// Create JSON encoder and encode our currentlyPlaying variable
	jsonEncoder := json.NewEncoder(w)
	err = jsonEncoder.Encode(searchByURL)
	if err != nil {
		log.Printf("ERROR in handleSearchByURL: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByURL: '%s' \n", err.Error()))
		return
	}
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	var err error;
	// Parse our request URL
	sURL, err := url.Parse(config.SpotifyPlayerURL + "/play")
	if err != nil {
		log.Printf("ERROR in handlePlay: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handlePlay: '%s' \n", err.Error()))
		return
	}
	// Grab our access token from middleware cache
	accessToken, err := middlewares.GetAccessToken()
	if err != nil {
		log.Printf("ERROR in handlePlay: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handlePlay: '%s' \n", err.Error()))
		return
	}
	// ADd our headers in a map
	headers := map[string]string {
		"Authorization": "Bearer "+accessToken,
		"Content-Type": "application/json",
	}
	// Call spotify API function to get response body
	_, err = services.SendSpotifyPlayerRequest(*sURL, http.MethodPut, nil, headers)
	if err != nil {
		if err.Error() == "invalid access token" {
			utils.RespondWithStatusUnavailable(w)
			return
		}
		log.Printf("ERROR in handlePlay: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handlePlay: '%s' \n", err.Error()))
		return
	}
}

func handlePause(w http.ResponseWriter, r *http.Request) {
	var err error;
	// Parse our request URL
	sURL, err := url.Parse(config.SpotifyPlayerURL + "/pause")
	if err != nil {
		log.Printf("ERROR in handlePause: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handlePause: '%s' \n", err.Error()))
		return
	}
	// Grab our access token from middleware cache
	accessToken, err := middlewares.GetAccessToken()
	if err != nil {
		log.Printf("ERROR in handlePause: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handlePause: '%s' \n", err.Error()))
		return
	}
	// ADd our headers in a map
	headers := map[string]string {
		"Authorization": "Bearer "+accessToken,
		"Content-Type": "application/json",
	}
	// Call spotify API function to get response body
	_, err = services.SendSpotifyPlayerRequest(*sURL, http.MethodPut, nil, headers)
	if err != nil {
		if err.Error() == "invalid access token" {
			utils.RespondWithStatusUnavailable(w)
			return
		}
		log.Printf("ERROR in handlePause: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handlePause: '%s' \n", err.Error()))
		return
	}
}

func handleAddToQueue(w http.ResponseWriter, r *http.Request) {
	var err error;
	trackURI := r.URL.Query().Get("track_uri")
	if trackURI == "" {
		log.Println("ERROR: 'trackURI' was null", )
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleAddToQueue: '%s' \n", fmt.Errorf("'urlQueryParam' was null")))
		return
	}
	// Parse our request URL
	sURL, err := url.Parse(config.SpotifyPlayerURL + "/queue")
	if err != nil {
		log.Printf("ERROR in handleAddToQueue: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleAddToQueue: '%s' \n", err.Error()))
		return
	}
	// Grab our access token from middleware cache
	accessToken, err := middlewares.GetAccessToken()
	if err != nil {
		log.Printf("ERROR in handleAddToQueue: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleAddToQueue: '%s' \n", err.Error()))
		return
	}
	// Add our headers in a map
	headers := map[string]string {
		"Authorization": "Bearer "+accessToken,
	}
	// Add our query params in a map
	queryParams := map[string]string {
		"uri": trackURI,
	}
	// Call spotify API function to get response body
	_, err = services.SendSpotifyPlayerRequest(*sURL, http.MethodPost, queryParams, headers)
	if err != nil {
		if err.Error() == "invalid access token" {
			utils.RespondWithStatusUnavailable(w)
			return
		}
		log.Printf("ERROR in handleAddToQueue: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleAddToQueue: '%s' \n", err.Error()))
		return
	}
}