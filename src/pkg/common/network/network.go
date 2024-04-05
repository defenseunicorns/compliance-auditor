package network

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

var HttpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// ParseUrl parses a URL string into a url.URL object.
func ParseUrl(inputURL string) (*url.URL, error) {
	url, err := url.ParseRequestURI(inputURL)
	if err != nil {
		return nil, errors.New("invalid URL")
	}
	return url, nil
}

// Fetch fetches the response body from a given URL after validating it.
func Fetch(inputURL string) ([]byte, error) {
	// Validate the URL
	_, err := ParseUrl(inputURL)
	if err != nil {
		return nil, err
	}

	// Make the HTTP GET request
	resp, err := HttpClient.Get(inputURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("received non-OK HTTP status")
	}
	if err != nil {
		return nil, err
	}

	return body, nil
}
