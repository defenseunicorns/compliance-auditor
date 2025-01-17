package validation

import (
	"context"
	"fmt"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"

	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

// Contains the store for the validations (their results, once executed) and associated requirements
type ValidationStore struct {
	validationMap map[string]*types.LulaValidation
	requirements  []Requirement
}

func NewValidationStore() *ValidationStore {
	return &ValidationStore{
		validationMap: make(map[string]*types.LulaValidation),
		requirements:  make([]Requirement, 0),
	}
}

// AddValidation adds a validation to the store
func (v *ValidationStore) AddValidation(validation *common.Validation) (id string, err error) {
	if validation.Metadata == nil {
		validation.Metadata = &common.Metadata{}
	}

	if validation.Metadata.UUID == "" {
		validation.Metadata.UUID = uuid.NewUUID()
	}

	lulaValidation, err := validation.ToLulaValidation(validation.Metadata.UUID)
	v.validationMap[validation.Metadata.UUID] = &lulaValidation

	if err != nil {
		return "", err
	}

	return validation.Metadata.UUID, nil
}

// AddLulaValidation adds a LulaValidation to the store
func (v *ValidationStore) AddLulaValidation(validation *types.LulaValidation, id string) {
	v.validationMap[id] = validation
}

// GetLulaValidation gets the LulaValidation from the store
func (v *ValidationStore) GetLulaValidation(id string) (validation *types.LulaValidation, err error) {
	if validation, ok := v.validationMap[id]; ok {
		return validation, nil
	}

	return validation, fmt.Errorf("validation #%s not found", id)
}

// AddRequirement adds a requirement to the store
func (v *ValidationStore) AddRequirement(requirement Requirement) {
	v.requirements = append(v.requirements, requirement)
}

// GetRequirements returns the requirements in the store
func (v *ValidationStore) GetRequirements() []Requirement {
	return v.requirements
}

// Count returns the number of validations in the store
func (v *ValidationStore) CountValidations() int {
	return len(v.validationMap)
}

// Count returns the number of requirements in the store
func (v *ValidationStore) CountRequirements() int {
	return len(v.requirements)
}

// Count requirements and associated validations?

func (v *ValidationStore) GetExecutable() []string {
	executableValidations := make([]string, 0)
	for k, val := range v.validationMap {
		if val != nil && val.Domain != nil {
			if (*val.Domain).IsExecutable() {
				executableValidations = append(executableValidations, k)
			}
		}
	}
	return executableValidations
}

// RunValidations executes the Lula Validations in the store
func (v *ValidationStore) RunValidations(ctx context.Context, confirmExecution, strict bool) error {
	for k, validation := range v.validationMap {
		err := validation.Validate(ctx, types.ExecutionAllowed(confirmExecution))
		if err != nil {
			message.Debugf("Error running validation %s: %v", k, err)
			// If strict is true, fail out
			if strict {
				return err
			}

			// Update validation with failed results
			validation.Result.State = "not-satisfied"
			validation.Result.Observations = map[string]string{
				"Error running validation": err.Error(),
			}
		}
		// Add the rest here ...

	}
	return nil
}

// TODO: Add SaveResources function
// TODO: Add RunTests function

// GetResults returns the results of the validations in the store
func (v *ValidationStore) GetResults() map[string]types.Result {
	results := make(map[string]types.Result)
	for _, validation := range v.validationMap {
		if validation.Result != nil {
			results[validation.UUID] = *validation.Result
		}
	}
	return results
}
