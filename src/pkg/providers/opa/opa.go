package opa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	kube "github.com/defenseunicorns/lula/src/pkg/common/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

// TODO: What is the new version of the information we are displaying on the command line?

func Validate(ctx context.Context, domain string, data types.Target) (types.Result, error) {
	if domain == "kubernetes" {
		payload := data.Payload

		err := kube.EvaluateWait(payload.Wait)
		if err != nil {
			return types.Result{}, err
		}

		collection, err := kube.QueryCluster(ctx, payload.Resources)
		if err != nil {
			return types.Result{}, err
		}

		// TODO: Add logging optionality for understanding what resources are actually being validated
		results, err := GetValidatedAssets(ctx, payload.Rego, collection, payload.Output)
		if err != nil {
			return types.Result{}, err
		}

		return results, nil

	} else if domain == "api" {
		payload := data.Payload

		collection := make(map[string]interface{}, 0)

		for _, request := range payload.Requests {
			transport := &http.Transport{}
			client := &http.Client{Transport: transport}

			resp, err := client.Get(request.URL)
			if err != nil {
				return types.Result{}, err
			}
			if resp.StatusCode != 200 {
				return types.Result{},
					fmt.Errorf("expected status code 200 but got %d\n", resp.StatusCode)
			}

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return types.Result{}, err
			}

			contentType := resp.Header.Get("Content-Type")
			if contentType == "application/json" {

				var prettyBuff bytes.Buffer
				json.Indent(&prettyBuff, body, "", "  ")
				prettyJson := prettyBuff.String()

				var tempData interface{}
				err = json.Unmarshal([]byte(prettyJson), &tempData)
				if err != nil {
					return types.Result{}, err
				}
				collection[request.Name] = tempData

			} else {
				return types.Result{}, fmt.Errorf("content type %s is not supported", contentType)
			}
		}

		results, err := GetValidatedAssets(ctx, payload.Rego, collection, payload.Output)
		if err != nil {
			return types.Result{}, err
		}
		return results, nil

	}

	return types.Result{}, fmt.Errorf("domain %s is not supported", domain)
}

// GetValidatedAssets performs the validation of the dataset against the given rego policy
func GetValidatedAssets(ctx context.Context, regoPolicy string, dataset map[string]interface{}, output types.Output) (types.Result, error) {
	var matchResult types.Result

	if len(dataset) == 0 {
		// Not an error but no entries to validate
		// TODO: add a warning log
		return matchResult, nil
	}

	compiler, err := ast.CompileModules(map[string]string{
		"validate.rego": regoPolicy,
	})
	if err != nil {
		log.Fatal(err)
		return matchResult, fmt.Errorf("failed to compile rego policy: %w", err)
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
		return matchResult, fmt.Errorf("failed to evaluate rego policy: %w", err)
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
			return matchResult, fmt.Errorf("failed to evaluate rego policy: %w", err)
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
