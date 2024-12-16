package oscal

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
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
// Source is the profile source url that should be used to pull implemented-requirements from the component definition. "Target" not currently supported.
// Profile is the profile model that should be used to populate the SSP.
// Compdefs are all component definitions that should be merged into the SSP.
// This will return an error if the profile does not contain any controls.
func GenerateSystemSecurityPlan(command, source string, targetRemarks []string, profile *oscalTypes.Profile, compdefs ...*oscalTypes.ComponentDefinition) (*SystemSecurityPlan, error) {
	compdef, err := MergeVariadicComponentDefinition(compdefs...)
	if err != nil {
		return nil, err
	}
	componentsMap := ComponentsToMap(compdef)

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
	}

	// Update parties from component definition if not nil
	// TODO: Handle parties on component definition merge op
	if compdef != nil {
		model.Metadata.Parties = compdef.Metadata.Parties
	}

	// Update the import-profile
	model.ImportProfile = oscalTypes.ImportProfile{
		Href: source,
	}

	// Add placeholder system characteristics
	model.SystemCharacteristics = oscalTypes.SystemCharacteristics{
		SystemName: "Generated System",
		Status: oscalTypes.Status{
			State:   "operational", // Defaulting to operational, will need to revisit how this should be set
			Remarks: "TODO: Validate state and remove this remark",
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
					Description: "TODO: Update information types",
				},
			},
		},
	}

	// Add placeholder system-implementation
	model.SystemImplementation = oscalTypes.SystemImplementation{
		Components: []oscalTypes.SystemComponent{
			{
				UUID:    uuid.NewUUID(),
				Title:   "Generated Component",
				Remarks: "TODO: Update generated component",
				Type:    "software",
				Status: oscalTypes.SystemComponentStatus{
					State:   "operational", // Defaulting to operational, will need to revisit how this should be set
					Remarks: "TODO: Validate state and remove this remark",
				},
			},
		},
		Users: []oscalTypes.SystemUser{
			{
				UUID:    uuid.NewUUID(),
				Title:   "Generated User",
				Remarks: "TODO: Update generated user",
			},
		},
	}

	// Get all source/target-mapped control-ids -> by-components
	componentControlMap := CreateSourceControlsMap(compdef)

	// This will most likely remove all "target" keys, since those will not be resolvable links
	componentControlMapUUIDs := RemapSourceToUUID(componentControlMap)

	// Get all controls from profile -> implemented-requirements
	profileControlMap, err := ResolveProfileControls(profile, source, "", nil, nil)
	if err != nil {
		return nil, err
	}

	// Get all implemented requirements
	controlMap, ok := profileControlMap[profile.UUID]
	if !ok {
		return nil, fmt.Errorf("profile %s not found", profile.UUID)
	}

	// Return an error if the profile does not contain any controls (or should I add a placeholder?)
	if len(controlMap) == 0 {
		return nil, fmt.Errorf("profile %s does not contain any controls", profile.UUID)
	}

	model.ControlImplementation = oscalTypes.ControlImplementation{
		ImplementedRequirements: make([]oscalTypes.ImplementedRequirement, 0, len(controlMap)),
	}

	componentsUsed := make([]string, 0)
	// Update the control-implementation.implemented-requirements
	for id, control := range controlMap {
		// Create an implemented-requirement
		implementedRequirement, err := createImplementedRequirement(&control, targetRemarks)
		if err != nil {
			return nil, err
		}

		// Check if any source UUIDs in componentControlMapUUIDs are in the profileControlMap -> update the respective implemented-requirement
		for u, componentControls := range componentControlMapUUIDs {
			// Check if the UUID is in the profileControlMap
			if _, ok := profileControlMap[u]; ok {
				// Check if the control-id is linked to the component
				if _, ok := componentControls[id]; ok {
					// Update the implemented-requirement.by-components
					for _, component := range componentControls[id] {
						if implementedRequirement.ByComponents == nil {
							implementedRequirement.ByComponents = new([]oscalTypes.ByComponent)
						}
						*implementedRequirement.ByComponents = append(*implementedRequirement.ByComponents, component)
						if !slices.Contains(componentsUsed, component.ComponentUuid) {
							// Add new components only
							componentsUsed = append(componentsUsed, component.ComponentUuid)
						}
					}
				}
			}
		}

		// Add the implemented-requirement to the SSP
		model.ControlImplementation.ImplementedRequirements = append(model.ControlImplementation.ImplementedRequirements, implementedRequirement)
	}

	// Update the system-implementation.components with any components used in the implemented-requirements
	if len(componentsUsed) > 0 {
		actualComponents := make([]oscalTypes.SystemComponent, 0, len(componentsUsed))
		for _, cUUID := range componentsUsed {
			component, ok := componentsMap[cUUID]
			if !ok {
				continue
			}
			actualComponents = append(actualComponents, oscalTypes.SystemComponent{
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
					Remarks: "TODO: Validate state and remove this remark",
				},
			})
		}
		model.SystemImplementation.Components = actualComponents
	}

	ssp := &SystemSecurityPlan{
		Model: &model,
	}

	// TODO: Perform a model merge of anything in oscal-parameters (ssp.MergeConfigurationData) overwrites anything as a "placeholder" or from the component-definition

	return ssp, nil

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

