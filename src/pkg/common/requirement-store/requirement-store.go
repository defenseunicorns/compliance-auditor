package requirementstore

import (
	"fmt"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

type RequirementStore struct {
	RequirementMap    map[string]oscalTypes_1_1_2.ImplementedRequirementControlImplementation
	ValidationLinkMap types.LulaValidationLinksMap
	ValidationStore   *validationstore.ValidationStore
	FindingsMap       map[string]oscalTypes_1_1_2.Finding
}

// NewRequirementStore creates a new requirement store from component defintion
func NewRequirementStore(componentDef *oscalTypes_1_1_2.ComponentDefinition, validationStore *validationstore.ValidationStore) *RequirementStore {
	return &RequirementStore{
		RequirementMap:    oscal.ComponentDefinitionToRequirementMap(componentDef),
		ValidationLinkMap: make(types.LulaValidationLinksMap),
		ValidationStore:   validationStore,
		FindingsMap:       make(map[string]oscalTypes_1_1_2.Finding),
	}
}

// UpdateRequirementStoreWithLulaValidations adds Lula validations to the store
func (r *RequirementStore) UpdateRequirementStoreWithLulaValidations() {
	// get all Lula validations linked to the requirement
	for _, requirement := range r.RequirementMap {
		if requirement.Links != nil {
			for _, link := range *requirement.Links {
				if common.IsLulaLink(link) {
					ids, err := r.ValidationStore.AddFromLink(link)
					if err != nil {
						message.Debugf("Error adding validation from link %s: %v", link.Href, err)
						// Create a new validation and add it to the validationLinkMap
						tmpValidation := &types.LulaValidation{
							Evaluated: true,
							Result: types.Result{
								State: "not-satisfied",
								Observations: map[string]string{
									fmt.Sprintf("Error adding validation from link %s", link.Href): err.Error(),
								},
							},
						}
						r.ValidationLinkMap[requirement.UUID] = append(r.ValidationLinkMap[requirement.UUID], tmpValidation)
						continue
					}
					for _, id := range ids {
						validation, err := r.ValidationStore.GetLulaValidation(id)
						if err != nil {
							message.Debugf("Error getting validation %s: %v", id, err)
							// Create a new validation and add it to the validationLinkMap
							tmpValidation := &types.LulaValidation{
								Evaluated: true,
								Result: types.Result{
									State: "not-satisfied",
									Observations: map[string]string{
										fmt.Sprintf("Error getting validation %s", id): err.Error(),
									},
								},
							}
							r.ValidationLinkMap[requirement.UUID] = append(r.ValidationLinkMap[requirement.UUID], tmpValidation)
							continue
						}
						r.ValidationLinkMap[requirement.UUID] = append(r.ValidationLinkMap[requirement.UUID], validation)
					}
				}
			}
		}
	}
}

// RunValidations runs the validations in the store
func (r *RequirementStore) RunValidations() {
	for _, validations := range r.ValidationStore. {
}

// GenerateFindings generates the findings in the store
func (r *RequirementStore) GenerateFindings() {
	// For each implemented requirement and linked validation, create a finding/observation
}

// Drop unused validations from store (only relevant when store created from the backmatter if it contains unused validations)
// func (r *RequirementStore) DropUnusedValidations() {
// 	TODO
// }

// AddFinding adds a finding to the store

// ReturnObservations returns the observations from the store (subset of findings)
