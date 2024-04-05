package network_test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/common/network"
)

func TestParseUrl(t *testing.T) {
	type args struct {
		inputURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid URL",
			args: args{
				inputURL: "https://raw.githubusercontent.com/defenseunicorns/go-oscal/main/docs/adr/0001-record-architecture-decisions.md",
			},
			wantErr: false,
		},
		{
			name: "invalid url",
			args: args{
				inputURL: "backmatter/resources",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := network.ParseUrl(tt.args.inputURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUrl() error = %v, wantErr %v", err, tt.wantErr)
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
			name:    "Not found",
			url:     "https://raw.githubusercontent.com/defenseunicorns/go-oscal/main/docs/adr/0000-record-architecture-decisions.md",
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
