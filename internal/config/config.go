package config

import (
	"fmt"
	"os"
)

// Define global variables for configuration
var (
    SpotifyURL   string
    ClientID     string
    API_URL      string
    ClientSecret string
    Scopes       string
)

// init function to load environment variables
func init() {
    SpotifyURL = os.Getenv("SPOTIFY_URL")
    ClientID = os.Getenv("CLIENT_ID")
    API_URL = os.Getenv("API_URL")
    ClientSecret = os.Getenv("CLIENT_SECRET")
    Scopes = os.Getenv("SCOPES")

    envarList := [5]string{SpotifyURL, ClientID, API_URL, ClientSecret, Scopes}

    for i := range envarList {
        if envarList[i] == "" {
            fmt.Printf("Envar of index %d is null \n", i)
            os.Exit(1)
        }
    }
}