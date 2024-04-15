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

func ComponentFromCatalog(source string, catalog oscalTypes_1_1_2.Catalog, targetControls []string) (componentDefinition oscalTypes_1_1_2.ComponentDefinition, err error) {

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
		if controlArray, ok := controlMap[group.ID]; ok {
			for _, control := range *group.Controls {
				id := strings.Split(control.ID, "-")
				// If this is a control we have identified
				if contains(controlArray, id[1]) {
					newRequirement, err := ControlToImplementedRequirement(control)
					if err != nil {
						return componentDefinition, err
					}
					implmentedRequirements = append(implmentedRequirements, newRequirement)
				}

			}
		}

	}

	componentDefinition.Components = &[]oscalTypes_1_1_2.DefinedComponent{
		{
			UUID:  uuid.NewUUID(),
			Type:  "software",
			Title: "Software Title",
			ControlImplementations: &[]oscalTypes_1_1_2.ControlImplementationSet{
				{
					UUID:                    uuid.NewUUIDWithSource(source),
					Source:                  source,
					ImplementedRequirements: implmentedRequirements,
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
func ControlToImplementedRequirement(control oscalTypes_1_1_2.Control) (implementedRequirement oscalTypes_1_1_2.ImplementedRequirementControlImplementation, err error) {

	var controlDescription string
	paramMap := make(map[string]parameter)

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

	for _, part := range *control.Parts {
		// I feel like we need recursion here
		// let's just start with name == "statement" for now
		if part.Name == "statement" {
			if part.Parts != nil {
				controlDescription += addPart(part.Parts, paramMap, 0)
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
		if strings.Contains(prose, "{{ insert: param,") {
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
