package opa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"net/http"

	"github.com/defenseunicorns/lula/src/pkg/common/kubernetes"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/mitchellh/mapstructure"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

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

	// Need []map[string]interface{} for rego validation
	var mapData []map[string]interface{}
	
	// query kubernetes for resource data if domain == "kubernetes"
	// TODO: evaluate processes for manifests/helm charts
	if domain == "kubernetes" {
		var resources []unstructured.Unstructured
		resources, err = kube.QueryCluster(ctx, payload)
		if err != nil {
			return types.Result{}, err
		}
		for _, item := range resources {
			mapData = append(mapData, item.Object)
		}

	} else if domain == "api" {
		transport := &http.Transport{}
		client := &http.Client{Transport: transport}

		resp, err := client.Get("http://localhost/api-field")
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var prettyBuff bytes.Buffer
		json.Indent(&prettyBuff, body, "", "  ")
		prettyString := prettyBuff.String()
		prettyString = "["+prettyString+"]"

		err = json.Unmarshal([]byte(prettyString), &mapData)
		if err != nil {
			log.Fatal(err)
		}
		

	} else {
		return types.Result{}, fmt.Errorf("domain %s is not supported", domain)
	}

	// TODO: Add logging optionality for understanding what resources are actually being validated
	results, err := GetValidatedAssets(ctx, payload.Rego, mapData)
	if err != nil {
		return types.Result{}, err
	}
	// return results

	return results, nil
}

// GetValidatedAssets performs the validation of the dataset against the given rego policy
func GetValidatedAssets(ctx context.Context, regoPolicy string, dataset []map[string]interface{}) (types.Result, error) {
	var wg sync.WaitGroup
	var matchResult types.Result

	compiler, err := ast.CompileModules(map[string]string{
		"validate.rego": regoPolicy,
	})
	if err != nil {
		log.Fatal(err)
		return matchResult, fmt.Errorf("failed to compile rego policy: %w", err)
	}

	fmt.Printf("Applying policy against %s resources\n", strconv.Itoa(len(dataset)))
	for _, asset := range dataset {
		wg.Add(1)
		go func(asset map[string]interface{}) {
			defer wg.Done()

			regoCalc := rego.New(
				rego.Query("data.validate"),
				rego.Compiler(compiler),
				rego.Input(asset),
			)

			resultSet, err := regoCalc.Eval(ctx)
			if err != nil || resultSet == nil || len(resultSet) == 0 {
				wg.Done()
			}

			for _, result := range resultSet {
				for _, expression := range result.Expressions {
					expressionBytes, err := json.Marshal(expression.Value)
					if err != nil {
						wg.Done()
					}

					var expressionMap map[string]interface{}
					err = json.Unmarshal(expressionBytes, &expressionMap)
					if err != nil {
						wg.Done()
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
		}(asset)
	}

	wg.Wait()

	return matchResult, nil
}
