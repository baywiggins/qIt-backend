package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/baywiggins/qIt-backend/internal/api/middlewares"
	"github.com/baywiggins/qIt-backend/internal/config"
	"github.com/baywiggins/qIt-backend/internal/services"
)

// REMEMBER TO ACCOUNT FOR CONDITION IF TOKEN EXPIRES BETWEEN API CALLS
// IF EVER GET ACCESS DENIED FROM SPOTIFY API JUST REFRESH THE TOKEN!!!!!

func HandleSpotifyControllerRoutes() {
	http.Handle("GET /spotify/currently-playing", middlewares.LoggingMiddleware(http.HandlerFunc(handleCurrentlyPlaying)))
	http.Handle("GET /test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sURL := "https://api.spotify.com/v1/me/player/play"
		authCode, err := middlewares.GetAccessToken()
		if err != nil{
			fmt.Println("err getting authcode")
			return
		}

		u, err := url.Parse(sURL)
		if err != nil {
			fmt.Println("err parsing URL")
			return
		}
		query := u.Query()
		u.RawQuery = query.Encode()

		req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer([]byte{}))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Authorization", "Bearer "+authCode)
		req.Header.Set("Content-Type", "application/json")

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Print response status
		fmt.Println("Response status:", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))

	}))

	http.Handle("GET /test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sURL := "https://api.spotify.com/v1/me/player/pause"
		authCode, err := middlewares.GetAccessToken()
		if err != nil{
			fmt.Println("err getting authcode")
			return
		}

		u, err := url.Parse(sURL)
		if err != nil {
			fmt.Println("err parsing URL")
			return
		}
		query := u.Query()
		u.RawQuery = query.Encode()

		req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer([]byte{}))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Authorization", "Bearer "+authCode)
		req.Header.Set("Content-Type", "application/json")

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Print response status
		fmt.Println("Response status:", resp.Status)
	}))
}

func handleCurrentlyPlaying(w http.ResponseWriter, r *http.Request) {
	var err error;
	sURL, err := url.Parse(config.SpotifyPlayerURL)
	if err != nil {
		fmt.Println("handddd;e;e")
	}
	accessToken, err := middlewares.GetAccessToken()
	if err != nil {
		fmt.Println("pls handle this")
	}
	
	headers := map[string]string {
		"Authorization": "Bearer "+accessToken,
	}

	body, err := services.SendSpotifyPlayerRequest(*sURL, http.MethodGet, nil, headers)
	if err != nil {
		fmt.Println("handsadasdasdddd;e;e")
	}

	fmt.Println(string(body))
}