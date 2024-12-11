package middlewares

import (
	"database/sql"
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
	"github.com/baywiggins/qIt-backend/internal/models"
	"github.com/baywiggins/qIt-backend/pkg/utils"
)

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	Scope string `json:"scope"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var tokenCache = cache.New(3600*time.Second, 12*time.Hour)

func GetSpotifyAuthURL(state string) (string, error) {
	var err error;

	apiURL := "http://" + config.API_URL
	var authURL string = config.SpotifyAuthURL + "/authorize"

	// Params of our auth URL we will be returning
	var params = map[string]string{
		"client_id": config.ClientID,
		"response_type": "code",
		"redirect_uri": apiURL + "/spotify/auth/callback",
		"state": state,
		"scope": config.Scopes,
		"show_dialog": "true",
	}

	fmt.Println(params)

	u, err := url.Parse(authURL)
	if err != nil {
		return "", fmt.Errorf("error in GetSpotifyAuthURL: '%s'", err.Error())
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	fmt.Println("auth url: " + u.String())

	return u.String(), err
}

// Function for Spotify controller endpoints to use
func GetAccessToken() (string, error) {
	var err error;

	accessToken, accessFound := tokenCache.Get("accessToken")
	_, refreshFound := tokenCache.Get("refreshToken")

	if !accessFound && !refreshFound {
		// If neither are found, user needs to authenticate with Spotify
		return "", errors.New("user must authenticate with spotify first")
	} else if !accessFound && refreshFound{
		// If access token not found, but refresh token is, get token from refresh function
		tokenData, err := GetRefreshToken()
		if err != nil {
			return "", fmt.Errorf("error in GetAccessToken: '%s'", err.Error())
		}
		return tokenData.AccessToken, err
	} else {
		return accessToken.(string), err
	}
}

func GetAccessTokenFromSpotify(code string, state string, db *sql.DB) error {
	var err error;
	// Get the URL to send POST request to
	tokenURL := config.SpotifyAuthURL + "/api/token"
	data := url.Values{}
	// Add query params for the POST request
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "http://" + config.API_URL + "/spotify/auth/callback")

	// Get our token data JSON
	tokenData, err := sendPostRequestWithAuthorization(tokenURL, data)
	if err != nil {
		return fmt.Errorf("error in GetAccessTokenFromSpotify: '%s'", err.Error())
	}

	tokenCache.Set("accessToken", tokenData.AccessToken, cache.DefaultExpiration)
	tokenCache.Set("refreshToken", tokenData.RefreshToken, cache.DefaultExpiration)

	encryptedAuthToken, err := utils.Encrypt(tokenData.AccessToken)
	if err != nil {
		return fmt.Errorf("error in GetAccessTokenFromSpotify: '%s'", err.Error())
	}
	encryptedRefreshToken, err := utils.Encrypt(tokenData.RefreshToken)
	if err != nil {
		return fmt.Errorf("error in GetAccessTokenFromSpotify: '%s'", err.Error())
	}
	// Add to database
	sta := models.StateToAuth{
		UserState: state,
		AuthToken: encryptedAuthToken,
		RefreshToken: encryptedRefreshToken,
		ExpirationTime: time.Now().Add(time.Hour),
	}
	if err := models.InsertStateToAuth(db, sta); err != nil {
		return fmt.Errorf("error in GetAccessTokenFromSpotify: '%s'", err.Error())
	}
	
	// If token obtained successfully, update user row to reflect account creation finish
	models.UpdateCreationStatusByState(db, state)
	
	return err
}

func GetRefreshToken() (SpotifyTokenResponse, error) {
	var err error;
	// TODO Make sure there is a refresh token in the cache, and if not, deal with normal auth flow
	refreshToken, found := tokenCache.Get("refreshToken")
	if !found {
		return SpotifyTokenResponse{}, errors.New("refreshToken not found in cache in GetRefreshToken")
	}
	// Get the URL to send POST request to
	tokenURL := config.SpotifyAuthURL + "/api/token"
	data := url.Values{}
	// Add query params to POST request
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken.(string))

	refreshTokenData, err := sendPostRequestWithAuthorization(tokenURL, data)
	if err != nil {
		return SpotifyTokenResponse{}, fmt.Errorf("error in GetRefreshToken: '%s'", err.Error())
	}

	// Update our cache again with new token
	tokenCache.Set("accessToken", refreshTokenData.AccessToken, cache.DefaultExpiration)
	if refreshTokenData.RefreshToken != "" {
		tokenCache.Set("refreshToken", refreshTokenData.RefreshToken, cache.NoExpiration)
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
		return SpotifyTokenResponse{}, fmt.Errorf("error in sendPostRequestWithAuthorization: Spotify responded with error: %s", body)
	}
	return tokenData, err
}