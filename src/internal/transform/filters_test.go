package transform_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/defenseunicorns/lula/src/internal/transform"
)

func TestBuildFilters(t *testing.T) {
	runTest := func(t *testing.T, node []byte, pathSlice []string, expected []yaml.Filter) {
		t.Helper()

		n := createRNode(t, node)

		filters, err := transform.BuildFilters(n, pathSlice)
		require.NoError(t, err)

		require.Equal(t, expected, filters)
	}

	tests := []struct {
		name      string
		pathSlice []string
		nodeBytes []byte
		expected  []yaml.Filter
	}{
		{
			name: "test-simple-path",
			pathSlice: []string{
				"a",
				"b",
			},
			expected: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.PathGetter{Path: []string{"b"}},
			},
		},
		{
			name: "test-sequence-path",
			pathSlice: []string{
				"a",
				"[b=c]",
			},
			expected: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementMatcher{Keys: []string{"b"}, Values: []string{"c"}},
			},
		},
		{
			name: "test-composite-path",
			pathSlice: []string{
				"a",
				"[b.c=d]",
			},
			nodeBytes: []byte(`
a:
  - b: 
      c: d
  - b:
      c: e
`),
			expected: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementIndexer{Index: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.nodeBytes, tt.pathSlice, tt.expected)
		})
	}
}
