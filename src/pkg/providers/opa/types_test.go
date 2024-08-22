package opa_test

import (
	"context"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
)

func TestCreateOpaProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		spec    *opa.OpaSpec
		wantErr bool
	}{
		{
			name: "valid spec",
			spec: &opa.OpaSpec{
				Rego: "package validate\n\ndefault validate = false",
			},
			wantErr: false,
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
			wantErr: false,
		},
		{
			name:    "nil spec",
			spec:    nil,
			wantErr: true,
		},
		{
			name: "empty rego",
			spec: &opa.OpaSpec{
				Rego: "",
			},
			wantErr: true,
		},
		{
			name: "invalid validation path",
			spec: &opa.OpaSpec{
				Rego: "package validate\n\ndefault validate = false",
				Output: &opa.OpaOutput{
					Validation: "invalid-path",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid observation path",
			spec: &opa.OpaSpec{
				Rego: "package validate\n\ndefault validate = false",
				Output: &opa.OpaOutput{
					Observations: []string{"invalid-path"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := opa.CreateOpaProvider(context.Background(), tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOpaProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
