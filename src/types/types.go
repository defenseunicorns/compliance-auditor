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

	// LulaValidationType is the type of validation that is being performed
	LulaValidationType LulaValidationType

	// Evaluated is a boolean that represents if the validation has been evaluated
	Evaluated bool

	// Result is the result of the validation
	Result Result
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
