package oscal

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"sigs.k8s.io/yaml"

	"github.com/defenseunicorns/lula/src/pkg/common"
)

type SystemSecurityPlan struct {
	Model *oscalTypes.SystemSecurityPlan
}

func NewSystemSecurityPlan() *SystemSecurityPlan {
	var systemSecurityPlan SystemSecurityPlan
	systemSecurityPlan.Model = nil
	return &systemSecurityPlan
}

func (ssp *SystemSecurityPlan) GetType() string {
	return "system-security-plan"
}

func (ssp *SystemSecurityPlan) GetCompleteModel() *oscalTypes.OscalModels {
	return &oscalTypes.OscalModels{
		SystemSecurityPlan: ssp.Model,
	}
}

// MakeDeterministic ensures the elements of the SSP are sorted deterministically
func (ssp *SystemSecurityPlan) MakeDeterministic() error {
	if ssp.Model == nil {
		return fmt.Errorf("cannot make nil model deterministic")
	} else {
		// Sort the SystemImplementation.Components by title
		slices.SortStableFunc(ssp.Model.SystemImplementation.Components, func(a, b oscalTypes.SystemComponent) int {
			return strings.Compare(a.Title, b.Title)
		})

		// Sort the ControlImplementation.ImplementedRequirements by control-id
		slices.SortStableFunc(ssp.Model.ControlImplementation.ImplementedRequirements, func(a, b oscalTypes.ImplementedRequirement) int {
			return CompareControlsInt(a.ControlId, b.ControlId)
		})

		// Sort the ControlImplementation.ImplementedRequirements.ByComponent by title
		for _, implementedRequirement := range ssp.Model.ControlImplementation.ImplementedRequirements {
			if implementedRequirement.ByComponents != nil {
				slices.SortStableFunc(*implementedRequirement.ByComponents, func(a, b oscalTypes.ByComponent) int {
					return strings.Compare(a.ComponentUuid, b.ComponentUuid)
				})
			}
		}

		// sort backmatter
		if ssp.Model.BackMatter != nil {
			backmatter := *ssp.Model.BackMatter
			sortBackMatter(&backmatter)
			ssp.Model.BackMatter = &backmatter
		}
	}
	return nil
}

// HandleExisting updates the existing SSP if a file is provided
func (ssp *SystemSecurityPlan) HandleExisting(path string) error {
	exists, err := common.CheckFileExists(path)
	if err != nil {
		return err
	}
	if exists {
		path = filepath.Clean(path)
		existingFileBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}
		ssp := NewSystemSecurityPlan()
		err = ssp.NewModel(existingFileBytes)
		if err != nil {
			return err
		}
		model, err := MergeSystemSecurityPlanModels(ssp.Model, ssp.Model)
		if err != nil {
			return err
		}
		ssp.Model = model
	}
	return nil
}

// NewModel updates the SSP model with the provided data
func (ssp *SystemSecurityPlan) NewModel(data []byte) error {
	var oscalModels oscalTypes.OscalModels

	err := multiModelValidate(data)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &oscalModels)
	if err != nil {
		return err
	}

	ssp.Model = oscalModels.SystemSecurityPlan

	return nil
}

