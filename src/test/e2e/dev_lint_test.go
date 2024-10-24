package test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/dev"
)

func TestLintCommand(t *testing.T) {

	// Define the test cases
	testCases := []struct {
		name       string
		inputFiles []string
		valid      []bool
	}{
		{
			name:       "Valid multi validation file",
			inputFiles: []string{"../../test/e2e/scenarios/dev-lint/multi.validation.yaml"},
			valid:      []bool{true, true},
		},
		{
			name:       "Valid OPA validation file",
			inputFiles: []string{"../../test/e2e/scenarios/dev-lint/opa.validation.yaml"},
			valid:      []bool{true},
		},
		{
			name:       "Valid Kyverno validation file",
			inputFiles: []string{"../../test/e2e/scenarios/dev-lint/validation.kyverno.yaml"},
			valid:      []bool{true},
		},
		{
			name:       "Invalid OPA validation file",
			inputFiles: []string{"../../test/e2e/scenarios/dev-lint/invalid.opa.validation.yaml"},
			valid:      []bool{false},
		},
		{
			name:       "Multiple files",
			inputFiles: []string{"../../test/e2e/scenarios/dev-lint/validation.kyverno.yaml", "../../test/e2e/scenarios/dev-lint/invalid.opa.validation.yaml"},
			valid:      []bool{true, false},
		},
		{
			name:       "Remote validation file",
			inputFiles: []string{"https://raw.githubusercontent.com/defenseunicorns/lula/main/src/test/e2e/scenarios/dev-validate/validation.kyverno.yaml"},
			valid:      []bool{true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validationResults := dev.DevLintCommand(tc.inputFiles)
			for i, result := range validationResults {
				if result.Valid != tc.valid[i] {
					t.Errorf("Expected valid to be %v, but got %v", tc.valid[i], result.Valid)
				}
			}
		})
	}
}
