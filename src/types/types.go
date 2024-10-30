package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/defenseunicorns/lula/src/internal/transform"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

// Define base errors for validations
var (
	ErrExecutionNotAllowed = errors.New("execution not allowed")
	ErrDomainGetResources  = errors.New("domain GetResources error")
	ErrProviderEvaluate    = errors.New("provider Evaluate error")
)

type LulaValidationType string

const (
	LulaValidationTypeNormal  LulaValidationType = "Lula Validation"
	DefaultLulaValidationType LulaValidationType = LulaValidationTypeNormal
)

type LulaValidation struct {
	// Name of the Validation
	Name string

	// UUID of the validation - tied to the component-definition.backmatter
	UUID string

	// Provider is the provider that is evaluating the validation
	Provider *Provider

	// Domain is the domain that provides the evidence for the validation
	Domain *Domain

	// DomainResources is the set of resources that the domain is providing
	DomainResources *DomainResources

	// LulaValidationType is the type of validation that is being performed
	LulaValidationType LulaValidationType

	// Evaluated is a boolean that represents if the validation has been evaluated
	Evaluated bool

	// Tests is a slice of tests that are defined for the validation
	Tests *[]Test `json:"tests" yaml:"tests"`

	// Result is the result of the validation
	Result *Result
}

// CreateFailingLulaValidation creates a placeholder LulaValidation object that is always failing
func CreateFailingLulaValidation(name string) *LulaValidation {
	return &LulaValidation{
		Name:      name,
		Evaluated: true,
		Result:    &Result{Failing: 1},
	}
}

// CreatePassingLulaValidation creates a placeholder LulaValidation object that is always passing
func CreatePassingLulaValidation(name string) *LulaValidation {
	return &LulaValidation{
		Name:      name,
		Evaluated: true,
		Result:    &Result{Passing: 1},
	}
}

// LulaValidationMap is a map of LulaValidation objects
type LulaValidationMap = map[string]LulaValidation

// Lula Validation Options settings
type lulaValidationOptions struct {
	staticResources  DomainResources
	executionAllowed bool
	isInteractive    bool
	onlyResources    bool
	spinner          *message.Spinner
}

type LulaValidationOption func(*lulaValidationOptions)

// WithStaticResources sets the static resources for the LulaValidation object
func WithStaticResources(resources DomainResources) LulaValidationOption {
	return func(opts *lulaValidationOptions) {
		opts.staticResources = resources
	}
}

// ExecutionAllowed sets the value of the executionAllowed field in the LulaValidation object
func ExecutionAllowed(executionAllowed bool) LulaValidationOption {
	return func(opts *lulaValidationOptions) {
		opts.executionAllowed = executionAllowed
	}
}

// Interactive is a function that returns a boolean indicating if the validation should be interactive
func Interactive(isInteractive bool) LulaValidationOption {
	return func(opts *lulaValidationOptions) {
		opts.isInteractive = isInteractive
	}
}

// WithSpinner returns a LulaValidationOption that sets the spinner for the LulaValidation object
func WithSpinner(spinner *message.Spinner) LulaValidationOption {
	return func(opts *lulaValidationOptions) {
		opts.spinner = spinner
	}
}

// RequireExecutionConfirmation is a function that returns a boolean indicating if the validation requires confirmation before execution
func GetResourcesOnly(onlyResources bool) LulaValidationOption {
	return func(opts *lulaValidationOptions) {
		opts.onlyResources = onlyResources
	}
}

// Perform the validation, and store the result in the LulaValidation struct
func (val *LulaValidation) Validate(ctx context.Context, opts ...LulaValidationOption) error {
	if !val.Evaluated {
		var result Result
		var err error
		var resources DomainResources

		// Update the validation
		val.DomainResources = &resources
		val.Result = &result
		val.Evaluated = true

		// Set Validation config from options passed
		config := &lulaValidationOptions{
			staticResources:  nil,
			executionAllowed: false,
			isInteractive:    false,
			onlyResources:    false,
			spinner:          nil,
		}
		for _, opt := range opts {
			opt(config)
		}

		// Check if confirmation needed before execution
		if (*val.Domain).IsExecutable() && config.staticResources == nil {
			if !config.executionAllowed {
				if config.isInteractive {
					// Run confirmation user prompt
					if confirm := message.PromptForConfirmation(config.spinner); !confirm {
						return fmt.Errorf("%w: requested execution denied", ErrExecutionNotAllowed)
					}
				} else {
					return fmt.Errorf("%w: non-interactive execution not allowed", ErrExecutionNotAllowed)
				}
			}
		}

		// Get the resources
		if config.staticResources != nil {
			resources = config.staticResources
		} else {
			resources, err = (*val.Domain).GetResources(ctx)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrDomainGetResources, err)
			}
			if config.onlyResources {
				return nil
			}
		}

		// Perform the evaluation using the provider
		result, err = (*val.Provider).Evaluate(resources)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrProviderEvaluate, err)
		}
	}
	return nil
}

