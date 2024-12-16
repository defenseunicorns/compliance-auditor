package opa

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

var (
	ErrCompileRego  = errors.New("failed to compile rego policy")
	ErrEvaluateRego = errors.New("failed to evaluate rego policy")
)

// mainPolicyModuleName is the name of the OPA module containing the main policy from the spec.rego field.
const mainPolicyModuleName = "validate.rego"

// GetValidatedAssets performs the validation of the dataset against the given rego policy
func GetValidatedAssets(ctx context.Context, regoPolicy string, regoModules map[string]string, dataset map[string]interface{}, output *OpaOutput) (types.Result, error) {
	var matchResult types.Result

	if len(dataset) == 0 {
		return matchResult, errors.New("opa validation not performed - no resources to validate")
	}

	if output == nil {
		output = &OpaOutput{}
	}

	modules := make(map[string]string, len(regoModules)+1)
	for k, v := range regoModules {
		modules[k] = v
	}
	modules[mainPolicyModuleName] = regoPolicy

	compiler, err := ast.CompileModules(modules)
	if err != nil {
		message.Debugf("failed to compile rego policy: %s", err.Error())
		return matchResult, fmt.Errorf("%w: %w", ErrCompileRego, err)
	}

	// Get validation decision
	validation := "validate.validate"
	if output.Validation != "" {
		validation = output.Validation
	}

	regoCalcValid := rego.New(
		rego.Query(fmt.Sprintf("data.%s", validation)),
		rego.Compiler(compiler),
		rego.Input(dataset),
	)

	resultValid, err := regoCalcValid.Eval(ctx)
	if err != nil {
		return matchResult, fmt.Errorf("%w: %w", ErrEvaluateRego, err)
	}
	// Checking result length is non-zero: will be zero if validation returns false
	if len(resultValid) != 0 {
		// Extra check on validation value = true, to ensure it's a boolean return since it could be anything
		if matched, ok := resultValid[0].Expressions[0].Value.(bool); ok && matched {
			matchResult.Passing += 1
		} else {
			matchResult.Failing += 1
			if !ok {
				message.Debugf("Validation field expected bool and got %s", reflect.TypeOf(resultValid[0].Expressions[0].Value))
			}
		}
	} else {
		matchResult.Failing += 1
	}

	// Get additional observations, if they exist - only supports string output
	observations := make(map[string]string)
	for _, obv := range output.Observations {
		regoCalcObv := rego.New(
			rego.Query(fmt.Sprintf("data.%s", obv)),
			rego.Compiler(compiler),
			rego.Input(dataset),
		)

		resultObv, err := regoCalcObv.Eval(ctx)
		if err != nil {
			return matchResult, fmt.Errorf("%w: %w", ErrEvaluateRego, err)
		}
		// To do: check if resultObv is empty - basically some extra error handling if a user defines an output but it's not coming out of the rego
		if len(resultObv) != 0 {
			if matched, ok := resultObv[0].Expressions[0].Value.(string); ok {
				observations[obv] = matched
			} else {
				message.Debugf("Observation field %s expected string and got %s", obv, reflect.TypeOf(resultObv[0].Expressions[0].Value))
			}
		} else {
			message.Debugf("Observation field %s not output from rego", obv)
		}
	}
	matchResult.Observations = observations

	return matchResult, nil
}
