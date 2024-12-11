package transform_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/defenseunicorns/lula/src/internal/transform"
)

func TestResolvePathWithFilters(t *testing.T) {
	runTest := func(t *testing.T, node []byte, path string, expectedPathParts []transform.PathPart, expectedFilters []yaml.Filter) {
		t.Helper()

		n := createRNode(t, node)

		pathParts, filters, err := transform.ResolvePathWithFilters(n, path)
		require.NoError(t, err)

		require.Equal(t, expectedPathParts, pathParts)
		require.Equal(t, expectedFilters, filters)
	}

	tests := []struct {
		name              string
		path              string
		nodeBytes         []byte
		expectedPathParts []transform.PathPart
		expectedFilters   []yaml.Filter
	}{
		{
			name: "simple-path",
			path: "a.b",
			expectedPathParts: []transform.PathPart{
				{Type: transform.PartTypeMap, Value: "a"},
				{Type: transform.PartTypeScalar, Value: "b"},
			},
			expectedFilters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.PathGetter{Path: []string{"b"}},
			},
		},
		{
			name: "filter-path",
			path: "a[b=c]",
			expectedPathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeSelector, Value: "b=c"},
			},
			expectedFilters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementMatcher{Keys: []string{"b"}, Values: []string{"c"}},
			},
		},
		{
			name: "complex-filter-path",
			path: "a[b.c=d]",
			expectedPathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeIndex, Value: "0"},
			},
			nodeBytes: []byte(`
a:
  - b:
      c: d
  - b:
      c: e
`),
			expectedFilters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementIndexer{Index: 0},
			},
		},
		{
			name: "composite-multi-filter-path",
			path: "a[b.c=d,e=f]",
			expectedPathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeIndex, Value: "1"},
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
			expectedFilters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementIndexer{Index: 1},
			},
		},
		{
			name: "multi-filter-path",
			path: "a[b=d,e=f]",
			expectedPathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeIndex, Value: "1"},
			},
			nodeBytes: []byte(`
a:
  - b: d
    e: c
  - b: d
    e: f
`),
			expectedFilters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementIndexer{Index: 1},
			},
		},
		{
			name: "encapsulated-key",
			path: "a[\"b.c=d\"]",
			expectedPathParts: []transform.PathPart{
				{Type: transform.PartTypeMap, Value: "a"},
				{Type: transform.PartTypeScalar, Value: "b.c=d"},
			},
			expectedFilters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.PathGetter{Path: []string{"b.c=d"}},
			},
		},
		{
			name: "encapsulated-filter-path",
			path: `a["b.c"=d]`,
			expectedPathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeSelector, Value: `"b.c"=d`},
			},
			expectedFilters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementMatcher{Keys: []string{"b.c"}, Values: []string{"d"}},
			},
		},
		{
			name: "root-path",
			expectedPathParts: []transform.PathPart{
				{Type: transform.PartTypeScalar, Value: ""},
			},
			expectedFilters: []yaml.Filter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.nodeBytes, tt.path, tt.expectedPathParts, tt.expectedFilters)
		})
	}
}
