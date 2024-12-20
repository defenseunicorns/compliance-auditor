package transform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defenseunicorns/lula/src/internal/transform"
)

func TestPathToParts(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected []transform.PathPart
	}{
		{
			name: "simple-path",
			path: "a.b.c",
			expected: []transform.PathPart{
				{Type: transform.PartTypeMap, Value: "a"},
				{Type: transform.PartTypeMap, Value: "b"},
				{Type: transform.PartTypeScalar, Value: "c"},
			},
		},
		{
			name: "filter-path",
			path: "a[b=c]",
			expected: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeSelector, Value: "b=c"},
			},
		},
		{
			name: "composite-filter-path",
			path: "a[b.c=d]",
			expected: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeSelector, Value: "b.c=d"},
			},
		},
		{
			name: "multi-filter-path",
			path: "a[b=d,e=f]",
			expected: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeSelector, Value: "b=d,e=f"},
			},
		},
		{
			name: "encapsulated-filter-path",
			path: `a["b.c"=d]`,
			expected: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeSelector, Value: `"b.c"=d`},
			},
		},
		{
			name: "encapsulated-key",
			path: `a["b.c=d"]`,
			expected: []transform.PathPart{
				{Type: transform.PartTypeMap, Value: "a"},
				{Type: transform.PartTypeScalar, Value: "b.c=d"},
			},
		},
		{
			name: "filter-by-index",
			path: `a[0]`,
			expected: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeIndex, Value: "0"},
			},
		},
		{
			name: "complicated-path-all-types",
			path: `a[b=c].d[e=f].h.j`,
			expected: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeSelector, Value: "b=c"},
				{Type: transform.PartTypeSequence, Value: "d"},
				{Type: transform.PartTypeSelector, Value: "e=f"},
				{Type: transform.PartTypeMap, Value: "h"},
				{Type: transform.PartTypeScalar, Value: "j"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transform.PathToParts(tt.path)
			assert.Equal(t, tt.expected, got)
		})
	}
}
