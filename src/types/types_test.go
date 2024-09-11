package types_test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/types"
)

func TestRunTests(t *testing.T) {
	// TODO: read in a resources.json file, set that as domain resources in a LulaValidation, and run tests
	// actually I think I need to create a validation too, and pull that in to test with...
	tests := []struct {
		name         string
		tests        []types.Test
		resourcePath string
	}{
		{
			name: "test-with-list-target",
			tests: []types.Test{
				{
					Name: "test-1",
					Permutations: []types.Permutation{
						{
							ListTarget: "namespaces[metadata.name=istio-system]",
							Path:       "metadata.labels",
							Value: map[string]interface{}{
								"app.kubernetes.io/managed-by": "lula",
							},
						},
					},
				},
			},
			resourcePath: "../test/unit/types/resources-all-namespaces.yaml",
		},
		{
			name: "test-flat-path",
			tests: []types.Test{
				{
					Name: "test-1",
					Permutations: []types.Permutation{
						{
							Path: "namespaces.metadata[name=istio-system].labels",
							Value: map[string]interface{}{
								"app.kubernetes.io/managed-by": "lula",
							},
						},
					},
				},
			},
			resourcePath: "../test/unit/types/resources-all-namespaces.yaml",
		},
		{
			name: "test-pods",
			tests: []types.Test{
				{
					Name: "test-1",
					Permutations: []types.Permutation{
						{
							Target: "pods[metadata.namespace=something].metadata[namespace=istio-system].labels[blah=blah]"
							Path: "pods.metadata[name=istio-system].labels",
							Value: map[string]interface{}{
								"app.kubernetes.io/managed-by": "lula",
							},
						},
					},
				},
			},
			resourcePath: "../test/unit/types/resources-all-namespaces.yaml",
		},
	}

}
