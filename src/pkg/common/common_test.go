package common_test

import (
	"context"
	"os"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/domains/api"
	kube "github.com/defenseunicorns/lula/src/pkg/domains/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/providers/kyverno"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"sigs.k8s.io/yaml"
)

const validKubernetesPath = "../../test/unit/common/valid-kubernetes-spec.yaml"
const validApiPath = "../../test/unit/common/valid-api-spec.yaml"
const validOpaPath = "../../test/unit/common/valid-opa-spec.yaml"
const validKyvernoPath = "../../test/unit/common/valid-kyverno-spec.yaml"

// Helper function to load test data
func loadTestData(t *testing.T, path string) []byte {
	t.Helper() // Marks this function as a test helper
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file '%s': %v", path, err)
	}
	return data
}

func TestGetDomain(t *testing.T) {
	validKubernetesBytes := loadTestData(t, validKubernetesPath)
	validApiBytes := loadTestData(t, validApiPath)

	var validKubernetes kube.KubernetesSpec
	if err := yaml.Unmarshal(validKubernetesBytes, &validKubernetes); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	var validApi api.ApiSpec
	if err := yaml.Unmarshal(validApiBytes, &validApi); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	// Define test cases
	tests := []struct {
		name     string
		domain   common.Domain
		expected string
	}{
		{
			name: "kubernetes domain",
			domain: common.Domain{
				Type:           "kubernetes",
				KubernetesSpec: validKubernetes,
			},
			expected: "kube.KubernetesDomain",
		},
		{
			name: "api domain",
			domain: common.Domain{
				Type:    "api",
				ApiSpec: validApi,
			},
			expected: "api.ApiDomain",
		},
		{
			name: "unsupported domain",
			domain: common.Domain{
				Type: "unsupported",
			},
			expected: "nil",
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.GetDomain(tt.domain, ctx)

			switch tt.expected {
			case "kube.KubernetesDomain":
				if _, ok := result.(kube.KubernetesDomain); !ok {
					t.Errorf("Expected result to be kube.KubernetesDomain, got %T", result)
				}
			case "api.ApiDomain":
				if _, ok := result.(api.ApiDomain); !ok {
					t.Errorf("Expected result to be api.ApiDomain, got %T", result)
				}
			case "nil":
				if result != nil {
					t.Errorf("Expected result to be nil, got %T", result)
				}
			}
		})
	}
}

func TestGetProvider(t *testing.T) {
	validOpaBytes := loadTestData(t, validOpaPath)
	validKyvernoBytes := loadTestData(t, validKyvernoPath)

	var validOpa opa.OpaSpec
	if err := yaml.Unmarshal(validOpaBytes, &validOpa); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	var validKyverno kyverno.KyvernoSpec
	if err := yaml.Unmarshal(validKyvernoBytes, &validKyverno); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	tests := []struct {
		name     string
		provider common.Provider
		expected string
	}{
		{
			name: "opa provider",
			provider: common.Provider{
				Type:    "opa",
				OpaSpec: validOpa,
			},
			expected: "opa.OpaProvider",
		},
		{
			name: "kyverno provider",
			provider: common.Provider{
				Type:        "kyverno",
				KyvernoSpec: validKyverno,
			},
			expected: "kyverno.KyvernoProvider",
		},
		{
			name: "unsupported provider",
			provider: common.Provider{
				Type: "unsupported",
			},
			expected: "nil",
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.GetProvider(tt.provider, ctx)

			switch tt.expected {
			case "opa.OpaProvider":
				if _, ok := result.(opa.OpaProvider); !ok {
					t.Errorf("Expected result to be opa.OpaProvider, got %T", result)
				}
			case "kyverno.KyvernoProvider":
				if _, ok := result.(kyverno.KyvernoProvider); !ok {
					t.Errorf("Expected result to be kyverno.KyvernoProvider, got %T", result)
				}
			case "nil":
				if result != nil {
					t.Errorf("Expected result to be nil, got %T", result)
				}
			}
		})
	}
}