// Get Components -> ImplementedRequirements map[string]ComponentsIRs

type ByComponentsMap map[string][]oscalTypes.ByComponent

// CreateSourceControlsMap maps the source/framework -> control-id -> []by-component
func CreateSourceControlsMap(compdef *oscalTypes.ComponentDefinition) map[string]ByComponentsMap {
	// Map that groups all sources/frameworks -> control-id -> by-component
	sourceControlsMap := make(map[string]ByComponentsMap)

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
						_, ok := sourceControlsMap[framework]
						if !ok {
							sourceControlsMap[framework] = make(map[string][]oscalTypes.ByComponent)
						}

						// For each implemented requirement, add it to the map
						for _, implementedRequirement := range controlImplementation.ImplementedRequirements {
							_, ok := sourceControlsMap[framework][implementedRequirement.ControlId]
							if ok {
								// If found, update the existing implemented requirement
								// TODO: add other "ByComponents" fields?
								sourceControlsMap[framework][implementedRequirement.ControlId] = append(sourceControlsMap[framework][implementedRequirement.ControlId], oscalTypes.ByComponent{
									ComponentUuid: component.UUID,
									UUID:          uuid.NewUUID(),
									Description:   implementedRequirement.Description,
									Links:         implementedRequirement.Links,
								})
							} else {
								// Otherwise create a new by-components slice
								sourceControlsMap[framework][implementedRequirement.ControlId] = []oscalTypes.ByComponent{
									{
										ComponentUuid: component.UUID,
										UUID:          uuid.NewUUID(),
										Description:   implementedRequirement.Description,
										Links:         implementedRequirement.Links,
									},
								}
							}
						}
					}
				}
			}
		}
	}

	return sourceControlsMap
}

// RemapSourceToUUID takes a map[string]any and remaps any source keys (profiles and catalogs) to their UUID
// NOTE: Doesn't support backmatter resources as source keys, only network resolvable links
func RemapSourceToUUID[V any](inMap map[string]V) map[string]V {
	outMap := make(map[string]V)

	for source, v := range inMap {
		// Fetch OSCAL Model
		// TODO: how to handle local file paths?
		oscalModel, modelType, err := FetchOSCALModel(source, "")
		if err != nil {
			// If unable to fetch, skip
			continue
		}

		switch modelType {
		case "profile":
			profile := oscalModel.Profile
			outMap[profile.UUID] = v
		case "catalog":
			catalog := oscalModel.Catalog
			outMap[catalog.UUID] = v
		}
	}
	return outMap
}

func createImplementedRequirement(control *oscalTypes.Control, targetRemarks []string) (implementedRequirement oscalTypes.ImplementedRequirement, err error) {
	remarks, err := getControlRemarks(control, targetRemarks)
	if err != nil {
		return implementedRequirement, err
	}

	// Create the implemented-requirement
	return oscalTypes.ImplementedRequirement{
		UUID:      uuid.NewUUID(),
		ControlId: control.ID,
		Remarks:   remarks,
	}, nil
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
		for oIdx, originalRequirement := range original {
			if latestRequirement.ControlId == originalRequirement.ControlId {
				found = true
				// Update ByComponent
				// If no by-components in latest, skip
				if latestRequirement.ByComponents != nil && originalRequirement.ByComponents != nil {
					for _, latestByComponent := range *latestRequirement.ByComponents {
						foundByComponent := false
						// Latest component is already in original, do nothing
						// ** Assumption: There should never be a different Component reference specification to the same control, e.g., different links to append
						for _, originalByComponent := range *originalRequirement.ByComponents {
							if latestByComponent.ComponentUuid == originalByComponent.ComponentUuid {
								foundByComponent = true
								break
							}
						}
						//if not found, append
						if !foundByComponent {
							*originalRequirement.ByComponents = append(*originalRequirement.ByComponents, latestByComponent)
						}
					}
				} else if latestRequirement.ByComponents != nil && originalRequirement.ByComponents == nil {
					original[oIdx].ByComponents = latestRequirement.ByComponents
				}
				break // Break when latest reqt is found in original
			}
		}
		//if not found, append
		if !found {
			original = append(original, latestRequirement)
		}
	}
	return original
}
