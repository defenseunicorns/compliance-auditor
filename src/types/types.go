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

// LulaValidationMap is a map of LulaValidation objects
type LulaValidationMap = map[string]LulaValidation

// LulaValidationLinksMap is map of an array of LulaValidations
type LulaValidationLinksMap = map[string][]*LulaValidation

// Lula Validation Options settings
type lulaValidationOptions struct {
	staticResources  DomainResources
	confirmExecution bool
}

type LulaValidationOption func(*lulaValidationOptions)

// WithStaticResources sets the static resources for the LulaValidation object
func WithStaticResources(resources DomainResources) LulaValidationOption {
	return func(opts *lulaValidationOptions) {
		opts.staticResources = resources
	}
}

// RequireExecutionConfirmation is a function that returns a boolean indicating if the validation requires confirmation before execution
func RequireExecutionConfirmation(confirmationFlag bool) LulaValidationOption {
	return func(opts *lulaValidationOptions) {
		opts.confirmExecution = !confirmationFlag
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
			confirmExecution: true,
		}
		for _, opt := range opts {
			opt(config)
		}

		// Check if confirmation is required before execution
		if config.confirmExecution && val.Domain.IsExecutable() {
			// Run confirmation user prompt
			confirm := message.PromptForConfirmation()
			if !confirm {
				return errors.New("execution not confirmed")
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
	return !val.Domain.IsExecutable()
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