// GenerateSystemSecurityPlan generates an OSCALModel System Security Plan.
// Command is the command that was used to generate the SSP.
// Source is the catalog source url that should be extracted from the component definition.
// Compdef is the partially* composed component definition and all merged component-definitions.
// TODOs: implement *partially = just imported component-definitions, remapped validation links;
// implement system-characteristics, parties->users->components, component status, (probably more);
// support for target instead of source?
func GenerateSystemSecurityPlan(command string, source string, compdef *oscalTypes.ComponentDefinition) (*SystemSecurityPlan, error) {
	if compdef == nil {
		return nil, fmt.Errorf("component definition is nil")
	}

	// Create the OSCAL SSP model for use and later assignment to the oscal.SystemSecurityPlan implementation
	var model oscalTypes.SystemSecurityPlan

	// Single time used for all time related fields
	rfc3339Time := time.Now()

	// Always create a new UUID for the assessment results (for now)
	model.UUID = uuid.NewUUID()

	// Creation of the generation prop
	props := []oscalTypes.Property{
		{
			Name:  "generation",
			Ns:    LULA_NAMESPACE,
			Value: command,
		},
		{
			Name:  "framework", // Should this be here? ** Assuming framework:source is 1:1
			Ns:    LULA_NAMESPACE,
			Value: source,
		},
	}

	// Create metadata object with requires fields and a few extras
	// Adding props to metadata as it is less available within the model
	model.Metadata = oscalTypes.Metadata{
		Title:        "System Security Plan",
		Version:      "0.0.1",
		OscalVersion: OSCAL_VERSION,
		Remarks:      "System Security Plan generated from Lula",
		Published:    &rfc3339Time,
		LastModified: rfc3339Time,
		Props:        &props,
		Parties:      compdef.Metadata.Parties, // TODO: Should these be handled on compdef merge?
	}

	// Update the import-profile
	model.ImportProfile = oscalTypes.ImportProfile{
		Href: source,
	}

	// Add system characteristics
	model.SystemCharacteristics = oscalTypes.SystemCharacteristics{
		SystemName: "Generated System",
		Status: oscalTypes.Status{
			State:   "operational", // Defaulting to operational, will need to revisit how this should be set
			Remarks: "<TODO: Validate state and remove this remark>",
		},
		SystemIds: []oscalTypes.SystemId{
			{
				ID: "generated-system",
			},
		},
		SystemInformation: oscalTypes.SystemInformation{
			InformationTypes: []oscalTypes.InformationType{
				{
					UUID:        uuid.NewUUID(),
					Title:       "Generated System Information",
					Description: "<TODO: Update information types>",
				},
			},
		},
	}

	// Decompose the component defn and add to the right sections of the SSP
	// TODO: external mapping of status? users? etc
	// only pull components from the selected source...
	implementedRequirementMap := CreateImplementedRequirementsByFramework(compdef)
	componentsMap := ComponentsToMap(compdef)

	if implementedRequirements, ok := implementedRequirementMap[source]; ok {
		// Update the control-implementation.implemented-requirements & system-implementation.components
		model.ControlImplementation = oscalTypes.ControlImplementation{
			ImplementedRequirements: make([]oscalTypes.ImplementedRequirement, 0),
		}
		model.SystemImplementation = oscalTypes.SystemImplementation{
			Components: make([]oscalTypes.SystemComponent, 0),
			Users: []oscalTypes.SystemUser{
				{
					UUID:    uuid.NewUUID(),
					Title:   "Generated User",
					Remarks: "<TODO: Update generated user>",
				},
			},
		}
		componentsAdded := make([]string, 0)

		for _, implementedRequirement := range implementedRequirements {
			// Append to the control-implementation.implemented-requirements
			model.ControlImplementation.ImplementedRequirements = append(model.ControlImplementation.ImplementedRequirements, implementedRequirement)

			// Append to the system-implementation.components
			for _, byComponent := range *implementedRequirement.ByComponents {
				if !slices.Contains(componentsAdded, byComponent.ComponentUuid) {
					if component, ok := componentsMap[byComponent.ComponentUuid]; ok {
						model.SystemImplementation.Components = append(model.SystemImplementation.Components, oscalTypes.SystemComponent{
							UUID:             component.UUID,
							Type:             component.Type,
							Title:            component.Title,
							Description:      component.Description,
							Props:            component.Props,
							Links:            component.Links,
							ResponsibleRoles: component.ResponsibleRoles,
							Protocols:        component.Protocols,
							Status: oscalTypes.SystemComponentStatus{
								State:   "operational", // Defaulting to operational, will need to revisit how this should be set
								Remarks: "<TODO: Validate state and remove this remark>",
							},
						})
					}

					componentsAdded = append(componentsAdded, byComponent.ComponentUuid)
				}
			}
		}
	} else {
		return nil, fmt.Errorf("no implemented requirements found for source %s", source)
	}

	return &SystemSecurityPlan{
		Model: &model,
	}, nil

}

// MergeSystemSecurityPlanModels merges two SystemSecurityPlan models
// Requires that the source of the models are the same
func MergeSystemSecurityPlanModels(original *oscalTypes.SystemSecurityPlan, latest *oscalTypes.SystemSecurityPlan) (*oscalTypes.SystemSecurityPlan, error) {
	// Input nil checks
	if original == nil && latest != nil {
		return latest, nil
	} else if original != nil && latest == nil {
		return original, nil
	} else if original == nil && latest == nil {
		return nil, fmt.Errorf("both models are nil")
	}

	// Check that the sources are the same, if not then can't be merged
	if original.ImportProfile.Href != latest.ImportProfile.Href {
		return nil, fmt.Errorf("cannot merge models with different sources")
	}

	// Merge unique Components in the SystemImplementation
	original.SystemImplementation.Components = mergeSystemComponents(original.SystemImplementation.Components, latest.SystemImplementation.Components)

	// Merge unique ImplementedRequirements in the ControlImplementation
	original.ControlImplementation.ImplementedRequirements = mergeImplementedRequirements(original.ControlImplementation.ImplementedRequirements, latest.ControlImplementation.ImplementedRequirements)

	// Merge the back-matter resources
	if original.BackMatter != nil && latest.BackMatter != nil {
		original.BackMatter = &oscalTypes.BackMatter{
			Resources: mergeResources(original.BackMatter.Resources, latest.BackMatter.Resources),
		}
	} else if original.BackMatter == nil && latest.BackMatter != nil {
		original.BackMatter = latest.BackMatter
	}

	// Update the uuid
	original.UUID = uuid.NewUUID()

	return original, nil
}

