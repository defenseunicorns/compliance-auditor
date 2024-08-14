package kyverno_test

import (
	"context"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/providers/kyverno"
	kjson "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
)

func TestCreateKyvernoProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		spec    *kyverno.KyvernoSpec
		wantErr bool
	}{
		{
			name: "valid spec",
			spec: &kyverno.KyvernoSpec{
				Policy: &kjson.ValidatingPolicy{
					Spec: kjson.ValidatingPolicySpec{},
				},
			},
			wantErr: false,
		},
		{
			name: "valid spec with output",
			spec: &kyverno.KyvernoSpec{
				Policy: &kjson.ValidatingPolicy{
					Spec: kjson.ValidatingPolicySpec{},
				},
				Output: &kyverno.KyvernoOutput{
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
			name:    "nil policy",
			spec:    &kyverno.KyvernoSpec{},
			wantErr: true,
		},
		{
			name: "empty policy",
			spec: &kyverno.KyvernoSpec{
				Policy: &kjson.ValidatingPolicy{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := kyverno.CreateKyvernoProvider(context.Background(), tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateKyvernoProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
