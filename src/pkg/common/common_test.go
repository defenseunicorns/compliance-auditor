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
	kjson "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
	"sigs.k8s.io/yaml"
)

const multiValidationPath = "../../test/e2e/scenarios/remote-validations/multi-validations.yaml"
const singleValidationPath = "../../test/e2e/scenarios/remote-validations/validation.opa.yaml"
const whitespaceValidationPath = "../../test/e2e/scenarios/remote-validations/validation.whitespace.yaml"

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
	t.Parallel()

	tests := []struct {
		name           string
		domain         common.Domain
		expectedErr    bool
		expectedDomain string
	}{
		{
			name: "valid kubernetes domain",
			domain: common.Domain{
				Type: "kubernetes",
				KubernetesSpec: &kube.KubernetesSpec{
					Resources: []kube.Resource{
						{
							Name: "podsvt",
							ResourceRule: &kube.ResourceRule{
								Version:    "v1",
								Resource:   "pods",
								Namespaces: []string{"validation-test"},
							},
						},
					},
				},
			},
			expectedErr:    false,
			expectedDomain: "kube.KubernetesDomain",
		},
		{
			name: "invalid kubernetes domain",
			domain: common.Domain{
				Type: "kubernetes",
				KubernetesSpec: &kube.KubernetesSpec{
					Resources: []kube.Resource{
						{
							Name: "podsvt",
							ResourceRule: &kube.ResourceRule{
								Version:    "v1",
								Namespaces: []string{"validation-test"},
							},
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "valid api domain",
			domain: common.Domain{
				Type: "api",
				ApiSpec: &api.ApiSpec{
					Requests: []api.Request{
						{
							Name: "local",
							URL:  "http://localhost",
						},
					},
				},
			},
			expectedErr:    false,
			expectedDomain: "api.ApiDomain",
		},
		{
			name: "invalid api domain",
			domain: common.Domain{
				Type: "api",
				ApiSpec: &api.ApiSpec{
					Requests: []api.Request{
						{
							Name: "local",
						},
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "invalid type domain",
			domain: common.Domain{
				Type: "foo",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := common.GetDomain(&tt.domain)
			if (err != nil) != tt.expectedErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			switch tt.expectedDomain {
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
	t.Parallel()

	tests := []struct {
		name             string
		provider         common.Provider
		expectedErr      bool
		expectedProvider string
	}{
		{
			name: "valid opa provider",
			provider: common.Provider{
				Type: "opa",
				OpaSpec: &opa.OpaSpec{
					Rego: "package validate\n\ndefault validate = false",
				},
			},
			expectedErr:      false,
			expectedProvider: "opa.OpaProvider",
		},
		{
			name: "invalid opa provider",
			provider: common.Provider{
				Type:    "opa",
				OpaSpec: &opa.OpaSpec{},
			},
			expectedErr: true,
		},
		{
			name: "valid kyverno provider",
			provider: common.Provider{
				Type: "kyverno",
				KyvernoSpec: &kyverno.KyvernoSpec{
					Policy: &kjson.ValidatingPolicy{
						Spec: kjson.ValidatingPolicySpec{},
					},
				},
			},
			expectedErr:      false,
			expectedProvider: "kyverno.KyvernoProvider",
		},
		{
			name: "invalid kyverno provider",
			provider: common.Provider{
				Type:        "kyverno",
				KyvernoSpec: &kyverno.KyvernoSpec{},
			},
			expectedErr: true,
		},
		{
			name: "invalid type provider",
			provider: common.Provider{
				Type: "foo",
			},
			expectedErr: true,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := common.GetProvider(&tt.provider, ctx)
			if (err != nil) != tt.expectedErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			switch tt.expectedProvider {
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
		uuid    string
		wantErr bool
	}{
		{
			name:    "Valid Validation string",
			data:    validationStrings[0],
			uuid:    "88AB3470-B96B-4D7C-BC36-02BF9563C46C",
			wantErr: false,
		},
		{
			name:    "Invalid Validation, successfully unmarshalled",
			data:    "Test: test",
			uuid:    "a50c374a-deee-4032-9a0e-38e624f49c3d", // check that still returns a valid UUID even if invalid validation string
			wantErr: true,
		},
		{
			name:    "Empty Data",
			data:    "",
			uuid:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lulaValidation, err := common.ValidationFromString(tt.data, tt.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidationFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if lulaValidation.UUID != tt.uuid {
				t.Errorf("ValidationFromString() UUID = %v, want %v", lulaValidation.UUID, tt.uuid)
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
				UUID: "1f639c6b-4e86-4c66-88b2-22dbf6d7ac02",
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

		if resource.UUID != validation.Metadata.UUID {
			t.Errorf("ToResource() resource UUID %s should match created validation UUID %s", resource.UUID, validation.Metadata.UUID)
		}
	})

	t.Run("It trims whitespace from the validation", func(t *testing.T) {
		t.Parallel()
		validationBytes := loadTestData(t, whitespaceValidationPath)

		validation, err := common.ReadValidationsFromYaml(validationBytes)
		if err != nil {
			t.Fatalf("yaml.Unmarshal failed: %v", err)
		}

		if len(validation) > 1 {
			t.Errorf("Expected 1 validation, got %d", len(validation))
		}

		resource, err := validation[0].ToResource()
		if err != nil {
			t.Errorf("ToResource() error = %v", err)
		}
		strings.Contains(resource.Description, " \n")
		if strings.Contains(resource.Description, " \n") {
			t.Errorf("ToResource() description = should not contain whitespace followed by newline")
		}
	})
}

func TestReadValidationsFromYaml(t *testing.T) {
	multipleValidations := loadTestData(t, multiValidationPath)
	singleValidations := loadTestData(t, singleValidationPath)

	tests := []struct {
		name          string
		validations   []byte
		expectedCount int
	}{
		{
			name:          "multiple validations",
			validations:   multipleValidations,
			expectedCount: 2,
		},
		{
			name:          "single validation",
			validations:   singleValidations,
			expectedCount: 1,
		},
		{
			name:          "no validations",
			validations:   []byte{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validations, err := common.ReadValidationsFromYaml(tt.validations)
			if err != nil {
				t.Errorf("Error reading validations from yaml: %v", err)
			}
			if len(validations) != tt.expectedCount {
				t.Errorf("Expected %d validations, but got %d", tt.expectedCount, len(validations))
			}
		})
	}
}

func TestIsVersionValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		versionConstraint string
		version           string
		expectedValid     bool
		expectedErr       bool
	}{
		{
			name:              "Valid constraint and version",
			versionConstraint: ">= 1.0.0, < 2.0.0",
			version:           "1.5.0",
			expectedValid:     true,
			expectedErr:       false,
		},
		{
			name:              "Valid constraint and version (exact match)",
			versionConstraint: "1.0.0",
			version:           "1.0.0",
			expectedValid:     true,
			expectedErr:       false,
		},
		{
			name:              "Invalid version",
			versionConstraint: ">= 1.0.0, < 2.0.0",
			version:           "2.5.0",
			expectedValid:     false,
			expectedErr:       false,
		},
		{
			name:              "Invalid constraint syntax",
			versionConstraint: ">=> 1.0.0",
			version:           "1.5.0",
			expectedValid:     false,
			expectedErr:       true,
		},
		{
			name:              "Invalid version syntax",
			versionConstraint: ">= 1.0.0, < 2.0.0",
			version:           "invalid",
			expectedValid:     false,
			expectedErr:       true,
		},
		{
			name:              "Unset version",
			versionConstraint: ">= 1.0.0, < 2.0.0",
			version:           "unset",
			expectedValid:     true,
			expectedErr:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := common.IsVersionValid(tc.versionConstraint, tc.version)
			if (err != nil) != tc.expectedErr {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, err)
			}
			if valid != tc.expectedValid {
				t.Fatalf("expected valid: %v, got: %v", tc.expectedValid, valid)
			}
		})
	}
}
