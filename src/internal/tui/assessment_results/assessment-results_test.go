package assessmentresults_test

import (
	"testing"

	assessmentresults "github.com/defenseunicorns/lula/src/internal/tui/assessment_results"
)

// TODO: add test for GetResults

func TestGetReadableDesc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		desc     string
		expected string
	}{
		{
			name:     "Test get readable desc",
			desc:     "[TEST]: 67456ae8-4505-4c93-b341-d977d90cb125 - istio-health-check",
			expected: "istio-health-check",
		},
		{
			name:     "Test get readable desc - no uuid",
			desc:     "test description",
			expected: "test description",
		},
		{
			name:     "Test get readable desc - no description",
			desc:     "[TEST]: 12345678-1234-1234-1234-123456789012",
			expected: "[TEST]: 12345678-1234-1234-1234-123456789012",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assessmentresults.GetReadableObservationName(tt.desc)
			if got != tt.expected {
				t.Errorf("GetReadableObservationName() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
