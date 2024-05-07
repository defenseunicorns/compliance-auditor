package common_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
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
				KubernetesSpec: &validKubernetes,
			},
			expected: "kube.KubernetesDomain",
		},
		{
			name: "api domain",
			domain: common.Domain{
				Type:    "api",
				ApiSpec: &validApi,
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
			result := common.GetDomain(&tt.domain, ctx)

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
				OpaSpec: &validOpa,
			},
			expected: "opa.OpaProvider",
		},
		{
			name: "kyverno provider",
			provider: common.Provider{
				Type:        "kyverno",
				KyvernoSpec: &validKyverno,
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
			result := common.GetProvider(&tt.provider, ctx)

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

func TestValidationFromString(t *testing.T) {
	validBackMatterMapBytes := loadTestData(t, "../../test/unit/common/oscal/valid-back-matter-map.yaml")

	var validBackMatterMap map[string]string
	if err := yaml.Unmarshal(validBackMatterMapBytes, &validBackMatterMap); err != nil {
		t.Fatalf("yaml.Unmarshal failed: %v", err)
	}

	validationStrings := make([]string, 0)
	for _, v := range validBackMatterMap {
		validationStrings = append(validationStrings, v)
	}

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "Valid Validation string",
			data:    validationStrings[0],
			wantErr: false,
		},
		{
			name:    "Invalid Validation string",
			data:    "Test: test",
			wantErr: true,
		},
		{
			name:    "Empty Data",
			data:    "",
			wantErr: true,
		},
		// Additional test cases can be added here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := common.ValidationFromString(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidationFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}

func TestSwitchCwd(t *testing.T) {

	tempDir := t.TempDir()

	tests := []struct {
		name     string
		path     string
		expected string
		wantErr  bool
	}{
		{
			name:     "Valid path",
			path:     tempDir,
			expected: tempDir,
			wantErr:  false,
		},
		{
			name:     "Path is File",
			path:     "./common_test.go",
			expected: "./",
			wantErr:  false,
		},
		{
			name:     "Invalid path",
			path:     "/nonexistent",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty Path",
			path:     "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFunc, err := common.SetCwdToFileDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SwitchCwd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				defer resetFunc()
				wd, _ := os.Getwd()
				expected, err := filepath.Abs(tt.expected)
				if err != nil {
					t.Errorf("SwitchCwd() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !strings.HasSuffix(wd, expected) {
					t.Errorf("SwitchCwd() working directory = %v, want %v", wd, tt.expected)
				}
			}
		})
	}
}

func TestValidationToResource(t *testing.T) {
	t.Parallel()
	t.Run("It populates a resource from a validation", func(t *testing.T) {
		t.Parallel()
		validation := &common.Validation{
			Metadata: &common.Metadata{
				UUID: "1234",
				Name: "Test Validation",
			},
			Provider: &common.Provider{
				Type: "test",
			},
			Domain: &common.Domain{
				Type: "test",
			},
		}

		resource, err := validation.ToResource()
		if err != nil {
			t.Errorf("ToResource() error = %v", err)
		}

		if resource.Title != validation.Metadata.Name {
			t.Errorf("ToResource() title = %v, want %v", resource.Title, validation.Metadata.Name)
		}

		if resource.UUID != validation.Metadata.UUID {
			t.Errorf("ToResource() UUID = %v, want %v", resource.UUID, validation.Metadata.UUID)
		}

		if resource.Description == "" {
			t.Errorf("ToResource() description = %v, want %v", resource.Description, "")
		}
	})

	t.Run("It adds a UUID if one does not exist", func(t *testing.T) {
		t.Parallel()
		validation := &common.Validation{
			Metadata: &common.Metadata{
				Name: "Test Validation",
			},
			Provider: &common.Provider{
				Type: "test",
			},
			Domain: &common.Domain{
				Type: "test",
			},
		}

		resource, err := validation.ToResource()
		if err != nil {
			t.Errorf("ToResource() error = %v", err)
		}

		if resource.UUID == validation.Metadata.UUID {
			t.Errorf("ToResource() description = \"\", want a valid UUID")
		}
	})
}
