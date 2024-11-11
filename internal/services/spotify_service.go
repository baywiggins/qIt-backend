package services

import (
	"encoding/json"
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
		return nil, fmt.Errorf("error in SendSpotifyPlayerRequest: '%s", err.Error())
	}

	// Add headers to request
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in SendSpotifyPlayerRequest: '%s", err.Error())
	}
	defer resp.Body.Close()

	// Read response body data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error in SendSpotifyPlayerRequest: '%s", err.Error())
	}

	return body, err
}

// JSON unmarshaling
func UnmarshalJSON[T any](body []byte) (T, error) {
    var result T
    err := json.Unmarshal(body, &result)
    if err != nil {
        return result, fmt.Errorf("error in UnmarshalJSON: '%s'", err.Error())
    }
    return result, nil
}