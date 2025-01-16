# 8. Validation Refactor

Date: 2024-01-16

## Status

Review

## Context

Execution of the Lula Validation Engine (i.e., the underlying functionality that performs evidence collection from a system of interest and evaluation of that against some policy) is tightly coupled to the OSCAL format. This is dis-advantageous for a few reasons:
* Other data formats that wrap Lula Validations (e.g., Assessment Plans or a future Lula-custom aggregated validation format) are not supported. At minimum the OSCAL Assessment Plan should be a supported input to correctly instantiate the [OSCAL e2e model flow](https://pages.nist.gov/OSCAL/resources/concepts/layer/assessment/assessment-results/).
* The current package structure does not allow for core functionality of Lula to be easily reused or extended. For example, a Kubernetes controller may want to use some Lula features in business logic, but not couple to the OSCAL format.

## Decision

To alleviate the rigidity of the current design, parts of the Lula library will need to be refactored, resulting in breaking changes, but ideally creating an implementation that is more modular, reusable, and generally not dependent on the data structure for Component Definitions and Assessment Results.

The primary outcome of this will be exposed functionality that can be thought of in following categories:

* Validation retrieval from source - Encapsulates extracting Lula Validations from source and loading the validations and associated requirements into the `ValidationStore`
* Execution of validations - Running the Validation Engine against `ValidationStore` validations
* Reporting to output - Realizing the results of the validation (i.e., the policy outputs) along with the requirements into the desired output format

## Consequences

Currently the validation logic exists in `src/pkg/common`, much of this logic will be refactored into `src/pkg/validation`. New interfaces will be introduced to support the decoupled validaiton logic.

The following interfaces are proposed:

```go
// A Validation Producer interface defines the requirements, how to meet them, and associated validations
type ValidationProducer interface {
	// Populate populates the validation store with the validations from the producer
	// as well as the associated requirements, as defined by the producer
	Populate(store *ValidationStore) error
}

// Requirement is an interface that defines the requirements for validation
// These are specific to the provider and contain custom logic that evaluate the "satisfaction" of the
// validation against said requirement.
type Requirement interface {
	// ID returns the unique requirement ID
	ID() string

	// Description returns a requirement description
	Description() string

	// EvaluateSuccess determines if the requirement is satisfied (true) and returns a message
	// describing the criteria evaluation
	EvaluateSuccess() (bool, string)

	// GetValidations returns the validations associated with the requirement
	GetValidations() []*types.LulaValidation
}

// ResultConsumer is the interface that must be implemented by any consumer of the validation
// store and results. It is responsible for evaluating the results and generating the output
// speific to the consumer.
type ResultConsumer interface {
	// Evaluate Results are the custom implementation for the consumer, which should take the
	// requirements, as specified by the producer, plus the data in the validation store
	// and evaluate them + generate the output
	EvaluateResults(store *ValidationStore) error

	// Generate Output is the custom implementation for the consumer that should create
	// a custom output
	GenerateOutput() error
}
```

The `ValidationStore` would be redefined to be non-specific to a Component Definition/Back-Matter data storage, but instead be a more generic store of the validations and requirements.

```go
// Contains the store for the validations (their results, once executed) and associated requirements
type ValidationStore struct {
	validationMap map[string]*types.LulaValidation
	requirements  []Requirement
}
```

The validation store would allow for adding validations and requirements, as well as executing the validations and retrieving the results.

### Example Implementation

An example implementation of the above interfaces is on the `test-lib-refactor` branch of the lula repo. Specifically, this explores implementing a Component Definition as a `ValidationProducer` and the Assessment Results as a `ResultConsumer`. As well as re-defining the `ValidationStore` to be a more generic store of validations and requirements.