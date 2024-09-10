package inject_test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/internal/inject"
	"github.com/stretchr/testify/assert"
	goyaml "gopkg.in/yaml.v3"
)

// TestInjectMapData tests the InjectMapData function
func TestInjectMapData(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		target   []byte
		subset   []byte
		expected []byte
	}{
		{
			name: "test-merge-subset-with-list",
			path: "metadata",
			target: []byte(`
name: target
metadata:
  some-data: target-data
  only-target-field: data
  some-submap:
    only-target-field: target-data
    sub-data: this-should-be-overwritten
  some-list:
    - item1
`),
			subset: []byte(`
some-data: subset-data
some-submap:
  sub-data: my-submap-data
  more-data: some-more-data
some-list:
  - item2
  - item3
`),
			expected: []byte(`
name: target
metadata:
  some-data: subset-data
  only-target-field: data
  some-submap:
    only-target-field: target-data
    sub-data: my-submap-data
    more-data: some-more-data
  some-list:
    - item1
    - item2
    - item3
`),
		},
		{
			name: "test-merge-at-root",
			path: "",
			target: []byte(`
name: target
some-information: some-data
some-map:
  test-key: test-value
`),
			subset: []byte(`
more-information: more-data
some-map:
  test-key: subset-value
`),
			expected: []byte(`
name: target
more-information: more-data
some-information: some-data
some-map:
  test-key: subset-value
`),
		},
		{
			name: "test-merge-at-non-existant-path",
			path: "metadata.test",
			target: []byte(`
name: target
some-information: some-data
`),
			subset: []byte(`
name: some-name
more-metdata: here
`),
			expected: []byte(`
name: target
some-information: some-data
metadata:
  test:
    name: some-name
    more-metdata: here
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := inject.InjectMapData(convertBytesToMap(t, tt.target), convertBytesToMap(t, tt.subset), tt.path)
			if err != nil {
				t.Errorf("InjectMapData() error = %v", err)
			}
			assert.Equal(t, result, convertBytesToMap(t, tt.expected), "The maps should be equal")
		})
	}
}

// convertBytesToMap converts a byte slice to a map[string]interface{}
func convertBytesToMap(t *testing.T, data []byte) map[string]interface{} {
	var dataMap map[string]interface{}
	if err := goyaml.Unmarshal(data, &dataMap); err != nil {
		t.Errorf("yaml.Unmarshal failed: %v", err)
	}
	return dataMap
}
