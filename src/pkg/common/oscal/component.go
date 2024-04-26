package oscal

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/message"

	"sigs.k8s.io/yaml"
)

type selection struct {
	HowMany string
	Choice  []string
}

type parameter struct {
	ID     string
	Label  string
	Select *selection
}

// NewOscalComponentDefinition consumes a byte array and returns a new single OscalComponentDefinitionModel object
// Standard use is to read a file from the filesystem and pass the []byte to this function
func NewOscalComponentDefinition(source string, data []byte) (componentDefinition oscalTypes_1_1_2.ComponentDefinition, err error) {
	var oscalModels oscalTypes_1_1_2.OscalModels

	if strings.HasSuffix(source, ".yaml") {
		err = yaml.Unmarshal(data, &oscalModels)
		if err != nil {
			message.Debugf("Error marshalling yaml: %s\n", err.Error())
			return componentDefinition, err
		}
	} else if strings.HasSuffix(source, ".json") {
		err = json.Unmarshal(data, &oscalModels)
		if err != nil {
			message.Debugf("Error marshalling json: %s\n", err.Error())
			return componentDefinition, err
		}
	} else {
		return componentDefinition, fmt.Errorf("unsupported file type: %s", source)
	}

	return *oscalModels.ComponentDefinition, nil
}

// This function should perform a merge of a component-definition with a specific component/control-implementation combo where maintaining the original component-definition is the primary concern.
func MergeComponentDefinitionOnComponent(original oscalTypes_1_1_2.ComponentDefinition, newComponent oscalTypes_1_1_2.DefinedComponent, newControlImplementation oscalTypes_1_1_2.ControlImplementationSet) (oscalTypes_1_1_2.ComponentDefinition, error) {
	var found bool = false
	// Given I have the component title and control-implementation source - How will I replace a targeted control-implementation in a targeted component?

	// Step 1 - identify the component - if it exists
	// If it doesn't exist in the original - add and return
	// So we create a temporary slice of components to re-assign?
	tempComponents := make([]oscalTypes_1_1_2.DefinedComponent, 0)
	var targetComponent oscalTypes_1_1_2.DefinedComponent
	for _, component := range *original.Components {
		if component.Title == newComponent.Title {
			targetComponent = component
			found = true
		} else {
			tempComponents = append(tempComponents, component)
		}
	}

	// New component was not found - append and return
	if !found {
		tempComponents = append(tempComponents, newComponent)
		original.Components = &tempComponents
		return original, nil
	}

	// reset found bool
	found = false

	// New Component was found - continue with merge
	// Step 2 - identify the control-implementation - if it exists
	// If it doesn't exist in the original - add and return
	tempControlImplementations := make([]oscalTypes_1_1_2.ControlImplementationSet, 0)
	var targetControlImplementation oscalTypes_1_1_2.ControlImplementationSet
	for _, controlImplementation := range *targetComponent.ControlImplementations {
		if controlImplementation.Source == newControlImplementation.Source {
			targetControlImplementation = controlImplementation
			found = true
		} else {
			tempControlImplementations = append(tempControlImplementations, controlImplementation)
		}
	}

	if !found {
		tempControlImplementations = append(tempControlImplementations, newControlImplementation)
		targetComponent.ControlImplementations = &tempControlImplementations
		tempComponents = append(tempComponents, targetComponent)
		original.Components = &tempComponents
		return original, nil
	}
	// Step 3 - Update or Replace - do not remove
	newImpMap := getImplementedRequirementsMap(newControlImplementation.ImplementedRequirements)
	originalImpMap := getImplementedRequirementsMap(targetControlImplementation.ImplementedRequirements)

	// Create a placeholder for implemented-requirements
	implmentedRequirements := make([]oscalTypes_1_1_2.ImplementedRequirementControlImplementation, 0)
	for id, req := range newImpMap {
		if orig, ok := originalImpMap[id]; ok {
			// requirement exists in both - update remarks
			orig.Remarks = req.Remarks
			implmentedRequirements = append(implmentedRequirements, orig)
			delete(originalImpMap, id)
		} else {
			implmentedRequirements = append(implmentedRequirements, req)
		}
	}
	// Add the remaining original implemented-requirements
	for _, req := range originalImpMap {
		implmentedRequirements = append(implmentedRequirements, req)
	}
	targetControlImplementation.ImplementedRequirements = implmentedRequirements
	tempControlImplementations = append(tempControlImplementations, targetControlImplementation)
	targetComponent.ControlImplementations = &tempControlImplementations
	tempComponents = append(tempComponents, targetComponent)
	original.Components = &tempComponents
	return original, nil

	// TODO: maybe this generation is constrained to a specific control-implementation?
	// Meaning we generate a latest component definition with the control-implementation specified in catalog-source
	// Then check for its existence in the original component-definition - if not exists - add new control-implementation and return? otherwise if exists - continue with merge

	// How to check the delta between two component-definitions with []oscalTypes_1_1_2.ImplementedRequirementControlImplementation ?

	// Get a map[control-id]oscalTypes_1_1_2.ImplementedRequirementControlImplementation from each component-definition

	// originalMap, err := getImplementedRequirementsMap(original)
	// latestMap, err := getImplementedRequirementsMap(latest)

	// Create a new []oscalTypes_1_1_2.ImplementedRequirementControlImplementation var
	// implmentedRequirements := make([]oscalTypes_1_1_2.ImplementedRequirementControlImplementation, 0)

	// For each key in the latest map - check if it exists in the original map
	// If so - add the original to implementedRequirements else add the latest

	// reassign the control-implementations.implemented-requirements array

	// return a copy of the original artifact

	// return componentDefinition, err
}

