package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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
    DBName           string
    AESKey           string
    RedisAddress     string
)

// init function to load environment variables
func init() {
    err := godotenv.Load()

    if err != nil {
        fmt.Println("Unable to load dotenv")
        os.Exit(1)
    }

    SpotifyAuthURL   = os.Getenv("SPOTIFY_AUTH_URL")
    SpotifyPlayerURL = os.Getenv("SPOTIFY_PLAYER_URL")
    SpotifySearchURL = os.Getenv("SPOTIFY_SEARCH_URL")
    ClientID         = os.Getenv("CLIENT_ID")
    API_URL          = os.Getenv("API_URL")
    ClientSecret     = os.Getenv("CLIENT_SECRET")
    Scopes           = os.Getenv("SCOPES")
    DBName           = os.Getenv("DB_NAME")
    AESKey           = os.Getenv("AES_KEY")
    RedisAddress     = os.Getenv("REDIS_ADDRESS")

    envarList := [10]string{SpotifyAuthURL, SpotifyPlayerURL, SpotifySearchURL, ClientID, API_URL, ClientSecret, Scopes, DBName, AESKey, RedisAddress}

    for i := range envarList {
        if envarList[i] == "" {
            fmt.Printf("Envar of index %d is null \n", i)
            os.Exit(1)
        }
    }
}