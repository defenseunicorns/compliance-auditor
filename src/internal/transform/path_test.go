package transform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
				{Type: transform.PartTypeFilter, Value: "b=c"},
			},
		},
		{
			name: "composite-filter-path",
			path: "a[b.c=d]",
			expected: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeFilter, Value: "b.c=d"},
			},
		},
		{
			name: "multi-filter-path",
			path: "a[b=d,e=f]",
			expected: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeFilter, Value: "b=d,e=f"},
			},
		},
		{
			name: "encapsulated-filter-path",
			path: `a["b.c"=d]`,
			expected: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeFilter, Value: `"b.c"=d`},
			},
		},
		{
			name: "encapsulated-key",
			path: `a["b.c=d"]`,
			expected: []transform.PathPart{
				{Type: transform.PartTypeMap, Value: "a"},
				{Type: transform.PartTypeScalar, Value: `"b.c=d"`},
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
				{Type: transform.PartTypeFilter, Value: "b=c"},
				{Type: transform.PartTypeSequence, Value: "d"},
				{Type: transform.PartTypeFilter, Value: "e=f"},
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

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		changeType    transform.ChangeType
		expectedIndex int
		wantErr       bool
		errContains   string
	}{
		{
			name:          "root-update",
			path:          "",
			changeType:    transform.ChangeTypeUpdate,
			expectedIndex: 0,
		},
		{
			name:          "root-dot-update",
			path:          ".",
			changeType:    transform.ChangeTypeUpdate,
			expectedIndex: 0,
		},
		{
			name:          "simple-update",
			path:          "a.b",
			changeType:    transform.ChangeTypeUpdate,
			expectedIndex: 1,
		},
		{
			name:          "filter-update",
			path:          "a[b.c=d]",
			changeType:    transform.ChangeTypeUpdate,
			expectedIndex: -1,
		},
		{
			name:          "index-update",
			path:          "a[0]",
			changeType:    transform.ChangeTypeUpdate,
			expectedIndex: -1,
		},
		{
			name:          "simple-delete",
			path:          "a.b",
			changeType:    transform.ChangeTypeDelete,
			expectedIndex: 1,
		},
		{
			name:          "invalid-filter-delete",
			path:          "a[b.c=d]",
			changeType:    transform.ChangeTypeDelete,
			expectedIndex: -1,
			wantErr:       true,
			errContains:   "cannot delete a list entry",
		},
		{
			name:          "invalid-index-delete",
			path:          "a[0]",
			changeType:    transform.ChangeTypeDelete,
			expectedIndex: -1,
			wantErr:       true,
			errContains:   "cannot delete a list entry",
		},
		{
			name:          "invalid-path-delete",
			path:          ".",
			changeType:    transform.ChangeTypeDelete,
			expectedIndex: -1,
			wantErr:       true,
			errContains:   "invalid path, cannot delete root",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathParts, index, err := transform.ResolvePath(tt.path, tt.changeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolvePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				require.Contains(t, err.Error(), tt.errContains)
			}
			assert.Equal(t, tt.expectedIndex, index)
			assert.GreaterOrEqual(t, len(pathParts), 0)
		})
	}
}
