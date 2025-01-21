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

// SpotifyRequestParams encapsulates common parameters for Spotify API requests.
type SpotifyRequestParams struct {
	Path        string
	Method      string
	QueryParams map[string]string
	Headers     map[string]string
}

// HandleSpotifyControllerRoutes registers all Spotify-related routes.
func HandleSpotifyControllerRoutes(db *sql.DB) {
	handler := &Handler{DB: db}

	http.Handle("/spotify/playback-state", middlewares.LoggingMiddleware(http.HandlerFunc(handler.withSpotify(handler.handlePlaybackState, ""))))
	http.Handle("/spotify/currently-playing", middlewares.LoggingMiddleware(http.HandlerFunc(handler.withSpotify(handler.handleCurrentlyPlaying, "/currently-playing"))))
	http.Handle("/spotify/queue", middlewares.LoggingMiddleware(http.HandlerFunc(handler.withSpotify(handler.handleGetQueue, "/queue"))))
	http.Handle("/spotify/search/track", middlewares.LoggingMiddleware(http.HandlerFunc(handler.withSpotify(handler.handleSearchByTrack, ""))))
	http.Handle("/spotify/search/url", middlewares.LoggingMiddleware(http.HandlerFunc(handler.withSpotify(handler.handleSearchByURL, ""))))
	http.Handle("/spotify/play", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(handler.withSpotify(handler.handlePlay, "/play")))))
	http.Handle("/spotify/pause", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(handler.withSpotify(handler.handlePause, "/pause")))))
	http.Handle("/spotify/add-to-queue", middlewares.LoggingMiddleware(http.HandlerFunc(handler.withSpotify(handler.handleAddToQueue, "/queue"))))
}

// withSpotify is a wrapper for Spotify API endpoints with common setup.
func (h *Handler) withSpotify(handlerFunc func(http.ResponseWriter, *http.Request, SpotifyRequestParams), path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("uuid")
		if uuid == "" {
			log.Printf("Error in withSpotify: uuid not provided in request")
			utils.RespondWithError(w, http.StatusBadRequest, "uuid not provided in request")
			return
		}

		accessToken, err := middlewares.GetAccessToken(uuid, h.DB)
		if err != nil {
			log.Printf("Error in withSpotify: '%s'", err.Error())
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get access token")
			return
		}

		params := SpotifyRequestParams{
			Path:    path,
			Method:  http.MethodGet,
			Headers: map[string]string{"Authorization": "Bearer " + accessToken},
		}
		handlerFunc(w, r, params)
	}
}

// sendSpotifyRequest is a reusable helper to interact with the Spotify API.
func sendSpotifyRequest(params SpotifyRequestParams, search bool) ([]byte, error) {
	baseURL := config.SpotifyPlayerURL
	if (search) {
		baseURL = config.SpotifySearchURL
	}
	sURL, err := url.Parse(baseURL + params.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	return services.SendSpotifyPlayerRequest(*sURL, params.Method, params.QueryParams, params.Headers)
}

// handlePlaybackState handles the current state of playback
func (h* Handler) handlePlaybackState(w http.ResponseWriter, r *http.Request, params SpotifyRequestParams) {
	body, err := sendSpotifyRequest(params, false)
	if err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}


	var playbackState models.PlaybackState
	if err := json.Unmarshal(body, &playbackState); err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}

	// Set the correct header for JSON response
	w.Header().Set("Content-Type", "application/json")

	// Encode the parsed JSON back to the response writer
	if err := json.NewEncoder(w).Encode(playbackState); err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}
}

// handleCurrentlyPlaying handles the currently playing endpoint.
func (h *Handler) handleCurrentlyPlaying(w http.ResponseWriter, r *http.Request, params SpotifyRequestParams) {
	body, err := sendSpotifyRequest(params, false)
	if err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}

	var currentlyPlaying models.CurrentlyPlaying
	if err := json.Unmarshal(body, &currentlyPlaying); err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}

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

// handleGetQueue handles the queue retrieval endpoint.
func (h *Handler) handleGetQueue(w http.ResponseWriter, r *http.Request, params SpotifyRequestParams) {
	body, err := sendSpotifyRequest(params, false)
	if err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}

	currentQueue, err := services.UnmarshalJSON[models.CurrentQueue](body)
	if err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(currentQueue); err != nil {
		log.Printf("ERROR in handleGetQueue: %s\n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleGetQueue: '%s'\n", err.Error()))
	}
}

// handleSearchByTrack handles track search requests.
func (h *Handler) handleSearchByTrack(w http.ResponseWriter, r *http.Request, params SpotifyRequestParams) {
	track := r.URL.Query().Get("track")
	if track == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing 'track' parameter")
		return
	}
	params.QueryParams = map[string]string{
		"q":    "track:" + track,
		"type": "track",
		"limit": r.URL.Query().Get("limit"),
	}
	body, err := sendSpotifyRequest(params, true)
	if err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}

	searchByTrack, err := services.UnmarshalJSON[models.SearchByTrack](body)
	if err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(searchByTrack); err != nil {
		log.Printf("ERROR in handleSearchByTrack: %s\n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByTrack: '%s'\n", err.Error()))
	}
}

// handleSearchByURL handles searches by URL.
func (h *Handler) handleSearchByURL(w http.ResponseWriter, r *http.Request, params SpotifyRequestParams) {
	urlQueryParam := r.URL.Query().Get("url")
	if urlQueryParam == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing 'url' parameter")
		return
	}
	parsedURL, err := url.Parse(urlQueryParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid 'url' parameter")
		return
	}
	params.Path = parsedURL.Path
	body, err := sendSpotifyRequest(params, true)
	if err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}

	searchByURL, err := services.UnmarshalJSON[models.SearchByTrack](body)
	if err != nil {
		utils.HandleSpotifyError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(searchByURL); err != nil {
		log.Printf("ERROR in handleSearchByURL: %s\n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleSearchByURL: '%s'\n", err.Error()))
	}
}

// handlePlay handles playback start.
func (h *Handler) handlePlay(w http.ResponseWriter, r *http.Request, params SpotifyRequestParams) {
	params.Method = http.MethodPut
	params.Headers["Content-Type"] = "application/json"
	if _, err := sendSpotifyRequest(params, false); err != nil {
		utils.HandleSpotifyError(w, err)
	}
}

// handlePause handles playback pause.
func (h *Handler) handlePause(w http.ResponseWriter, r *http.Request, params SpotifyRequestParams) {
	params.Method = http.MethodPut
	params.Headers["Content-Type"] = "application/json"
	if _, err := sendSpotifyRequest(params, false); err != nil {
		utils.HandleSpotifyError(w, err)
	}
}

// handleAddToQueue adds a track to the Spotify queue.
func (h *Handler) handleAddToQueue(w http.ResponseWriter, r *http.Request, params SpotifyRequestParams) {
	trackURI := r.URL.Query().Get("track_uri")
	if trackURI == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing 'track_uri' parameter")
		return
	}
	params.QueryParams = map[string]string{"uri": trackURI}
	params.Method = http.MethodPost
	if _, err := sendSpotifyRequest(params, false); err != nil {
		utils.HandleSpotifyError(w, err)
	}
}
