package config

import (
	"fmt"
	"os"
)

// Define global variables for configuration
var (
    SpotifyAuthURL   string
    SpotifyPlayerURL string
    SpotifySearchURL string
    ClientID         string
    API_URL          string
    ClientSecret     string
    Scopes           string
)

// init function to load environment variables
func init() {
    SpotifyAuthURL = os.Getenv("SPOTIFY_AUTH_URL")
    SpotifyPlayerURL = os.Getenv("SPOTIFY_PLAYER_URL")
    SpotifySearchURL = os.Getenv("SPOTIFY_SEARCH_URL")
    ClientID = os.Getenv("CLIENT_ID")
    API_URL = os.Getenv("API_URL")
    ClientSecret = os.Getenv("CLIENT_SECRET")
    Scopes = os.Getenv("SCOPES")

    envarList := [7]string{SpotifyAuthURL, SpotifyPlayerURL, SpotifySearchURL, ClientID, API_URL, ClientSecret, Scopes}

    for i := range envarList {
        if envarList[i] == "" {
            fmt.Printf("Envar of index %d is null \n", i)
            os.Exit(1)
        }
    }
}