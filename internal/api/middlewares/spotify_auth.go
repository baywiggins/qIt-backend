package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/baywiggins/qIt-backend/internal/config"
)

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	Scope string `json:"scope"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var tokenCache = cache.New(3600*time.Second, 12*time.Hour)

func GetSpotifyAuthURL() (string, error) {
	var err error;

	apiURL := "http://" + config.API_URL
	var authURL string = config.SpotifyURL + "/authorize"

	// Params of our auth URL we will be returning
	// Change values to be gotten from the environment
	var params = map[string]string{
		"client_id": config.ClientID,
		"response_type": "code",
		"redirect_uri": apiURL + "/spotify/auth/callback",
		"scope": config.Scopes,
		"show_dialog": "true",
	}

	u, err := url.Parse(authURL)
	if err != nil {
		return "", fmt.Errorf("error in GetSpotifyAuthURL: '%s'", err.Error())
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), err
}

func GetAccessToken() string {
	accessToken, accessFound := tokenCache.Get("accessToken")
	refreshToken, refreshFound := tokenCache.Get("refreshToken")
}

func GetAccessTokenFromSpotify(code string) (SpotifyTokenResponse, error) {
	var err error;
	// Get the URL to send POST request to
	tokenURL := config.SpotifyURL + "/api/token"
	data := url.Values{}
	// Add query params for the POST request
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "http://" + config.API_URL + "/spotify/auth/callback")

	// Get our token data JSON
	tokenData, err := sendPostRequestWithAuthorization(tokenURL, data)
	if err != nil {
		return SpotifyTokenResponse{}, fmt.Errorf("error in GetAccessTokenFromSpotify: '%s'", err.Error())
	}

	// Add to cache
	tokenCache.Set("accessToken", tokenData.AccessToken, cache.DefaultExpiration)
	tokenCache.Set("refreshToken", tokenData.RefreshToken, cache.NoExpiration)
	
	return tokenData, err
}

func GetRefreshToken() (SpotifyTokenResponse, error) {
	var err error;
	// TODO Make sure there is a refresh token in the cache, and if not, deal with normal auth flow
	refreshToken, found := tokenCache.Get("refreshToken")
	if !found {
		return SpotifyTokenResponse{}, errors.New("refreshToken not found in cache in GetRefreshToken")
	}
	// Get the URL to send POST request to
	tokenURL := config.SpotifyURL + "/api/token"
	data := url.Values{}
	// Add query params to POST request
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken.(string))

	refreshTokenData, err := sendPostRequestWithAuthorization(tokenURL, data)
	if err != nil {
		return SpotifyTokenResponse{}, fmt.Errorf("error in GetRefreshToken: '%s'", err.Error())
	}

	return refreshTokenData, err
}

func sendPostRequestWithAuthorization(url string, queryData url.Values) (SpotifyTokenResponse, error) {
	var err error;

	req, err := http.NewRequest("POST", url, strings.NewReader(queryData.Encode()))
	if err != nil {
		return SpotifyTokenResponse{}, fmt.Errorf("error in sendPostRequestWithAuthorization: '%s'", err.Error())
	}

	clientID := config.ClientID
	clientSecret := config.ClientSecret
	
	// Concatenate client ID and client secret with a colon
    credentials := clientID + ":" + clientSecret

    // Base64 encode the credentials
    encodedCredentials := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Set("Authorization", encodedCredentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return SpotifyTokenResponse{}, fmt.Errorf("error in sendPostRequestWithAuthorization: '%s'", err.Error())
	}

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SpotifyTokenResponse{}, fmt.Errorf("error in sendPostRequestWithAuthorization: '%s'", err.Error())
	}

	var tokenData SpotifyTokenResponse
	if err := json.Unmarshal(body, &tokenData); err != nil {
		return SpotifyTokenResponse{}, fmt.Errorf("error in sendPostRequestWithAuthorization: '%s'", err.Error())
	}
	
	if resp.StatusCode != 200 {
		return SpotifyTokenResponse{}, fmt.Errorf("error in sendPostRequestWithAuthorization: Spotify responded with error")
	}
	return tokenData, err
}