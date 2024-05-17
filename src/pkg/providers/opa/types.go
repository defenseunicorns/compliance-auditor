package opa

import (
	"context"

	"github.com/defenseunicorns/lula/src/types"
)

type OpaProvider struct {
	// Context is the context that the OPA policy is being evaluated in
	Context context.Context `json:"context" yaml:"context"`

	// Spec is the specification of the OPA policy
	Spec *OpaSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func (o OpaProvider) Evaluate(resources types.DomainResources) (types.Result, error) {
	results, err := GetValidatedAssets(o.Context, o.Spec.Rego, resources, o.Spec.Output)
	if err != nil {
		return types.Result{}, err
	}
	return results, nil
}

type OpaSpec struct {
	Rego   string     `json:"rego" yaml:"rego"`
	Output *OpaOutput `json:"output,omitempty" yaml:"output,omitempty"`
}

type OpaOutput struct {
	Validation   string   `json:"validation" yaml:"validation"`
	Observations []string `json:"observations" yaml:"observations"`
}