// Creates a component-definition from a catalog and identified (or all) controls. Allows for specification of what the content of the remarks section should contain.
func ComponentFromCatalog(source string, catalog oscalTypes_1_1_2.Catalog, componentTitle string, targetControls []string, targetRemarks []string, allControls bool) (componentDefinition oscalTypes_1_1_2.ComponentDefinition, err error) {

	message.Debugf("target controls %v", targetControls)
	// store all of the implemented requirements
	implmentedRequirements := make([]oscalTypes_1_1_2.ImplementedRequirementControlImplementation, 0)

	// How will we identify a control based on the array of controls passed?
	controlMap := make(map[string][]string)
	for _, control := range targetControls {
		id := strings.Split(control, "-")
		controlMap[id[0]] = append(controlMap[id[0]], id[1])
	}

	// We then want to iterate through the catalog, identify the controls and map control information to the implemented-requirements

	for _, group := range *catalog.Groups {
		// Is this a group of controls that we are targeting
		if controlArray, ok := controlMap[group.ID]; ok || allControls {
			for _, control := range *group.Controls {
				id := strings.Split(control.ID, "-")
				// Check if the control is the primary control
				if contains(controlArray, id[1]) || allControls {
					newRequirement, err := ControlToImplementedRequirement(control, targetRemarks)
					if err != nil {
						return componentDefinition, err
					}
					implmentedRequirements = append(implmentedRequirements, newRequirement)
				}
				// Now check if this control has sub-controls - this can likely be improved by checking for the '.' notation in a control
				// We need
				if control.Controls != nil {
					for _, subControl := range *control.Controls {
						subId := strings.Split(subControl.ID, "-")
						if contains(controlArray, subId[1]) || allControls {
							newRequirement, err := ControlToImplementedRequirement(subControl, targetRemarks)
							if err != nil {
								return componentDefinition, err
							}
							implmentedRequirements = append(implmentedRequirements, newRequirement)
						}
					}
				}
			}
		}

	}

	componentDefinition.Components = &[]oscalTypes_1_1_2.DefinedComponent{
		{
			UUID:        uuid.NewUUID(),
			Type:        "software",
			Title:       componentTitle,
			Description: "Component Description",
			ControlImplementations: &[]oscalTypes_1_1_2.ControlImplementationSet{
				{
					UUID:                    uuid.NewUUIDWithSource(source),
					Source:                  source,
					ImplementedRequirements: implmentedRequirements,
					Description:             "Control Implementation Description",
				},
			},
		},
	}
	rfc3339Time := time.Now()

	componentDefinition.UUID = uuid.NewUUID()

	componentDefinition.Metadata = oscalTypes_1_1_2.Metadata{
		OscalVersion: OSCAL_VERSION,
		LastModified: rfc3339Time,
		Published:    &rfc3339Time,
		Remarks:      "Lula Generated Component Definition",
		Title:        "Component Title",
		Version:      "0.0.1",
	}

	return componentDefinition, nil

}

