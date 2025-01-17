package validation

import (
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/types"
)

// A Validation Producer interface defines the requirements, how to meet them, and associated validations
type ValidationProducer interface {
	// Populate populates the validation store with the validations from the producer
	// as the requirements defined by the producer
	Populate(store *ValidationStore) error
}

// ComponentDefinitionProducer is a producer of validations referenced via OSCAL component definition
type ComponentDefinitionProducer struct {
	componentDefinition *oscal.ComponentDefinition
	requirements        []*ComponentDefinitionRequirement
}

func NewComponentProducer(compdef *oscal.ComponentDefinition, path, target string) *ComponentDefinitionProducer {
	// get oscal model from data
	// run some kind of composition here/import routines
	// how to incorporate target? or... is this external to this? Like maybe cmd validate wraps this and
	// passes the target?
	// get all requirements?
	return &ComponentDefinitionProducer{
		componentDefinition: compdef,
		requirements:        make([]*ComponentDefinitionRequirement, 0),
	}
}

func (c *ComponentDefinitionProducer) Populate(store *ValidationStore) error {
	// TODO: Get all validations from component definition

	// These could be stored in the backmatter or in a separate file

	// Get all requirements <> validations -> populate c.requirements

	return nil
}

// SimpleProducer
type SimpleProducer struct {
	validations []common.Validation
}

func NewSimpleProducer(validations []common.Validation) *SimpleProducer {
	return &SimpleProducer{
		validations: validations,
	}
}

func (p *SimpleProducer) Populate(store *ValidationStore) error {
	lulaValidations := make([]*types.LulaValidation, 0, len(p.validations))
	for _, validation := range p.validations {
		id, err := store.AddValidation(&validation)
		if err != nil {
			return err
		}
		lulaValidation, err := store.GetLulaValidation(id)
		if err != nil {
			return err
		}
		lulaValidations = append(lulaValidations, lulaValidation)
	}
	simpleReqt := NewSimpleRequirement(lulaValidations, "simple")

	store.AddRequirement(simpleReqt)

	return nil
}
