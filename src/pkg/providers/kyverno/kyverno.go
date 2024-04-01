package kyverno

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	kube "github.com/defenseunicorns/lula/src/pkg/common/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
	kjson "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"

	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
)

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

		results, err := GetValidatedAssets(ctx, payload.Kyverno, collection, payload.Output)
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
					fmt.Errorf("expected status code 200 but got %d", resp.StatusCode)
			}

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return types.Result{}, err
			}

			contentType := resp.Header.Get("Content-Type")
			if contentType == "application/json" || contentType == "application/yaml" {

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

		results, err := GetValidatedAssets(ctx, payload.Kyverno, collection, payload.Output)
		if err != nil {
			return types.Result{}, err
		}
		return results, nil

	}

	return types.Result{}, fmt.Errorf("domain %s is not supported", domain)
}

func GetValidatedAssets(ctx context.Context, kyvernoPolicies *kjson.ValidatingPolicy, resources map[string]interface{}, output types.Output) (types.Result, error) {
	var matchResult types.Result

	if len(resources) == 0 {
		return matchResult, nil
	}

	if kyvernoPolicies == nil {
		return matchResult, fmt.Errorf("kyverno policy is not provided")
	}

	validationSet := make(map[string]map[string]bool)
	if output.Validation != "" {
		validationPairs := strings.Split(output.Validation, ",")

		for _, pair := range validationPairs {
			pair := strings.Split(pair, ".")

			if len(pair) != 2 {
				message.Debugf("Invalid validation pair: %v", pair)
				continue
			}

			validationPolicy := strings.TrimSpace(pair[0])
			validationRule := strings.TrimSpace(pair[1])
			if _, ok := validationSet[validationPolicy]; !ok {
				validationSet[validationPolicy] = make(map[string]bool)
			}
			validationSet[validationPolicy][validationRule] = true
		}
	}

	observationSet := make(map[string]map[string]bool)
	if len(output.Observations) > 0 {
		for _, observationPair := range output.Observations {
			pair := strings.Split(observationPair, ".")

			if len(pair) != 2 {
				message.Debugf("Invalid validation pair: %v", pair)
				continue
			}

			observationPolicy := strings.TrimSpace(pair[0])
			observationRule := strings.TrimSpace(pair[1])
			if _, ok := observationSet[observationPolicy]; !ok {
				observationSet[observationPolicy] = make(map[string]bool)
			}
			observationSet[observationPolicy][observationRule] = true
		}
	}

	policyarr := []*kjson.ValidatingPolicy{kyvernoPolicies}
	message.Debug(*policyarr[0])

	engine := jsonengine.New()
	response := engine.Run(ctx, jsonengine.Request{
		Resource: resources,
		Policies: policyarr,
	})

	observations := make(map[string]string)
	for i, policy := range response.Policies {
		for j, rule := range policy.Rules {
			if rule.Error != nil {
				message.Debugf("Error while evaluating rule: %v", rule.Error)
				continue
			}

			if _, ok := validationSet[policy.Policy.Name][rule.Rule.Name]; output.Validation == "" || ok {
				if len(rule.Violations) > 0 {
					matchResult.Failing += 1
				} else {
					matchResult.Passing += 1
				}
			}

			if _, ok := observationSet[policy.Policy.Name][rule.Rule.Name]; len(output.Observations) == 0 || ok {
				if len(rule.Violations) > 0 {
					observations[fmt.Sprintf("%s,%s-%d,%d", policy.Policy.Name, rule.Rule.Name, i, j)] = fmt.Sprintf("FAIL: %s", rule.Violations[0].Message)
				} else {
					observations[fmt.Sprintf("%s,%s-%d,%d", policy.Policy.Name, rule.Rule.Name, i, j)] = "PASS"
				}
			}
		}
	}

	matchResult.Observations = observations
	return matchResult, nil
}
