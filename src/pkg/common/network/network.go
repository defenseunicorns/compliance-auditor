package network

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/defenseunicorns/lula/src/pkg/common"
)

var HttpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// parseUrl parses a URL string into a url.URL object.
func parseUrl(inputURL string) (*url.URL, error) {
	parsedUrl, err := url.ParseRequestURI(inputURL)
	if err != nil || parsedUrl.Scheme == "" || (parsedUrl.Scheme != "file" && parsedUrl.Host == "") {
		return nil, errors.New("invalid URL")
	}
	return parsedUrl, nil
}

// ParseChecksum parses a URL string into a url.URL object.
// If the URL has a checksum, the checksum is removed from the URL and returned.
func ParseChecksum(src string) (*url.URL, string, error) {
	atSymbolCount := strings.Count(src, "@")
	var checksum string
	if atSymbolCount > 0 {
		parsed, err := parseUrl(src)
		if err != nil {
			return parsed, checksum, fmt.Errorf("unable to parse the URL: %s", src)
		}
		if atSymbolCount == 1 && parsed.User != nil {
			return parsed, checksum, nil
		}

		index := strings.LastIndex(src, "@")
		checksum = src[index+1:]
		src = src[:index]
	}

	url, err := parseUrl(src)
	if err != nil {
		return url, checksum, err
	}

	return url, checksum, nil
}

// Fetch fetches the response body from a given URL after validating it.
// If the URL scheme is "file", the file is fetched from the local filesystem.
// If the URL scheme is "http", "https", or "ftp", the file is fetched from the remote server.
// If the URL has a checksum, the file is validated against the checksum.
func Fetch(inputURL string) (bytes []byte, err error) {
	url, checksum, err := ParseChecksum(inputURL)
	if err != nil {
		return bytes, err
	}

	if url.Scheme == "file" {
		bytes, err = FetchLocalFile(url)
		if err != nil {
			return bytes, err
		}
	} else {

		// Make the HTTP GET request
		resp, err := HttpClient.Get(inputURL)
		if err != nil {
			return bytes, err
		}
		defer resp.Body.Close()

		// Read the response body
		bytes, err = io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return bytes, errors.New("received non-OK HTTP status")
		}
		if err != nil {
			return bytes, err
		}
	}

	if checksum != "" {
		// Validate the bytes against the SHA
		err = common.ValidateChecksum(bytes, checksum)
		if err != nil {
			return bytes, err
		}
	}

	return bytes, nil
}

// FetchLocalFile fetches a local file from a given URL.
// If the URL scheme is not "file", an error is returned.
// If the URL is relative, the component definition directory is prepended if set, otherwise the current working directory is prepended.
func FetchLocalFile(url *url.URL) ([]byte, error) {
	if url.Scheme != "file" {
		return nil, errors.New("expected file URL scheme")
	}
	requestUri := url.RequestURI()

	// If the request uri is absolute, use it directly
	if _, err := os.Stat(requestUri); err != nil {
		// if relative pre-pend cwd
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		requestUri = filepath.Join(cwd, requestUri)
	}

	bytes, err := os.ReadFile(requestUri)
	return bytes, err
}
