package types

import (
	"errors"
	"fmt"

	"github.com/defenseunicorns/lula/src/pkg/message"
)

type LulaValidationType string

const (
	LulaValidationTypeNormal  LulaValidationType = "Lula Validation"
	DefaultLulaValidationType LulaValidationType = LulaValidationTypeNormal
)

type LulaValidation struct {
	// Name of the Validation
	Name string

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
func (val *LulaValidation) Validate(opts ...LulaValidationOption) error {
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
						return errors.New("execution not allowed")
					}
				} else {
					return errors.New("execution not allowed")
				}
			}
		}

		// Get the resources
		if config.staticResources != nil {
			resources = config.staticResources
		} else {
			resources, err = (*val.Domain).GetResources()
			if err != nil {
				return fmt.Errorf("domain GetResources error: %v", err)
			}
			if config.onlyResources {
				return nil
			}
		}

		// Perform the evaluation using the provider
		result, err = (*val.Provider).Evaluate(resources)
		if err != nil {
			return fmt.Errorf("provider Evaluate error: %v", err)
		}
	}
	return nil
}

// Check if the validation requires confirmation before possible execution code is run
func (val *LulaValidation) RequireExecutionConfirmation() (confirm bool) {
	return !(*val.Domain).IsExecutable()
}

type DomainResources map[string]interface{}

type Domain interface {
	GetResources() (DomainResources, error)
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