// Consume a control - Identify statements - iterate through parts in order to create a description
func ControlToImplementedRequirement(control oscalTypes_1_1_2.Control, targetRemarks []string) (implementedRequirement oscalTypes_1_1_2.ImplementedRequirementControlImplementation, err error) {

	var controlDescription string
	paramMap := make(map[string]parameter)

	if control.Params != nil {
		for _, param := range *control.Params {

			if param.Select == nil {
				paramMap[param.ID] = parameter{
					ID:    param.ID,
					Label: param.Label,
				}
			} else {
				sel := *param.Select
				paramMap[param.ID] = parameter{
					ID: param.ID,
					Select: &selection{
						HowMany: sel.HowMany,
						Choice:  *sel.Choice,
					},
				}
			}
		}
	} else {
		message.Debugf("No parameters found for %s", control.ID)
	}

	if control.Parts != nil {
		for _, part := range *control.Parts {
			if contains(targetRemarks, part.Name) {
				controlDescription += fmt.Sprintf("%s:\n", strings.ToTitle(part.Name))
				if part.Prose != "" && strings.Contains(part.Prose, "{{ insert: param,") {
					controlDescription += replaceParams(part.Prose, paramMap)
				} else {
					controlDescription += part.Prose
				}
				if part.Parts != nil {
					controlDescription += addPart(part.Parts, paramMap, 0)
				}
			}
		}
	}

	// assemble implemented-requirements object
	implementedRequirement.Remarks = controlDescription
	implementedRequirement.Description = "<how the specified control may be implemented if the containing component or capability is instantiated in a system security plan>"
	implementedRequirement.ControlId = control.ID
	implementedRequirement.UUID = uuid.NewUUID()

	return implementedRequirement, nil
}

// Returns a map of the uuid - description of the back-matter resources
func BackMatterToMap(backMatter oscalTypes_1_1_2.BackMatter) (resourceMap map[string]string) {
	if backMatter.Resources == nil {
		return nil
	}

	resourceMap = make(map[string]string)
	for _, resource := range *backMatter.Resources {
		// perform a check to see if the key already exists (meaning duplicitive uuid use)
		_, exists := resourceMap[resource.UUID]
		if exists {
			message.Warnf("Duplicative UUID use detected - Overwriting UUID %s", resource.UUID)
		}

		resourceMap[resource.UUID] = resource.Description
	}
	return resourceMap

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Function to allow for recursively adding prose to the description string
func addPart(part *[]oscalTypes_1_1_2.Part, paramMap map[string]parameter, level int) string {

	var result, label string

	for _, part := range *part {
		// need to get the label first - unsure if there will ever be more than one?
		for _, prop := range *part.Props {
			if prop.Name == "label" {
				label = prop.Value
			}
		}
		var tabs string
		for range level {
			tabs += "\t"
		}
		prose := part.Prose
		if prose == "" {
			result += fmt.Sprintf("%s%s\n", tabs, label)
		} else if strings.Contains(prose, "{{ insert: param,") {
			result += fmt.Sprintf("%s%s %s\n", tabs, label, replaceParams(prose, paramMap))
		} else {
			result += fmt.Sprintf("%s%s %s\n", tabs, label, prose)
		}
		if part.Parts != nil {
			result += addPart(part.Parts, paramMap, level+1)
		}

	}

	return result
}

func replaceParams(input string, params map[string]parameter) string {
	re := regexp.MustCompile(`{{\s*insert:\s*param,\s*([^}\s]+)\s*}}`)
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		paramName := strings.TrimSpace(re.FindStringSubmatch(match)[1])
		if param, ok := params[paramName]; ok {
			if param.Select == nil {
				return fmt.Sprintf("[Assignment: organization-defined %s]", param.Label)
			} else {
				return fmt.Sprintf("[Selection: (%s) organization-defined %s]", param.Select.HowMany, strings.Join(param.Select.Choice, "; "))
			}
		}
		return match
	})
	return result
}

func getImplementedRequirementsMap(irs []oscalTypes_1_1_2.ImplementedRequirementControlImplementation) map[string]oscalTypes_1_1_2.ImplementedRequirementControlImplementation {
	result := make(map[string]oscalTypes_1_1_2.ImplementedRequirementControlImplementation)
	for _, ir := range irs {
		result[ir.ControlId] = ir
	}
	return result
}
