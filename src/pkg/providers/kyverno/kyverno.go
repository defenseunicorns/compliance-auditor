package kyverno

import (
	"context"
	"fmt"
	"strings"

	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
	kjson "github.com/kyverno/kyverno-json/pkg/apis/policy/v1alpha1"

	jsonengine "github.com/kyverno/kyverno-json/pkg/json-engine"
)

func GetValidatedAssets(ctx context.Context, kyvernoPolicies *kjson.ValidatingPolicy, resources map[string]interface{}, output KyvernoOutput) (types.Result, error) {
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
