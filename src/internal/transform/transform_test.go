package transform_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	goyaml "gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/defenseunicorns/lula/src/internal/transform"
)

func createRNode(t *testing.T, data []byte) *yaml.RNode {
	t.Helper()

	node, err := yaml.FromMap(convertBytesToMap(t, data))
	require.NoError(t, err)

	return node
}

// convertBytesToMap converts a byte slice to a map[string]interface{}
func convertBytesToMap(t *testing.T, data []byte) map[string]interface{} {
	var dataMap map[string]interface{}
	err := goyaml.Unmarshal(data, &dataMap)
	require.NoError(t, err)

	return dataMap
}

// TestAdd tests the Add function
func TestAdd(t *testing.T) {
	runTest := func(t *testing.T, current []byte, new []byte, expected []byte) {
		t.Helper()

		node := createRNode(t, current)
		newNode := createRNode(t, new)

		err := transform.Add(node, newNode)
		require.NoError(t, err)

		var nodeMap map[string]interface{}
		err = node.YNode().Decode(&nodeMap)
		require.NoError(t, err)

		require.Equal(t, convertBytesToMap(t, expected), nodeMap)
	}

	tests := []struct {
		name     string
		current  []byte
		new      []byte
		expected []byte
	}{
		{
			name: "test-add-new-key-value",
			current: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
			new: []byte(`
k4: v5
`),

			expected: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
k4: v5
`),
		},
		{
			name: "test-add-existing-key-value",
			current: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
			new: []byte(`
k2: v5
`),

			expected: []byte(`
k1: v1
k2: v5
k3:
  - v3
  - v4
`),
		},
		{
			name: "test-add-list-entry",
			current: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
			new: []byte(`
k3:
  - v5
`),

			expected: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
  - v5
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.current, tt.new, tt.expected)
		})
	}
}

// TestUpdate tests the Update function
func TestUpdate(t *testing.T) {
	runTest := func(t *testing.T, current []byte, new []byte, expected []byte) {
		t.Helper()

		node := createRNode(t, current)
		newNode := createRNode(t, new)

		node, err := transform.Update(node, newNode)
		require.NoError(t, err)

		var nodeMap map[string]interface{}
		err = node.YNode().Decode(&nodeMap)
		require.NoError(t, err)

		require.Equal(t, convertBytesToMap(t, expected), nodeMap)
	}

	tests := []struct {
		name     string
		current  []byte
		new      []byte
		expected []byte
	}{
		{
			name: "test-update-new-key-value",
			current: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
			new: []byte(`
k4: v5
`),

			expected: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
k4: v5
`),
		},
		{
			name: "test-update-existing-key-value",
			current: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
			new: []byte(`
k2: v5
`),

			expected: []byte(`
k1: v1
k2: v5
k3:
  - v3
  - v4
`),
		},
		{
			name: "test-update-list-entry",
			current: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
			new: []byte(`
k3:
  - v5
`),

			expected: []byte(`
k1: v1
k2: v2
k3:
  - v5
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.current, tt.new, tt.expected)
		})
	}
}

// TestDelete tests the Delete function
func TestDelete(t *testing.T) {
	runTest := func(t *testing.T, lastSegment string, current []byte, expected []byte) {
		t.Helper()

		node := createRNode(t, current)
		filters, err := transform.BuildFilters(node, []string{})
		require.NoError(t, err)

		err = transform.Delete(node, lastSegment, filters)
		require.NoError(t, err)

		var nodeMap map[string]interface{}
		err = node.YNode().Decode(&nodeMap)
		require.NoError(t, err)

		require.Equal(t, convertBytesToMap(t, expected), nodeMap)
	}

	tests := []struct {
		name        string
		lastSegment string
		current     []byte
		expected    []byte
	}{
		{
			name:        "test-delete-key-value",
			lastSegment: "k2",
			current: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
			expected: []byte(`
k1: v1
k3:
  - v3
  - v4
`),
		},
		{
			name:        "test-delete-list-key",
			lastSegment: "k3",
			current: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
			expected: []byte(`
k1: v1
k2: v2
`),
		},
		{
			name:        "test-delete-non-existent-key",
			lastSegment: "k4",
			current: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
			expected: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.lastSegment, tt.current, tt.expected)
		})
	}
}

// TestSetNodeAtPath tests the SetNodeAtPath function
func TestSetNodeAtPath(t *testing.T) {
	runTest := func(t *testing.T, pathSlice []string, nodeBytes, newNodeBytes, expected []byte) {
		t.Helper()

		node := createRNode(t, nodeBytes)
		newNode := createRNode(t, newNodeBytes)
		filters, err := transform.BuildFilters(node, pathSlice)
		require.NoError(t, err)

		err = transform.SetNodeAtPath(node, newNode, filters, pathSlice)
		require.NoError(t, err)

		var nodeMap map[string]interface{}
		err = node.YNode().Decode(&nodeMap)
		require.NoError(t, err)

		require.Equal(t, convertBytesToMap(t, expected), nodeMap)
	}

	tests := []struct {
		name      string
		pathSlice []string
		node      []byte
		newNode   []byte
		expected  []byte
	}{
		{
			name:      "simple-path",
			pathSlice: []string{"a", "b"},
			node: []byte(`
a:
  b:
    c: z
  d: y
e: 
  f: g
`),
			newNode: []byte(`
c: x
`),
			expected: []byte(`
a:
  b:
    c: x
  d: y
e: 
  f: g
`),
		},
		{
			name:      "path-with-index-filter",
			pathSlice: []string{"a", "[b=y]"},
			node: []byte(`
a:
  - b: z
    c: 1
  - b: y
    c: 2
`),
			newNode: []byte(`
b: y
c: 3
`),
			expected: []byte(`
a:
  - b: z
    c: 1
  - b: y
    c: 3
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.pathSlice, tt.node, tt.newNode, tt.expected)
		})
	}
}