type ImplementedRequirementMap map[string]oscalTypes.ImplementedRequirement

// CreateImplementedRequirementsByFramework sorts the implemented requirements for each framework
func CreateImplementedRequirementsByFramework(compdef *oscalTypes.ComponentDefinition) map[string]ImplementedRequirementMap {
	frameworkImplementedRequirementsMap := make(map[string]ImplementedRequirementMap)

	// Sort components by framework
	if compdef != nil && compdef.Components != nil {
		for _, component := range *compdef.Components {
			if component.ControlImplementations != nil {
				for _, controlImplementation := range *component.ControlImplementations {
					// update list of frameworks in a given control-implementation
					frameworks := []string{controlImplementation.Source}
					status, value := GetProp("framework", LULA_NAMESPACE, controlImplementation.Props)
					if status {
						frameworks = append(frameworks, value)
					}

					for _, framework := range frameworks {
						// Initialize the map for the source and framework if it doesn't exist
						_, ok := frameworkImplementedRequirementsMap[framework]
						if !ok {
							frameworkImplementedRequirementsMap[framework] = make(map[string]oscalTypes.ImplementedRequirement)
						}

						// For each implemented requirement, add it to the map
						for _, implementedRequirement := range controlImplementation.ImplementedRequirements {
							existingIr, ok := frameworkImplementedRequirementsMap[framework][implementedRequirement.ControlId]
							if ok {
								// If found, update the existing implemented requirement
								// TODO: add other "ByComponents" fields?
								*existingIr.ByComponents = append(*existingIr.ByComponents, oscalTypes.ByComponent{
									ComponentUuid: component.UUID,
									UUID:          uuid.NewUUID(),
									Description:   implementedRequirement.Description,
									Links:         implementedRequirement.Links,
								})
							} else {
								// Otherwise create a new implemented-requirement
								frameworkImplementedRequirementsMap[framework][implementedRequirement.ControlId] = oscalTypes.ImplementedRequirement{
									UUID:      uuid.NewUUID(),
									ControlId: implementedRequirement.ControlId,
									Remarks:   implementedRequirement.Remarks,
									ByComponents: &[]oscalTypes.ByComponent{
										{
											ComponentUuid: component.UUID,
											UUID:          uuid.NewUUID(),
											Description:   implementedRequirement.Description,
											Links:         implementedRequirement.Links,
										},
									},
								}
							}
						}
					}
				}
			}
		}
	}
	return frameworkImplementedRequirementsMap
}

func mergeSystemComponents(original []oscalTypes.SystemComponent, latest []oscalTypes.SystemComponent) []oscalTypes.SystemComponent {
	// Check all latest, add to original if not present
	for _, latestComponent := range latest {
		found := false
		for _, originalComponent := range original {
			if latestComponent.UUID == originalComponent.UUID {
				found = true
				break
			}
		}
		//if not found, append
		if !found {
			original = append(original, latestComponent)
		}
	}
	return original
}

func mergeImplementedRequirements(original []oscalTypes.ImplementedRequirement, latest []oscalTypes.ImplementedRequirement) []oscalTypes.ImplementedRequirement {
	for _, latestRequirement := range latest {
		found := false
		for _, originalRequirement := range original {
			if latestRequirement.ControlId == originalRequirement.ControlId {
				found = true
				// Update ByComponent
				for _, latestByComponent := range *latestRequirement.ByComponents {
					foundByComponent := false
					// Latest component is already in original, do nothing
					// ** Assumption: There should never be a different Component reference specification to the same control, e.g., different links to append
					for _, originalByComponent := range *originalRequirement.ByComponents {
						if latestByComponent.UUID == originalByComponent.UUID {
							foundByComponent = true
							break
						}
					}
					//if not found, append
					if !foundByComponent {
						*originalRequirement.ByComponents = append(*originalRequirement.ByComponents, latestByComponent)
					}
				}
				break
			}
		}
		//if not found, append
		if !found {
			original = append(original, latestRequirement)
		}
	}
	return original
}
