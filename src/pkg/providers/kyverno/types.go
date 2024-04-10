package kyverno

import (
	"context"

	"github.com/defenseunicorns/lula/src/types"
	kjson "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"
)

type KyvernoProvider struct {
	// Context is the context that the Kyverno policy is being evaluated in
	Context context.Context `json:"context" yaml:"context"`

	// Spec is the specification of the Kyverno policy
	Spec KyvernoSpec `json:"spec" yaml:"spec"`
}

func (k KyvernoProvider) Evaluate(resources types.DomainResources) (types.Result, error) {
	results, err := GetValidatedAssets(k.Context, k.Spec.Policy, resources, k.Spec.Output)
	if err != nil {
		return types.Result{}, err
	}
	return results, nil
}

type KyvernoSpec struct {
	Policy *kjson.ValidatingPolicy `json:"policy" yaml:"policy"`
	Output KyvernoOutput           `json:"output" yaml:"output"`
}

type KyvernoOutput struct {
	Validation   string   `json:"validation" yaml:"validation"`
	Observations []string `json:"observations" yaml:"observations"`
}