// TestIntegrationCreateAndExecuteTransform tests the integration of creation and execution of transforms
func TestIntegrationCreateAndExecuteTransform(t *testing.T) {
	runTest := func(t *testing.T, path string, value string, changeType transform.ChangeType, root, valueMap, expected map[string]interface{}) {
		t.Helper()

		tt, err := transform.CreateTransformTarget(root)
		require.NoError(t, err)

		// Execute the transform
		result, err := tt.ExecuteTransform(path, changeType, value, valueMap)
		require.NoError(t, err)
		require.Equal(t, expected, result)
	}

	tests := []struct {
		name       string
		path       string
		changeType transform.ChangeType
		value      string
		valueByte  []byte
		target     []byte
		expected   []byte
	}{
		{
			name:       "update-struct-simple-path",
			path:       "metadata",
			changeType: transform.ChangeTypeUpdate,
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
			valueByte: []byte(`
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
    - item2
    - item3
`),
		},
		{
			name:       "update-data-simple-path",
			path:       "metadata.test",
			changeType: transform.ChangeTypeUpdate,
			target: []byte(`
name: target
some-information: some-data
metadata: {}
`),
			valueByte: []byte(`
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
		{
			name:       "update-at-index-string",
			path:       "foo.subset.[uuid=123].test",
			changeType: transform.ChangeTypeUpdate,
			target: []byte(`
foo:
  subset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: some data to be replaced
`),
			value: "just a string to inject",
			expected: []byte(`
foo:
  subset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: just a string to inject
`),
		},
		{
			name:       "update-at-index-string-with-encapsulation",
			path:       "foo.subset.[\"complex.key\"]",
			changeType: transform.ChangeTypeUpdate,
			target: []byte(`
foo:
  subset:
    complex.key: change-me
`),
			value: "new-value",
			expected: []byte(`
foo:
  subset:
    complex.key: new-value
`),
		},
		{
			name:       "update-at-double-index-map",
			path:       "foo.subset.[uuid=xyz].subsubset.[uuid=123]",
			changeType: transform.ChangeTypeUpdate,
			target: []byte(`
foo:
  subset:
  - uuid: abc
    subsubset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: just some data at 123
  - uuid: xyz
    subsubset:
      - uuid: 321
        test: more data
      - uuid: 123
        test: some data to be replaced
`),
			valueByte: []byte(`
test: just a string to inject
another-key: another-value
`),
			expected: []byte(`
foo:
  subset:
  - uuid: abc
    subsubset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: just some data at 123
  - uuid: xyz
    subsubset:
      - uuid: 321
        test: more data
      - uuid: 123
        test: just a string to inject
        another-key: another-value
`),
		},
		{
			name:       "update-list",
			path:       "foo.subset.[uuid=xyz]",
			changeType: transform.ChangeTypeUpdate,
			target: []byte(`
foo:
  subset:
  - uuid: abc
    subsubset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: just some data at 123
  - uuid: xyz
    subsubset:
      - uuid: 321
        test: more data
      - uuid: 123
        test: some data to be replaced
`),
			valueByte: []byte(`
subsubset:
- uuid: new-uuid
  test: new test data
`),
			expected: []byte(`
foo:
  subset:
  - uuid: abc
    subsubset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: just some data at 123
  - uuid: xyz
    subsubset:
      - uuid: new-uuid
        test: new test data
`),
		},
		{
			name:       "update-list-at-root",
			path:       ".",
			changeType: transform.ChangeTypeUpdate,
			target: []byte(`
foo:
  - uuid: abc
    subsubset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: just some data at 123
  - uuid: xyz
    subsubset:
      - uuid: 321
        test: more data
      - uuid: 123
        test: some data to be replaced
`),
			valueByte: []byte(`
foo:
- uuid: hi
  test: hi
`),
			expected: []byte(`
foo:
- uuid: hi
  test: hi
`),
		},
		{
			name:       "update-at-composite-filter",
			path:       "pods.[metadata.namespace=foo,metadata.name=bar].metadata.labels.app",
			changeType: transform.ChangeTypeUpdate,
			target: []byte(`
pods:
  - metadata:
      name: bar
      namespace: foo
      labels:
        app: replace-me
  - metadata:
      name: baz
      namespace: foo
    labels:
        app: dont-replace-me
`),
			value: "new-app",
			expected: []byte(`
pods:
  - metadata:
      name: bar
      namespace: foo
      labels:
        app: new-app
  - metadata:
      name: baz
      namespace: foo
    labels:
        app: dont-replace-me
`),
		},
		{
			name:       "update-at-composite-double-filter",
			path:       "pods.[metadata.namespace=foo,metadata.name=bar].spec.containers.[name=istio-proxy]",
			changeType: transform.ChangeTypeUpdate,
			target: []byte(`
pods:
  - metadata:
      name: bar
      namespace: foo
      labels:
        app: my-foo-app
    spec:
      containers:
        - name: istio-proxy
          image: replace-me
        - name: foo-app
          image: foo-app:v1
  - metadata:
      name: baz
      namespace: foo
      labels:
        app: my-foo-app
    spec:
      containers:
        - name: istio-proxy
          image: proxyv2
        - name: foo-app
          image: foo-app:v1
`),
			valueByte: []byte(`
image: new-image
`),
			expected: []byte(`
pods:
  - metadata:
      name: bar
      namespace: foo
      labels:
        app: my-foo-app
    spec:
      containers:
        - name: istio-proxy
          image: new-image
        - name: foo-app
          image: foo-app:v1
  - metadata:
      name: baz
      namespace: foo
      labels:
        app: my-foo-app
    spec:
      containers:
        - name: istio-proxy
          image: proxyv2
        - name: foo-app
          image: foo-app:v1
`),
		},
		{
			name:       "add-to-list",
			path:       "foo.subset.[uuid=xyz]",
			changeType: transform.ChangeTypeAdd,
			target: []byte(`
foo:
  subset:
  - uuid: abc
    subsubset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: just some data at 123
  - uuid: xyz
    subsubset:
      - uuid: 321
        test: more data
      - uuid: 123
        test: some data to be replaced
`),
			valueByte: []byte(`
subsubset:
  - uuid: new-uuid
    test: new test data
`),
			expected: []byte(`
foo:
  subset:
  - uuid: abc
    subsubset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: just some data at 123
  - uuid: xyz
    subsubset:
      - uuid: 321
        test: more data
      - uuid: 123
        test: some data to be replaced
      - uuid: new-uuid
        test: new test data
`),
		},
		{
			name:       "delete-from-struct",
			path:       "metadata.some-submap.sub-data",
			changeType: transform.ChangeTypeDelete,
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
			expected: []byte(`
name: target
metadata:
  some-data: target-data
  only-target-field: data
  some-submap:
    only-target-field: target-data
  some-list:
    - item1
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.path, tt.value, tt.changeType, convertBytesToMap(t, tt.target), convertBytesToMap(t, tt.valueByte), convertBytesToMap(t, tt.expected))
		})
	}
}
