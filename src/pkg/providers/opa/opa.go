package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	kube "github.com/defenseunicorns/lula/src/pkg/common/kubernetes"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/mitchellh/mapstructure"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

func Validate(ctx context.Context, domain string, data map[string]interface{}) (types.Result, error) {

	// Convert map[string]interface to a RegoTarget
	var payload types.Payload
	err := mapstructure.Decode(data, &payload)
	if err != nil {
		return types.Result{}, err
	}

	// TODO: Start here
	// need to create a single map[string]interface{}
	// What is the top level? or is there not one?
	// IE map["pods-vt"]interface{}
	// map[""pods-other]interface{}

	// Given that this is executed per-target - there may never be a need for a slice?
	collection := make(map[string]interface{}, 0)
	if domain == "kubernetes" {
		collection, err = kube.QueryCluster(ctx, payload.Resources)
		if err != nil {
			return types.Result{}, err
		}
	} else {
		return types.Result{}, fmt.Errorf("domain %s is not supported", domain)
	}

	// TODO: Add logging optionality for understanding what resources are actually being validated
	results, err := GetValidatedAssets(ctx, payload.Rego, collection)
	if err != nil {
		return types.Result{}, err
	}
	// return results

	return results, nil
}

// GetValidatedAssets performs the validation of the dataset against the given rego policy
func GetValidatedAssets(ctx context.Context, regoPolicy string, dataset map[string]interface{}) (types.Result, error) {
	var matchResult types.Result

	compiler, err := ast.CompileModules(map[string]string{
		"validate.rego": regoPolicy,
	})
	if err != nil {
		log.Fatal(err)
		return matchResult, fmt.Errorf("failed to compile rego policy: %w", err)
	}

	regoCalc := rego.New(
		rego.Query("data.validate"),
		rego.Compiler(compiler),
		rego.Input(dataset),
	)

	resultSet, err := regoCalc.Eval(ctx)

	if err != nil || resultSet == nil || len(resultSet) == 0 {
		return matchResult, fmt.Errorf("failed to evaluate rego policy: %w", err)
	}

	for _, result := range resultSet {
		for _, expression := range result.Expressions {
			expressionBytes, err := json.Marshal(expression.Value)
			if err != nil {
				return matchResult, fmt.Errorf("failed to marshal expression: %w", err)
			}

			var expressionMap map[string]interface{}
			err = json.Unmarshal(expressionBytes, &expressionMap)
			if err != nil {
				return matchResult, fmt.Errorf("failed to unmarshal expression: %w", err)
			}
			// TODO: add logging optionality here for developer experience
			if matched, ok := expressionMap["validate"]; ok && matched.(bool) {
				// fmt.Printf("Asset %s matched policy: %s\n\n", asset, expression)
				matchResult.Passing += 1
			} else {
				// fmt.Printf("Asset %s no matched policy: %s\n\n", asset, expression)
				matchResult.Failing += 1
			}
		}
	}

	return matchResult, nil
}
