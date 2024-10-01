package network_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/common/network"
)

func TestParseUrl(t *testing.T) {

	tests := []struct {
		name         string
		input        string
		wantErr      bool
		wantChecksum bool
	}{
		{
			name:         "valid URL",
			input:        "https://raw.githubusercontent.com/defenseunicorns/go-oscal/main/docs/adr/0001-record-architecture-decisions.md",
			wantErr:      false,
			wantChecksum: false,
		},
		{
			name:         "Invalid URL scheme",
			input:        "ht@tp://example.com",
			wantErr:      true,
			wantChecksum: false,
		},
		{
			name:         "Empty URL",
			input:        "",
			wantErr:      true,
			wantChecksum: false,
		},
		{
			name:         "URL with spaces",
			input:        "http://example .com",
			wantErr:      true,
			wantChecksum: false,
		},
		{
			name:         "Adds file if no scheme",
			input:        "path/to/file",
			wantErr:      false,
			wantChecksum: false,
		},
		{
			name:         "File url",
			input:        "file://../../../test/e2e/scenarios/remote-validations/validation.opa.yaml",
			wantErr:      false,
			wantChecksum: false,
		},
		{
			name:         "With Checksum",
			input:        "file://../../../test/e2e/scenarios/remote-validations/validation.opa.yaml@394f5efa7aa5c3163a631d0f2640efe836af07c77fa7b27749f00819dd869058",
			wantErr:      false,
			wantChecksum: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, checksum, err := network.ParseChecksum(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (checksum != "") != tt.wantChecksum {
				t.Errorf("ParseChecksum() checksum = %v, want %v", checksum, tt.wantChecksum)
				return
			}
		})
	}
}

func TestFetch(t *testing.T) {

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "valid URL",
			url:     "https://raw.githubusercontent.com/defenseunicorns/go-oscal/main/docs/adr/0001-record-architecture-decisions.md",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			url:     "backmatter/resources",
			wantErr: true,
		},
		{
			name:    "File",
			url:     "file://../../../test/e2e/scenarios/remote-validations/validation.opa.yaml",
			wantErr: false,
		},
		{
			name:    "File with checksum SHA-256",
			url:     "file://../../../test/e2e/scenarios/remote-validations/validation.opa.yaml@394f5efa7aa5c3163a631d0f2640efe836af07c77fa7b27749f00819dd869058",
			wantErr: false,
		},
		{
			name:    "Not found",
			url:     "https://raw.githubusercontent.com/defenseunicorns/go-oscal/main/docs/adr/0000-record-architecture-decisions.md",
			wantErr: true,
		},
		{
			name:    "Invalid Sha",
			url:     "file://../../../test/e2e/scenarios/remote-validations/validation.opa.yaml@2d4c18916f2fd70f9488b76690c2eed06789d5fd12e06152a01a8ef7600c41ef",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := network.Fetch(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 && !tt.wantErr {
				t.Errorf("Expected response body, got %v", got)
			}
		})
	}
}

func TestValidateChecksum(t *testing.T) {
	tests := []struct {
		name             string
		data             []byte
		expectedChecksum string
		wantErr          bool
	}{
		{
			name:             "MD5 checksum",
			data:             []byte("test"),
			expectedChecksum: "098f6bcd4621d373cade4e832627b4f6",
			wantErr:          false,
		},
		{
			name:             "SHA-1 checksum",
			data:             []byte("test"),
			expectedChecksum: "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
			wantErr:          false,
		},
		{
			name:             "SHA-256 checksum",
			data:             []byte("test"),
			expectedChecksum: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			wantErr:          false,
		},
		{
			name:             "SHA-512 checksum",
			data:             []byte("test"),
			expectedChecksum: "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff",
			wantErr:          false,
		},
		{
			name:             "Invalid checksum type",
			data:             []byte("test"),
			expectedChecksum: "invalid",
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := network.ValidateChecksum(tt.data, tt.expectedChecksum)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestParseChecksum(t *testing.T) {
	tests := []struct {
		name             string
		inputURL         string
		expectedURL      *url.URL
		expectedChecksum string
		wantErr          bool
	}{
		{
			name:             "Valid URL with checksum",
			inputURL:         "https://example.com/file.txt@123456",
			expectedURL:      &url.URL{Scheme: "https", Host: "example.com", Path: "/file.txt"},
			expectedChecksum: "123456",
			wantErr:          false,
		},
		{
			name:             "Valid URL without checksum",
			inputURL:         "https://example.com/file.txt",
			expectedURL:      &url.URL{Scheme: "https", Host: "example.com", Path: "/file.txt"},
			expectedChecksum: "",
			wantErr:          false,
		},
		{
			name:             "Invalid URL",
			inputURL:         "",
			expectedURL:      nil,
			expectedChecksum: "",
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, checksum, err := network.ParseChecksum(tt.inputURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(url, tt.expectedURL) {
				t.Errorf("ParseChecksum() url = %v, want %v", url, tt.expectedURL)
			}
			if checksum != tt.expectedChecksum {
				t.Errorf("ParseChecksum() checksum = %v, want %v", checksum, tt.expectedChecksum)
			}
		})
	}
}
