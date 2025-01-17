package validation

import (
	"context"
)

// Contains the business logic around collecting and returning Lula Validations

type Validator struct {
	// Extracts validations from any source and populates the store
	producer ValidationProducer

	// Processes the final results after validation execution.
	consumer ResultConsumer

	// Contains the validations and requirements
	store *ValidationStore

	// Variables to store validator configuration behaviors
	outputDir                    string
	requestExecutionConfirmation bool
	runExecutableValidations     bool
	saveResources                bool
	strict                       bool
	silent                       bool
}

// Create a new validator
func New(producer ValidationProducer, consumer ResultConsumer, opts ...Option) (*Validator, error) {
	var validator Validator

	for _, opt := range opts {
		if err := opt(&validator); err != nil {
			return nil, err
		}
	}

	validator.store = NewValidationStore()
	validator.producer = producer
	validator.consumer = consumer

	err := validator.producer.Populate(validator.store)
	if err != nil {
		return nil, err
	}

	return &validator, nil
}

// ExecuteValidations collects the validations, executes, and provides the results in the specified consumer
func (v *Validator) ExecuteValidations(ctx context.Context, runExecutableValidations bool) error {
	// Run the validations
	err := v.store.RunValidations(ctx, runExecutableValidations, v.strict)
	if err != nil {
		return err
	}

	// Consumer evaluates the results -> this should execute their custom output routines
	err = v.consumer.EvaluateResults(v.store)
	if err != nil {
		return err
	}

	return nil
}

func (v *Validator) GetStats() (int, int, int) {
	executableValidations := v.store.GetExecutable()
	return v.store.CountRequirements(),
		v.store.CountValidations(),
		len(executableValidations)
}

func (v *Validator) YieldResults() error {
	// Execute the consumer GenerateOutput function
	return v.consumer.GenerateOutput()
}
