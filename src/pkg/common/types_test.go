package common_test

import (
	"errors"
	"testing"

	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/common"
)

func TestToLulaValidation(t *testing.T) {
	t.Parallel()
	config.CLIVersion = "1.0.0" // Set the version for testing purposes

	tests := []struct {
		name            string
		inputYaml       []byte
		expectErr       bool
		expectedErrType error
	}{
		{
			name: "Valid validation",
			inputYaml: []byte(`
lula-version: "1.0.0"
metadata:
  name: "test-valid"
domain:
  type: "kubernetes"
  kubernetes-spec: 
    resources: []
provider:
  type: "opa"
  opa-spec:
    rego: "package validate\n\ndefault validate = false"
`),
			expectErr: false,
		},
		{
			name: "Invalid version",
			inputYaml: []byte(`
lula-version: "2.0.0"
metadata:
  name: "test-invalid-version"
domain:
  type: "kubernetes"
  kubernetes-spec:
    resources: []
provider:
  type: "opa"
  opa-spec:
    rego: "package validate\n\ndefault validate = false"
`),
			expectErr:       true,
			expectedErrType: common.ErrInvalidVersion,
		},
		{
			name: "Invalid schema",
			inputYaml: []byte(`
lula-version: "1.0.0"
metadata: {}
`),
			expectErr:       true,
			expectedErrType: common.ErrInvalidSchema,
		},
		{
			name: "Invalid domain schema, bad type",
			inputYaml: []byte(`
lula-version: "1.0.0"
metadata:
  name: "test-invalid-domain"
domain:
  type: "unknown"
provider:
  type: "opa"
  opa-spec:
    rego: "package validate\n\ndefault validate = false"
`),
			expectErr:       true,
			expectedErrType: common.ErrInvalidSchema,
		},
		{
			name: "Invalid domain schema, missing spec",
			inputYaml: []byte(`
lula-version: "1.0.0"
metadata:
  name: "test-invalid-domain"
domain:
  type: "kubernetes"
provider:
  type: "opa"
  opa-spec:
    rego: "package validate\n\ndefault validate = false"
`),
			expectErr:       true,
			expectedErrType: common.ErrInvalidSchema,
		},
		{
			name: "Invalid provider schema, bad type",
			inputYaml: []byte(`
lula-version: "1.0.0"
metadata:
  name: "test-invalid-provider"
domain:
  type: "kubernetes"
  kubernetes-spec:
    resources: []
provider:
  type: "unknown"
`),
			expectErr:       true,
			expectedErrType: common.ErrInvalidSchema,
		},
		{
			name: "Invalid provider schema, missing spec",
			inputYaml: []byte(`
lula-version: "1.0.0"
metadata:
  name: "test-invalid-provider"
domain:
  type: "kubernetes"
  kubernetes-spec:
    resources: []
provider:
  type: "opa"
`),
			expectErr:       true,
			expectedErrType: common.ErrInvalidSchema,
		},
		{
			name: "Bad kubernetes spec - missing resource",
			inputYaml: []byte(`
lula-version: "1.0.0"
metadata:
  name: "test-invalid-provider"
domain:
  type: "kubernetes"
  kubernetes-spec:
    resources:
      - name: "test"
        resource-rule:
          version: "test"
provider:
  type: "opa"
  opa-spec:
    rego: "package validate\n\ndefault validate = false"
`),
			expectErr:       true,
			expectedErrType: common.ErrInvalidDomain,
		},
		{
			name: "Bad opa spec - bad output validation format",
			inputYaml: []byte(`
lula-version: "1.0.0"
metadata:
  name: "test-invalid-provider"
domain:
  type: "kubernetes"
  kubernetes-spec:
    resources: []
provider:
  type: opa
  opa-spec:
    rego: "package validate\n\ndefault validate = false"
    output:
      validation: validate-result
`),
			expectErr:       true,
			expectedErrType: common.ErrInvalidProvider,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var validation common.Validation
			err := validation.UnmarshalYaml(tt.inputYaml)
			if err != nil {
				t.Fatalf("UnmarshalYaml failed: %v", err)
			}

			_, err = validation.ToLulaValidation()
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if (err != nil) && !errors.Is(err, tt.expectedErrType) {
				t.Fatalf("expected error type: %v, got: %v", tt.expectedErrType, err)
			}
		})
	}
}
