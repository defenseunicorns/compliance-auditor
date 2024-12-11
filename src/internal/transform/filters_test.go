package transform_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/defenseunicorns/lula/src/internal/transform"
)

func TestBuildFilters(t *testing.T) {
	runTest := func(t *testing.T, node []byte, pathParts []transform.PathPart, expected []yaml.Filter) {
		t.Helper()

		n := createRNode(t, node)

		filters, err := transform.BuildFilters(n, pathParts)
		require.NoError(t, err)

		require.Equal(t, expected, filters)
	}

	tests := []struct {
		name      string
		pathParts []transform.PathPart
		nodeBytes []byte
		expected  []yaml.Filter
	}{
		{
			name: "simple-path",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeMap, Value: "a"},
				{Type: transform.PartTypeScalar, Value: "b"},
			},
			expected: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.PathGetter{Path: []string{"b"}},
			},
		},
		{
			name: "filter-path",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeFilter, Value: "b=c"},
			},
			expected: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementMatcher{Keys: []string{"b"}, Values: []string{"c"}},
			},
		},
		{
			name: "complex-filter-path",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeFilter, Value: "b.c=d"},
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
		{
			name: "composite-multi-filter-path",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeFilter, Value: "b.c=d,e=f"},
			},
			nodeBytes: []byte(`
a:
  - b:
      c: d
    e: c
  - b:
      c: d
    e: f
`),
			expected: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementIndexer{Index: 1},
			},
		},
		{
			name: "multi-filter-path",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeFilter, Value: "b=d,e=f"},
			},
			nodeBytes: []byte(`
a:
  - b: d
    e: c
  - b: d
    e: f
`),
			expected: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementIndexer{Index: 1},
			},
		},
		{
			name: "encapsulated-key",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeMap, Value: "a"},
				{Type: transform.PartTypeScalar, Value: "b.c=d"},
			},
			expected: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.PathGetter{Path: []string{"b.c=d"}},
			},
		},
		{
			name: "encapsulated-filter-path",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeFilter, Value: `"b.c"=d`},
			},
			expected: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementMatcher{Keys: []string{"b.c"}, Values: []string{"d"}},
			},
		},
		{
			name: "root-path",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeScalar, Value: ""},
			},
			expected: []yaml.Filter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.nodeBytes, tt.pathParts, tt.expected)
		})
	}
}
