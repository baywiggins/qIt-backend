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
func GetAccessToken(uuid string, db *sql.DB) (string, error) {
	var err error;

	accessToken, refreshToken, expiresIn, err := models.GetStateToAuthRowByID(db, uuid)

	if err != nil {
		// error in db function
		return "", fmt.Errorf("error in GetStateToAuthRowByID: '%s'", err)
	}

	// Check if token is expired
	if time.Now().UTC().After(expiresIn) {
		fmt.Println("token is expired yipee")
		tr, err := getRefreshToken(refreshToken, uuid, db)
		fmt.Println(tr)

		return tr.AccessToken, err
	}

	if (accessToken == "" || refreshToken == "") {
		// If either aren't found, user needs to authenticate with Spotify
		return "", errors.New("user must authenticate with spotify first")
	} else {
		return accessToken, err
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
		ExpirationTime: time.Now().UTC().Add(time.Hour).Format(time.RFC3339),
	}
	if err := models.InsertStateToAuth(db, sta); err != nil {
		return fmt.Errorf("error in GetAccessTokenFromSpotify: '%s'", err.Error())
	}
	
	// If token obtained successfully, update user row to reflect account creation finish
	models.UpdateCreationStatusByState(db, state)
	
	return err
}

func getRefreshToken(refreshToken string, uuid string, db *sql.DB) (SpotifyTokenResponse, error) {
	var err error;
	// Get the URL to send POST request to
	tokenURL := config.SpotifyAuthURL + "/api/token"
	data := url.Values{}
	// Add query params to POST request
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	refreshTokenData, err := sendPostRequestWithAuthorization(tokenURL, data)
	if err != nil {
		return SpotifyTokenResponse{}, fmt.Errorf("error in GetRefreshToken: '%s'", err.Error())
	}

	// Check if new refresh token was provided in the request, if not just keep using the one we have
	newRefreshToken := refreshTokenData.RefreshToken
	if refreshTokenData.RefreshToken == "" {
		newRefreshToken = refreshToken
	}
	// Get new expiry time (1 hour after obtaining it)
	expirationTime := time.Now().Add(time.Hour).Format(time.RFC3339)

	// Update auth token info in DB
	models.UpdateAccessTokenByID(db, uuid, refreshTokenData.AccessToken, newRefreshToken, expirationTime)

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