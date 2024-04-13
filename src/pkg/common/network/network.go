package network

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	// If the URL is a file, fetch the file from the local filesystem
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
		err = ValidateChecksum(bytes, checksum)
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

// ValidateChecksum validates a given checksum against a given []bytes.
// Supports MD5, SHA-1, SHA-256, and SHA-512.
// Returns an error if the hash does not match.
func ValidateChecksum(data []byte, expectedChecksum string) error {
	var actualChecksum string
	switch len(expectedChecksum) {
	case md5.Size * 2:
		hash := md5.Sum(data)
		actualChecksum = hex.EncodeToString(hash[:])
	case sha1.Size * 2:
		hash := sha1.Sum(data)
		actualChecksum = hex.EncodeToString(hash[:])
	case sha256.Size * 2:
		hash := sha256.Sum256(data)
		actualChecksum = hex.EncodeToString(hash[:])
	case sha512.Size * 2:
		hash := sha512.Sum512(data)
		actualChecksum = hex.EncodeToString(hash[:])
	default:
		return errors.New("unsupported checksum type")
	}

	if actualChecksum != expectedChecksum {
		return errors.New("checksum validation failed")
	}

	return nil
}
