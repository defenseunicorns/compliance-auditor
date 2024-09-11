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
			// TODO: Should we extend the functionaly to allow for non-existent paths?
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
`),
		},
		{
			name: "test-inject-at-index",
			path: "foo.subset[uuid=123]",
			target: []byte(`
foo:
  subset:
    - uuid: 321
      test: some data
    - uuid: 123
      test: some data to be replaced
`),
			subset: []byte(`
test: just a string to inject
`),
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
			name: "test-inject-at-double-index",
			path: "foo.subset[uuid=xyz].subsubset[uuid=123]",
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
			subset: []byte(`
test: just a string to inject
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
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := inject.InjectMapData(convertBytesToMap(t, tt.target), convertBytesToMap(t, tt.subset), tt.path)
			if err != nil {
				t.Errorf("InjectMapData() error = %v", err)
			}
			assert.Equal(t, convertBytesToMap(t, tt.expected), result, "The maps should be equal")
		})
	}
}

// TestInjectMapData2 tests the InjectMapData function -> more complex filtering
// actually do I need this if I could support like:
// pods[metadata.namespace=monitoring,metadata.labels.app=kube-prometheus-stack-operator].spec.containers[name=istio-proxy]
// if I pass nothing in value, will it delete it?
// ^^ this one should be deleted though to support the policy...
func TestInjectMapData2(t *testing.T) {
	// Add a test if multiple nodes meet the criteria
	tests := []struct {
		name     string
		path     string
		target   []byte
		subset   []byte
		expected []byte
	}{
		{
			name: "test-inject-at-double-filter",
			path: "pods.metadata[namespace=foo,name=bar].labels",
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
			subset: []byte(`
app: new-app
`),
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
			name: "test-inject-at-nested-double-filter",
			path: "pods[metadata.namespace=foo,metadata.name=bar].spec.containers[name=istio-proxy]",
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
			subset: []byte(`
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := inject.InjectMapData(convertBytesToMap(t, tt.target), convertBytesToMap(t, tt.subset), tt.path)
			if err != nil {
				t.Errorf("InjectMapData() error = %v", err)
			}
			assert.Equal(t, convertBytesToMap(t, tt.expected), result, "The maps should be equal")
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
