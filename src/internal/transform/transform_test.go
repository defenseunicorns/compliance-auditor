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
	runTest := func(t *testing.T, filters []yaml.Filter, pathParts []transform.PathPart, operandIdx int, current []byte, expected []byte) {
		t.Helper()

		node := createRNode(t, current)

		err := transform.Delete(node, filters, pathParts, operandIdx)
		require.NoError(t, err)

		var nodeMap map[string]interface{}
		err = node.YNode().Decode(&nodeMap)
		require.NoError(t, err)

		require.Equal(t, convertBytesToMap(t, expected), nodeMap)
	}

	tests := []struct {
		name       string
		filters    []yaml.Filter
		pathParts  []transform.PathPart
		operandIdx int
		current    []byte
		expected   []byte
	}{
		{
			name: "delete-root-key-value",
			filters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"k1"}},
			},
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeScalar, Value: "k1"},
			},
			operandIdx: 0,
			current: []byte(`
k1: v1
k2: v2
k3:
  - v3
  - v4
`),
			expected: []byte(`
k2: v2
k3:
  - v3
  - v4
`),
		},
		{
			name: "delete-sub-key",
			filters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"k1"}},
				yaml.PathGetter{Path: []string{"k2"}},
			},
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeMap, Value: "k1"},
				{Type: transform.PartTypeScalar, Value: "k2"},
			},
			operandIdx: 1,
			current: []byte(`
k1:
  k2: v2
k3:
  - v3
  - v4
`),
			expected: []byte(`
k1: {}
k3:
  - v3
  - v4
`),
		},
		{
			name: "delete-non-existent-key",
			filters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"k4"}},
			},
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeScalar, Value: "k4"},
			},
			operandIdx: 0,
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
		{
			name: "delete-indexed-key",
			filters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"k3"}},
				yaml.PathGetter{Path: []string{"0"}},
			},
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "k3"},
				{Type: transform.PartTypeIndex, Value: "0"},
			},
			operandIdx: -1,
			current: []byte(`
k1: v1
k2: v2
k3:
  - k31: v3
  - k32: v4
`),
			expected: []byte(`
k1: v1
k2: v2
k3:
  - k32: v4
`),
		},
		{
			name: "delete-indexed-key-for-string",
			filters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"k3"}},
				yaml.PathGetter{Path: []string{"0"}},
			},
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "k3"},
				{Type: transform.PartTypeIndex, Value: "0"},
			},
			operandIdx: -1,
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
  - v4
`),
		},
		{
			name: "delete-selected-key",
			filters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"k3"}},
				yaml.ElementMatcher{Keys: []string{"k32"}, Values: []string{"v4"}},
			},
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "k3"},
				{Type: transform.PartTypeSelector, Value: "k32=v4"},
			},
			operandIdx: -1,
			current: []byte(`
k1: v1
k2: v2
k3:
  - k31: v3
  - k32: v4
`),
			expected: []byte(`
k1: v1
k2: v2
k3:
  - k31: v3
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.filters, tt.pathParts, tt.operandIdx, tt.current, tt.expected)
		})
	}
}

// TestSetNodeAtPath tests the SetNodeAtPath function
func TestSetNodeAtPath(t *testing.T) {
	runTest := func(t *testing.T, pathParts []transform.PathPart, filters []yaml.Filter, lastItemIdx int, nodeBytes, newNodeBytes, expected []byte) {
		t.Helper()

		node := createRNode(t, nodeBytes)
		newNode := createRNode(t, newNodeBytes)

		err := transform.SetNodeAtPath(node, newNode, filters, pathParts, lastItemIdx)
		require.NoError(t, err)

		var nodeMap map[string]interface{}
		err = node.YNode().Decode(&nodeMap)
		require.NoError(t, err)

		require.Equal(t, convertBytesToMap(t, expected), nodeMap)
	}

	tests := []struct {
		name         string
		pathParts    []transform.PathPart
		filters      []yaml.Filter
		finalItemIdx int
		node         []byte
		newNode      []byte
		expected     []byte
		expectErr    bool
		errContains  string
	}{
		{
			name: "simple-path:a.b",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeMap, Value: "a"},
				{Type: transform.PartTypeScalar, Value: "b"},
			},
			filters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.PathGetter{Path: []string{"b"}},
			},
			finalItemIdx: 1,
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
			name: "path-with-filter:a[b=y]",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeSelector, Value: "b=y"},
			},
			filters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementMatcher{Keys: []string{"b"}, Values: []string{"y"}},
			},
			finalItemIdx: -1,
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
		{
			name: "path-with-index-filter:a[0]",
			pathParts: []transform.PathPart{
				{Type: transform.PartTypeSequence, Value: "a"},
				{Type: transform.PartTypeIndex, Value: "0"},
			},
			filters: []yaml.Filter{
				yaml.PathGetter{Path: []string{"a"}},
				yaml.ElementIndexer{Index: 0},
			},
			finalItemIdx: -1,
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
  - b: y
    c: 3
  - b: y
    c: 2
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTest(t, tt.pathParts, tt.filters, tt.finalItemIdx, tt.node, tt.newNode, tt.expected)
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
			path:       `foo.subset["complex.key"]`,
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
			name:       "update-at-int-index",
			path:       "foo.subset.[0].test",
			changeType: transform.ChangeTypeUpdate,
			target: []byte(`
foo:
  subset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: some more data
`),
			value: "just a string to inject",
			expected: []byte(`
foo:
  subset:
    - uuid: 321
      test: just a string to inject
    - uuid: 123
      test: some more data
`),
		},
		{
			name:       "update-at-int-index-last",
			path:       "foo.subset.[0]",
			changeType: transform.ChangeTypeUpdate,
			target: []byte(`
foo:
  subset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: some more data
`),
			valueByte: []byte(`
uuid: new-uuid
test: just a string to inject
`),
			expected: []byte(`
foo:
  subset:
    - uuid: new-uuid
      test: just a string to inject
    - uuid: 123
      test: some more data
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
