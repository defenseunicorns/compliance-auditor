package requirementstore

import (
	"fmt"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

type RequirementStore struct {
	requirementMap map[string]oscal.Requirement
	findingMap     map[string]oscalTypes.Finding
}

type Stats struct {
	TotalRequirements        int
	TotalValidations         int
	ExecutableValidations    bool
	ExecutableValidationsMsg string
	TotalFindings            int
}

// NewRequirementStore creates a new requirement store from component defintion
func NewRequirementStore(controlImplementations *[]oscalTypes.ControlImplementationSet) *RequirementStore {
	return &RequirementStore{
		requirementMap: oscal.ControlImplementationstToRequirementsMap(controlImplementations),
		findingMap:     make(map[string]oscalTypes.Finding),
	}
}

// ResolveLulaValidations resolves the linked Lula validations with the requirements and populates the ValidationStore.validationMap
func (r *RequirementStore) ResolveLulaValidations(validationStore *validationstore.ValidationStore) {
	// get all Lula validations linked to the requirement
	var lulaValidation *types.LulaValidation
	for _, requirement := range r.requirementMap {
		if requirement.ImplementedRequirement.Links != nil {
			for _, link := range *requirement.ImplementedRequirement.Links {
				if common.IsLulaLink(link) {
					_, err := validationStore.GetLulaValidation(link.Href)
					if err != nil {
						message.Debugf("Error adding validation from link %s: %v", link.Href, err)
						// Create new LulaValidation and add to validationStore
						lulaValidation = types.CreateFailingLulaValidation("lula-validation-error")
						lulaValidation.Result.Observations = map[string]string{
							fmt.Sprintf("Error getting Lula validation %s", link.Href): err.Error(),
						}
						validationStore.AddLulaValidation(lulaValidation, link.Href)
					}
				}
			}
		}
	}
}

// GenerateFindings generates the findings in the store
func (r *RequirementStore) GenerateFindings(validationStore *validationstore.ValidationStore) map[string]oscalTypes.Finding {
	// For each implemented requirement and linked validation, create a finding/observation
	for _, requirement := range r.requirementMap {
		// This should produce a finding - check if an existing finding for the control-id has been processed
		var finding oscalTypes.Finding
		var pass, fail int

		// A single finding should be "control-id centric"
		if _, ok := r.findingMap[requirement.ImplementedRequirement.ControlId]; ok {
			finding = r.findingMap[requirement.ImplementedRequirement.ControlId]
			finding.Description += fmt.Sprintf("Control Implementation: %s / Implemented Requirement: %s\n%s\n", requirement.ControlImplementation.UUID, requirement.ImplementedRequirement.UUID, requirement.ImplementedRequirement.Description)
		} else {
			finding = oscalTypes.Finding{
				UUID:        uuid.NewUUID(),
				Title:       fmt.Sprintf("Validation Result - Control: %s", requirement.ImplementedRequirement.ControlId),
				Description: fmt.Sprintf("Control Implementation: %s / Implemented Requirement: %s\n%s\n", requirement.ControlImplementation.UUID, requirement.ImplementedRequirement.UUID, requirement.ImplementedRequirement.Description),
			}
		}

		if requirement.ImplementedRequirement.Links != nil {
			relatedObservations := make([]oscalTypes.RelatedObservation, 0, len(*requirement.ImplementedRequirement.Links))
			for _, link := range *requirement.ImplementedRequirement.Links {
				observation, passBool := validationStore.GetRelatedObservation(link.Href)
				relatedObservations = append(relatedObservations, observation)
				if passBool {
					pass++
				} else {
					fail++
				}
			}
			// If there are pre-existing related observations we need to append
			if finding.RelatedObservations != nil {
				relatedObservations = append(relatedObservations, *finding.RelatedObservations...)
			}
			finding.RelatedObservations = &relatedObservations
		}

		// Using language from Assessment Results model for Target Objective Status State
		var state string
		message.Debugf("Pass: %v / Fail: %v / Existing State: %s", pass, fail, finding.Target.Status.State)
		if finding.Target.Status.State == "not-satisfied" {
			state = "not-satisfied"
		} else if pass > 0 && fail <= 0 {
			state = "satisfied"
		} else {
			state = "not-satisfied"
		}

		finding.Target = oscalTypes.FindingTarget{
			Status: oscalTypes.ObjectiveStatus{
				State: state,
			},
			TargetId: requirement.ImplementedRequirement.ControlId,
			Type:     "objective-id",
		}

		r.findingMap[requirement.ImplementedRequirement.ControlId] = finding
	}
	return r.findingMap
}

// GetStats returns the stats of the store
func (r *RequirementStore) GetStats(validationStore *validationstore.ValidationStore) Stats {
	var executableValidations bool
	var executableValidationsMsg string
	if validationStore != nil {
		executableValidations, executableValidationsMsg = validationStore.DryRun()
	}

	return Stats{
		TotalRequirements:        len(r.requirementMap),
		TotalValidations:         validationStore.Count(),
		ExecutableValidations:    executableValidations,
		ExecutableValidationsMsg: executableValidationsMsg,
		TotalFindings:            len(r.findingMap),
	}
}

// Drop unused validations from store (only relevant when store created from the backmatter if it contains unused validations)
// func (r *RequirementStore) DropUnusedValidations() {
// 	TODO
// }
