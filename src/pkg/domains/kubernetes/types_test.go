package kube_test

import (
	"context"
	"testing"

	kube "github.com/defenseunicorns/lula/src/pkg/domains/kubernetes"
)

func TestCreateKubernetesDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		spec        *kube.KubernetesSpec
		expectedErr bool
	}{
		{
			name:        "nil spec",
			spec:        nil,
			expectedErr: true,
		},
		{
			name:        "empty spec",
			spec:        &kube.KubernetesSpec{},
			expectedErr: true,
		},
		{
			name: "empty resources",
			spec: &kube.KubernetesSpec{
				Resources: []kube.Resource{},
			},
			expectedErr: false, // currently this is not an error to allow space for PLACEHOLDER validations (TODO: implement empty resources error)
		},
		{
			name: "valid resources",
			spec: &kube.KubernetesSpec{
				Resources: []kube.Resource{
					{
						Name: "test",
						ResourceRule: &kube.ResourceRule{
							Version:  "test",
							Resource: "test",
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "valid resources with field",
			spec: &kube.KubernetesSpec{
				Resources: []kube.Resource{
					{
						Name: "test",
						ResourceRule: &kube.ResourceRule{
							Name:     "test",
							Version:  "test",
							Resource: "test",
							Field: &kube.Field{
								Jsonpath: "test",
							},
						},
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "invalid resources - no name",
			spec: &kube.KubernetesSpec{
				Resources: []kube.Resource{
					{
						ResourceRule: &kube.ResourceRule{
							Version:  "test",
							Resource: "test",
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "invalid resources - no resource-rule",
			spec: &kube.KubernetesSpec{
				Resources: []kube.Resource{
					{
						Name: "test",
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "invalid resource-rule, no version",
			spec: &kube.KubernetesSpec{
				Resources: []kube.Resource{
					{
						Name: "test",
						ResourceRule: &kube.ResourceRule{
							Resource: "test",
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "invalid resource-rule, no resource",
			spec: &kube.KubernetesSpec{
				Resources: []kube.Resource{
					{
						Name: "test",
						ResourceRule: &kube.ResourceRule{
							Version: "test",
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "invalid resource-rule, one name, many namespaces",
			spec: &kube.KubernetesSpec{
				Resources: []kube.Resource{
					{
						Name: "test",
						ResourceRule: &kube.ResourceRule{
							Name:       "test",
							Version:    "test",
							Resource:   "test",
							Namespaces: []string{"test-1", "test-2"},
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "invalid resource-rule, field without name",
			spec: &kube.KubernetesSpec{
				Resources: []kube.Resource{
					{
						Name: "test",
						ResourceRule: &kube.ResourceRule{
							Version:    "test",
							Resource:   "test",
							Namespaces: []string{"test-1", "test-2"},
							Field: &kube.Field{
								Jsonpath: "test",
							},
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "empty create-resources",
			spec: &kube.KubernetesSpec{
				CreateResources: []kube.CreateResource{},
			},
			expectedErr: false, // currently this is not an error to allow space for PLACEHOLDER validations (TODO: implement empty resources error)
		},
		{
			name: "valid create-resources with manifest",
			spec: &kube.KubernetesSpec{
				CreateResources: []kube.CreateResource{
					{
						Name: "test",
						Manifest: `
apiVersion: v1
kind: Pod
metadata:
  name: test
`,
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "valid create-resources with file",
			spec: &kube.KubernetesSpec{
				CreateResources: []kube.CreateResource{
					{
						Name: "test",
						File: "../file/path.yaml",
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "invalid create-resource, missing file and manifest",
			spec: &kube.KubernetesSpec{
				CreateResources: []kube.CreateResource{
					{
						Name: "test",
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "invalid create-resource, missing name",
			spec: &kube.KubernetesSpec{
				CreateResources: []kube.CreateResource{
					{
						File: "../file/path.yaml",
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "invalid create-resources, both file and manifest",
			spec: &kube.KubernetesSpec{
				CreateResources: []kube.CreateResource{
					{
						Name: "test",
						File: "../file/path.yaml",
						Manifest: `
apiVersion: v1
kind: Pod
metadata:
  name: test
`,
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "valid wait",
			spec: &kube.KubernetesSpec{
				Wait: &kube.Wait{
					Condition: "Ready",
					Kind:      "namespace/test",
				},
			},
			expectedErr: false,
		},
		{
			name: "invalid wait, both condition and jsonpath specified",
			spec: &kube.KubernetesSpec{
				Wait: &kube.Wait{
					Condition: "Ready",
					Jsonpath:  "test",
					Kind:      "namespace/test",
				},
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := kube.CreateKubernetesDomain(context.Background(), tt.spec)
			if (err != nil) != tt.expectedErr {
				t.Errorf("CreateKubernetesDomain() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}
