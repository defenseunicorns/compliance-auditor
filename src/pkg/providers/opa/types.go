package opa

import (
	"context"
	"fmt"
	"strings"

	"github.com/defenseunicorns/lula/src/types"
)

type OpaProvider struct {
	// Context is the context that the OPA policy is being evaluated in
	Context context.Context `json:"context" yaml:"context"`

	// Spec is the specification of the OPA policy
	Spec *OpaSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func CreateOpaProvider(ctx context.Context, spec *OpaSpec) (types.Provider, error) {
	// Check validity of spec
	if spec == nil {
		return nil, fmt.Errorf("spec is nil")
	}

	if spec.Rego == "" {
		return nil, fmt.Errorf("rego policy cannot be empty")
	}

	if spec.Output != nil {
		if spec.Output.Validation != "" {
			if !strings.Contains(spec.Output.Validation, ".") {
				return nil, fmt.Errorf("validation field must be a json path")
			}
		}
		if spec.Output.Observations != nil {
			for _, observation := range spec.Output.Observations {
				if !strings.Contains(observation, ".") {
					return nil, fmt.Errorf("observation field must be a json path")
				}
			}
		}
	}

	return OpaProvider{
		Context: ctx,
		Spec:    spec,
	}, nil
}

func (o OpaProvider) Evaluate(resources types.DomainResources) (types.Result, error) {
	results, err := GetValidatedAssets(o.Context, o.Spec.Rego, resources, o.Spec.Output)
	if err != nil {
		return types.Result{}, err
	}
	return results, nil
}

// OpaSpec is the specification of the OPA policy, required if the provider type is opa
type OpaSpec struct {
	// Required: Rego is the OPA policy
	Rego string `json:"rego" yaml:"rego"`
	// Optional: Output is the output of the OPA policy
	Output *OpaOutput `json:"output,omitempty" yaml:"output,omitempty"`
}

// OpaOutput Defines the output structure for OPA validation results, including validation status and additional observations.
type OpaOutput struct {
	// optional: Specifies the JSON path to a boolean value indicating the validation result.
	Validation string `json:"validation" yaml:"validation"`
	// optional: any additional observations to include (fields must resolve to strings)
	Observations []string `json:"observations" yaml:"observations"`
}
