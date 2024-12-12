package opa_test

import (
	"context"
	"errors"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
)

func TestCreateOpaProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		spec    *opa.OpaSpec
		wantErr error
	}{
		{
			name: "valid spec",
			spec: &opa.OpaSpec{
				Rego: "package validate\n\ndefault validate = false",
			},
		},
		{
			name: "valid spec with output",
			spec: &opa.OpaSpec{
				Rego: "package validate\n\ndefault validate = false",
				Output: &opa.OpaOutput{
					Validation: "validate.result",
					Observations: []string{
						"validate.observation",
					},
				},
			},
		},
		{
			name:    "nil spec",
			spec:    nil,
			wantErr: opa.ErrNilSpec,
		},
		{
			name: "empty rego",
			spec: &opa.OpaSpec{
				Rego: "",
			},
			wantErr: opa.ErrEmptyRego,
		},
		{
			name: "invalid validation path",
			spec: &opa.OpaSpec{
				Rego: "package validate\n\ndefault validate = false",
				Output: &opa.OpaOutput{
					Validation: "invalid-path",
				},
			},
			wantErr: opa.ErrInvalidValidationPath,
		},
		{
			name: "invalid observation path",
			spec: &opa.OpaSpec{
				Rego: "package validate\n\ndefault validate = false",
				Output: &opa.OpaOutput{
					Observations: []string{"invalid-path"},
				},
			},
			wantErr: opa.ErrInvalidObservationPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := opa.CreateOpaProvider(context.Background(), tt.spec)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateOpaProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