// RunTests executes any tests defined in the validation
// TODO: how to capture the test results? Want to execute all so don't want to return an error if one fails
func (v *LulaValidation) RunTests(ctx context.Context) error {
	if v.DomainResources == nil {
		return fmt.Errorf("domain resources are nil, tests cannot be run") // actually this probably isn't true...
	}

	// For each test, apply the transforms to the domain resources and run validate using those resources
	if v.Tests != nil {
		for _, test := range *v.Tests {
			// Create a fresh copy of the resources and validation to run each test on
			testResources := deepCopyMap(*v.DomainResources)
			testValidation := &LulaValidation{
				Name:     fmt.Sprintf("Test %s", test.Name),
				Provider: v.Provider,
				Domain:   v.Domain,
			}

			tt, err := transform.CreateTransformTarget(testResources)
			if err != nil {
				return fmt.Errorf("error creating transform target: %v", err)
			}

			for _, c := range test.Changes {
				testResources, err = tt.ExecuteTransform(c.Path, c.Type, c.Value, c.ValueMap)
				if err != nil {
					return fmt.Errorf("error executing transform: %v", err)
				}
			}

			err = testValidation.Validate(ctx, WithStaticResources(testResources))
			if err != nil {
				return err
			}

			// Check result
			if test.ExpectedResult == "pass" {
				if v.Result.Passing == 0 {
					return fmt.Errorf("expected passing test result, but got %d", v.Result.Passing)
				}
			} else if test.ExpectedResult == "fail" {
				if v.Result.Failing == 0 {
					return fmt.Errorf("expected failing test result, but got %d", v.Result.Failing)
				}
			}
		}
	} else {
		message.Debugf("No tests defined for validation %s", v.Name)
	}
	return nil
}

// Check if the validation requires confirmation before possible execution code is run
func (val *LulaValidation) RequireExecutionConfirmation() (confirm bool) {
	return !(*val.Domain).IsExecutable()
}

// Return domain resources as a json []byte
func (val *LulaValidation) GetDomainResourcesAsJSON() []byte {
	if val.DomainResources == nil {
		return []byte("{}")
	}
	jsonData, err := json.MarshalIndent(val.DomainResources, "", "  ")
	if err != nil {
		message.Debugf("Error marshalling domain resources to JSON: %v", err)
		jsonData = []byte(`{"Error": "Error marshalling to JSON"}`)
	}
	return jsonData
}

type DomainResources map[string]interface{}

type Domain interface {
	GetResources(context.Context) (DomainResources, error)
	IsExecutable() bool
}

type Provider interface {
	Evaluate(DomainResources) (Result, error)
}

// native type for conversion to targeted report format
type Result struct {
	UUID         string            `json:"uuid" yaml:"uuid"`
	ControlId    string            `json:"control-id" yaml:"control-id"`
	Description  string            `json:"description" yaml:"description"`
	Passing      int               `json:"passing" yaml:"passing"`
	Failing      int               `json:"failing" yaml:"failing"`
	State        string            `json:"state" yaml:"state"`
	Observations map[string]string `json:"observations" yaml:"observations"`
}

// Test is a struct that contains the name of the test, the permutations, and the expected result
type Test struct {
	Name           string   `json:"name" yaml:"name"`
	Changes        []Change `json:"changes" yaml:"changes"`
	ExpectedResult string   `json:"expected-result" yaml:"expected-result"`
}

type Change struct {
	Path     string                 `json:"path" yaml:"path"`
	Type     transform.ChangeType   `json:"type" yaml:"type"`
	Value    string                 `json:"value" yaml:"value"`
	ValueMap map[string]interface{} `json:"value-map" yaml:"value-map"`
}

func (c *Change) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("path is empty")
	}
	switch c.Type {
	case transform.ChangeTypeAdd, transform.ChangeTypeUpdate, transform.ChangeTypeDelete:
	default:
		return fmt.Errorf("invalid type")
	}

	return nil
}

func deepCopyMap(input map[string]interface{}) map[string]interface{} {
	if input == nil {
		return nil
	}

	// Create a new map to hold the copy
	copy := make(map[string]interface{})

	for key, value := range input {
		// Check the type of the value and copy accordingly
		switch v := value.(type) {
		case map[string]interface{}:
			// If the value is a map, recursively deep copy it
			copy[key] = deepCopyMap(v)
		case []interface{}:
			// If the value is a slice, deep copy each element
			copy[key] = deepCopySlice(v)
		default:
			// For other types (e.g., strings, ints), just assign directly
			copy[key] = v
		}
	}

	return copy
}

// Helper function to deep copy a slice of interface{}
func deepCopySlice(input []interface{}) []interface{} {
	if input == nil {
		return nil
	}

	// Create a new slice to hold the copy
	copy := make([]interface{}, len(input))

	for i, value := range input {
		// Check the type of the value and copy accordingly
		switch v := value.(type) {
		case map[string]interface{}:
			// If the value is a map, recursively deep copy it
			copy[i] = deepCopyMap(v)
		case []interface{}:
			// If the value is a slice, deep copy each element
			copy[i] = deepCopySlice(v)
		default:
			// For other types (e.g., strings, ints), just assign directly
			copy[i] = v
		}
	}

	return copy
}
