package validation

import (
	"fmt"
	"strings"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/defenseunicorns/lula/src/types"
)

// Requirement is an interface that defines the requirements for validation
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

// ComponentDefinitionRequirement implements a Requirement interface
// This is a requirement defined by a component definition
type ComponentDefinitionRequirement struct {
	// Unique ID of the control
	controlID string

	// Description of the control
	description string

	// Map of ImplementedRequirementUUID to DefinedComponent
	componentMap map[string]*oscalTypes.DefinedComponent

	// Map of ImplementedRequirementUUID to []*types.LulaValidation
	validationMap map[string][]*types.LulaValidation
}

func NewComponentDefinitionRequirement(controlID, description string, components map[string]*oscalTypes.DefinedComponent, validations map[string][]*types.LulaValidation) *ComponentDefinitionRequirement {
	return &ComponentDefinitionRequirement{
		controlID:     controlID,
		description:   description,
		componentMap:  components,
		validationMap: validations,
	}
}

func (r *ComponentDefinitionRequirement) ID() string {
	return r.controlID
}

func (r *ComponentDefinitionRequirement) Description() string {
	return r.description
}

// Success Criteria for component definition validations:
// all implemented requirements must have validations AND all validations must pass
// TODO: implement some future dynamic logic for success criteria?
func (r *ComponentDefinitionRequirement) EvaluateSuccess() (bool, string) {
	var msg strings.Builder
	var success bool

	for irUUID, component := range r.componentMap {
		validations, ok := r.validationMap[irUUID]
		if !ok {
			msg.WriteString(fmt.Sprintf("Component %s does not have any validations\n", component.Title))
			success = false
			continue
		}

		failingValidations := make([]string, 0)
		for _, validation := range validations {
			if validation.Result.Passing > 0 {
				continue
			}
			failingValidations = append(failingValidations, validation.UUID)
		}

		if len(failingValidations) > 0 {
			msg.WriteString(fmt.Sprintf("Component %s has the following validations that do not pass: %v\n", component.Title, failingValidations))
			success = false
			continue
		}
		success = true
	}

	return success, msg.String()
}

func (r *ComponentDefinitionRequirement) GetValidations() []*types.LulaValidation {
	var validations []*types.LulaValidation
	for _, v := range r.validationMap {
		validations = append(validations, v...)
	}
	return validations
}

// SimpleRequirement implements a Requirement interface
// It is a requirement that groups a set of validations
type SimpleRequirement struct {
	id          string
	validations []*types.LulaValidation
}

func NewSimpleRequirement(validations []*types.LulaValidation, id string) *SimpleRequirement {
	return &SimpleRequirement{
		id:          id,
		validations: validations,
	}
}

func (r *SimpleRequirement) ID() string {
	return r.id
}

func (r *SimpleRequirement) Description() string {
	return "All validations must pass"
}

// Requirement is successful if all linked validations are passing
// Note - this logic assumes that the absense of validations is a failure
func (r *SimpleRequirement) EvaluateSuccess() (bool, string) {
	var msg strings.Builder
	success := false
	passCount := 0

	for _, validation := range r.validations {
		if validation.Result.Passing > 0 {
			passCount++
			continue
		}
		msg.WriteString(fmt.Sprintf("validation %s failed\n", validation.UUID))
	}
	if passCount == len(r.validations) && passCount > 0 {
		return true, "all validations pass"
	}

	return success, msg.String()
}

func (r *SimpleRequirement) GetValidations() []*types.LulaValidation {
	return r.validations
}
