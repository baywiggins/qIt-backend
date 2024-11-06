package services

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func SendSpotifyPlayerRequest(u url.URL, httpMethod string, queryParams map[string]string, headers map[string]string) ([]byte, error) {
	var err error;

	// Add query params to URL
	query := u.Query()
	for k, v := range queryParams {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()

	// Create request
	req, err := http.NewRequest(httpMethod, u.String(), nil)
	if err != nil {
		fmt.Println("Add error handling pls lol")
		return nil, err
	}

	// Add headers to request
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Add error handling plsssss looooool")
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("pls bro")
		return nil, err
	}

	return body, err
}

func SendSpotifySearchRequest () {
	return
}