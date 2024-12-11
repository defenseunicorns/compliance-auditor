package opa_test

import (
	"context"
	"errors"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
)

func TestOpaModules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		spec        *opa.OpaSpec
		wantErr     error
		wantPassing int
	}{
		{
			name: "no modules",
			spec: &opa.OpaSpec{
				Rego: "package validate\n\ndefault validate = true",
			},
			wantPassing: 1,
		},
		{
			name: "empty module",
			spec: &opa.OpaSpec{
				Rego:    "package validate\n\ndefault validate = true",
				Modules: map[string]string{"empty": "testdata/empty.rego"},
			},
			wantErr: opa.ErrCompileRego,
		},
		{
			name: "invalid module",
			spec: &opa.OpaSpec{
				Rego:    "package validate\n\nvalidate = true if true",
				Modules: map[string]string{"empty": "testdata/empty.rego"},
			},
			wantErr: opa.ErrCompileRego,
		},
		{
			name: "invalid module path",
			spec: &opa.OpaSpec{
				Rego:    "package validate\n\ndefault validate = false",
				Modules: map[string]string{"test.module": "invalid-path"},
			},
			wantErr: opa.ErrDownloadModule,
		},
		{
			name: "reserved module validation",
			spec: &opa.OpaSpec{
				Rego:    "package validate\n\nimport data.lula.labels as lula_labels\n\nvalidate { lula_labels.has_lula_label(input.pod) }",
				Modules: map[string]string{"lula.labels": "testdata/lula.rego", "validate.rego": "not-used"},
			},
			wantErr: opa.ErrReservedModuleName,
		},
		{
			name: "module validation",
			spec: &opa.OpaSpec{
				Rego:    "package validate\n\nimport data.lula.labels as lula_labels\n\nvalidate { lula_labels.has_lula_label(input.pod) }",
				Modules: map[string]string{"lula.labels": "testdata/lula.rego"},
			},
			wantPassing: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			provider, err := opa.CreateOpaProvider(ctx, tt.spec)
			if err != nil {
				t.Errorf("CreateOpaProvider() error: %v", err)
			}

			result, err := provider.Evaluate(ctx, dummyPod)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if result.Passing != tt.wantPassing {
				t.Errorf("Passing = %d, want %d", result.Passing, tt.wantPassing)
			}
		})
	}
}

var dummyPod = map[string]interface{}{
	"pod": map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": map[string]string{
				"lula": "true",
			},
		},
	},
}
