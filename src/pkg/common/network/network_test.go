package network_test

import (
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
			name:         "invalid url",
			input:        "backmatter/resources",
			wantErr:      true,
			wantChecksum: false,
		},
		{
			name:         "File url",
			input:        "file://../../../../test/e2e/scenarios/remote-validations/validation.opa.yaml",
			wantErr:      false,
			wantChecksum: false,
		},
		{
			name:         "With Checksum",
			input:        "file://../../../../test/e2e/scenarios/remote-validations/validation.opa.yaml@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
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
			url:     "file://../../../../test/e2e/scenarios/remote-validations/validation.opa.yaml",
			wantErr: false,
		},
		{
			name:    "File with checksum SHA-256",
			url:     "file://../../../../test/e2e/scenarios/remote-validations/validation.opa.yaml@6f76e9e7b0e936d4c03ea2ece54c7d250d9d53d906d8ca325c64697435802c17",
			wantErr: false,
		},
		{
			name:    "File with checksum",
			url:     "file://../../../../test/e2e/scenarios/remote-validations/validation.opa.yaml@d8a6ba5ebccde299ecef1e54352e0eded97a6738",
			wantErr: false,
		},
		{
			name:    "Not found",
			url:     "https://raw.githubusercontent.com/defenseunicorns/go-oscal/main/docs/adr/0000-record-architecture-decisions.md",
			wantErr: true,
		},
		{
			name:    "Invalid Sha",
			url:     "file://../../../../test/e2e/scenarios/remote-validations/validation.opa.yaml@2d4c18916f2fd70f9488b76690c2eed06789d5fd12e06152a01a8ef7600c41ef",
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
