package types

type LulaValidationType string

const (
	LulaValidationTypeNormal  LulaValidationType = "Lula Validation"
	DefaultLulaValidationType LulaValidationType = LulaValidationTypeNormal
)

type LulaValidation struct {
	// Provider is the provider that is evaluating the validation
	Provider Provider

	// Domain is the domain that provides the evidence for the validation
	Domain Domain

	// DomainResources is the set of resources that the domain is providing
	DomainResources DomainResources

	// LulaValidationType is the type of validation that is being performed
	LulaValidationType LulaValidationType

	// Evaluated is a boolean that represents if the validation has been evaluated
	Evaluated bool

	// Result is the result of the validation
	Result Result
}

// LulaValidationMap is a map of LulaValidation objects
type LulaValidationMap = map[string]LulaValidation

// Perform the validation, and store the result in the LulaValidation struct
func (val *LulaValidation) Validate() error {
	if !val.Evaluated {
		// Extract the resources from the domain
		domainResources, err := val.Domain.GetResources()
		if err != nil {
			return err
		}
		// Bookkeeping of the domain resources for use elsewhere
		val.DomainResources = domainResources
		// Perform the evaluation using the provider
		result, err := val.Provider.Evaluate(domainResources)
		if err != nil {
			return err
		}
		// Store the result in the validation object
		val.Result = result
		val.Evaluated = true
	}
	return nil
}

type DomainResources map[string]interface{}

type Domain interface {
	GetResources() (DomainResources, error)
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
