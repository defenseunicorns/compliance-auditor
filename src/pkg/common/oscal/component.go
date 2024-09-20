package oscal

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/message"

	"sigs.k8s.io/yaml"
)

type Requirement struct {
	ImplementedRequirement *oscalTypes_1_1_2.ImplementedRequirementControlImplementation
	ControlImplementation  *oscalTypes_1_1_2.ControlImplementationSet
}

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
func NewOscalComponentDefinition(data []byte) (componentDefinition *oscalTypes_1_1_2.ComponentDefinition, err error) {
	var oscalModels oscalTypes_1_1_2.OscalModels

	// validate the data
	err = multiModelValidate(data)
	if err != nil {
		return componentDefinition, err
	}

	err = yaml.Unmarshal(data, &oscalModels)
	if err != nil {
		return componentDefinition, err
	}
	return oscalModels.ComponentDefinition, nil
}

// This function should perform a merge of two component-definitions where maintaining the original component-definition is the primary concern.
func MergeComponentDefinitions(original *oscalTypes_1_1_2.ComponentDefinition, latest *oscalTypes_1_1_2.ComponentDefinition) (*oscalTypes_1_1_2.ComponentDefinition, error) {

	originalMap := make(map[string]oscalTypes_1_1_2.DefinedComponent)

	if original.Components == nil {
		return original, fmt.Errorf("original component-definition is nil")
	}

	if latest.Components == nil {
		return original, fmt.Errorf("latest component-definition is nil")
	}

	for _, component := range *original.Components {
		originalMap[component.Title] = component
	}

	latestMap := make(map[string]oscalTypes_1_1_2.DefinedComponent)

	for _, component := range *latest.Components {
		latestMap[component.Title] = component
	}

	tempItems := make([]oscalTypes_1_1_2.DefinedComponent, 0)
	for key, value := range latestMap {
		if comp, ok := originalMap[key]; ok {
			// if the component exists - merge & append
			comp = *mergeComponents(&comp, &value)
			tempItems = append(tempItems, comp)
			delete(originalMap, key)
		} else {
			// append the component
			tempItems = append(tempItems, value)
		}
	}

	for _, item := range originalMap {
		tempItems = append(tempItems, item)
	}

	// merge the back-matter resources
	if original.BackMatter != nil && latest.BackMatter != nil {
		original.BackMatter = &oscalTypes_1_1_2.BackMatter{
			Resources: mergeResources(original.BackMatter.Resources, latest.BackMatter.Resources),
		}
	} else if original.BackMatter == nil && latest.BackMatter != nil {
		original.BackMatter = latest.BackMatter
	}

	original.Components = &tempItems
	original.Metadata.LastModified = time.Now()

	// Artifact will be modified - need to update the UUID
	original.UUID = uuid.NewUUID()

	return original, nil

}

func mergeComponents(original *oscalTypes_1_1_2.DefinedComponent, latest *oscalTypes_1_1_2.DefinedComponent) *oscalTypes_1_1_2.DefinedComponent {
	originalMap := make(map[string]oscalTypes_1_1_2.ControlImplementationSet)

	if original.ControlImplementations != nil {
		for _, item := range *original.ControlImplementations {
			originalMap[item.Source] = item
		}
	}

	latestMap := make(map[string]oscalTypes_1_1_2.ControlImplementationSet)

	if latest.ControlImplementations != nil {
		for _, item := range *latest.ControlImplementations {
			latestMap[item.Source] = item
		}
	}

	tempItems := make([]oscalTypes_1_1_2.ControlImplementationSet, 0)
	for key, value := range latestMap {
		if orig, ok := originalMap[key]; ok {
			// if the control implementation exists - merge & append
			orig = *mergeControlImplementations(&orig, &value)
			tempItems = append(tempItems, orig)
			delete(originalMap, key)
		} else {
			// append the component
			tempItems = append(tempItems, value)
		}
	}

	for _, item := range originalMap {
		tempItems = append(tempItems, item)
	}

	// merge Props without duplicating
	// note: this assumes uniqueness of named prop
	if latest.Props != nil {
		if original.Props == nil {
			original.Props = latest.Props
		} else {
			for _, prop := range *latest.Props {
				UpdateProps(prop.Name, prop.Ns, prop.Value, original.Props)
			}
		}
	}

	original.ControlImplementations = &tempItems
	return original
}

func mergeControlImplementations(original *oscalTypes_1_1_2.ControlImplementationSet, latest *oscalTypes_1_1_2.ControlImplementationSet) *oscalTypes_1_1_2.ControlImplementationSet {
	originalMap := make(map[string]oscalTypes_1_1_2.ImplementedRequirementControlImplementation)

	if original.ImplementedRequirements != nil {
		for _, item := range original.ImplementedRequirements {
			originalMap[item.ControlId] = item
		}
	}
	latestMap := make(map[string]oscalTypes_1_1_2.ImplementedRequirementControlImplementation)

	if latest.ImplementedRequirements != nil {
		for _, item := range latest.ImplementedRequirements {
			latestMap[item.ControlId] = item
		}
	}

	tempItems := make([]oscalTypes_1_1_2.ImplementedRequirementControlImplementation, 0)
	for key, latestImp := range latestMap {
		if orig, ok := originalMap[key]; ok {
			// requirement exists in both - update remarks as this is solely owned by the automation
			orig.Remarks = latestImp.Remarks
			// update the links as another critical field
			if orig.Links != nil && latestImp.Links != nil {
				orig.Links = mergeLinks(*orig.Links, *latestImp.Links)
			} else if orig.Links == nil && latestImp.Links != nil {
				orig.Links = latest.Links
			}

			tempItems = append(tempItems, orig)
			delete(originalMap, key)
		} else {
			// append the component
			tempItems = append(tempItems, latestImp)
		}
	}

	for _, item := range originalMap {
		tempItems = append(tempItems, item)
	}
	// merge Props without duplicating
	// note: this assumes uniqueness of named prop
	if latest.Props != nil {
		if original.Props == nil {
			original.Props = latest.Props
		} else {
			for _, prop := range *latest.Props {
				UpdateProps(prop.Name, prop.Ns, prop.Value, original.Props)
			}
		}
	}
	original.ImplementedRequirements = tempItems
	return original
}

// Merges two arrays of resources into a single array
func mergeResources(orig *[]oscalTypes_1_1_2.Resource, latest *[]oscalTypes_1_1_2.Resource) *[]oscalTypes_1_1_2.Resource {
	if orig == nil {
		return latest
	}

	if latest == nil {
		return orig
	}

	result := make([]oscalTypes_1_1_2.Resource, 0)

	tempResource := make(map[string]oscalTypes_1_1_2.Resource)
	for _, resource := range *orig {
		tempResource[resource.UUID] = resource
		result = append(result, resource)
	}

	for _, resource := range *latest {
		// Only append if does not exist
		if _, ok := tempResource[resource.UUID]; !ok {
			result = append(result, resource)
		}
	}

	return &result
}

// Merges two arrays of links into a single array
// TODO: account for overriding validations
func mergeLinks(orig []oscalTypes_1_1_2.Link, latest []oscalTypes_1_1_2.Link) *[]oscalTypes_1_1_2.Link {
	result := make([]oscalTypes_1_1_2.Link, 0)

	tempLinks := make(map[string]oscalTypes_1_1_2.Link)
	for _, link := range orig {
		// Both of these are string fields, href is required - resource fragment can help establish uniqueness
		key := fmt.Sprintf("%s%s", link.Href, link.ResourceFragment)
		tempLinks[key] = link
		result = append(result, link)
	}

	for _, link := range latest {
		key := fmt.Sprintf("%s%s", link.Href, link.ResourceFragment)
		// Only append if does not exist
		if _, ok := tempLinks[key]; !ok {
			result = append(result, link)
		}
	}

	return &result
}

// Creates a component-definition from a catalog and identified (or all) controls. Allows for specification of what the content of the remarks section should contain.
func ComponentFromCatalog(command string, source string, catalog *oscalTypes_1_1_2.Catalog, componentTitle string, targetControls []string, targetRemarks []string, framework string) (*oscalTypes_1_1_2.ComponentDefinition, error) {
	// store all of the implemented requirements
	implementedRequirements := make([]oscalTypes_1_1_2.ImplementedRequirementControlImplementation, 0)
	var componentDefinition = &oscalTypes_1_1_2.ComponentDefinition{}

	if len(targetControls) == 0 {
		return componentDefinition, fmt.Errorf("no controls identified for generation")
	}

	controlMap := make(map[string]bool)
	for _, control := range targetControls {
		controlMap[control] = false
	}

	// A catalog has groups and controls
	// A group has controls and groups (note the nesting of groups)
	// A control has controls (note the nesting)
	// Given the nesting of groups/controls we will need to recursively search groups/controls

	// Begin Recursive group search
	if catalog.Groups != nil {
		newReqs, err := searchGroups(catalog.Groups, controlMap, targetRemarks)
		if err != nil {
			return componentDefinition, err
		}
		implementedRequirements = append(implementedRequirements, newReqs...)
	}

	// Begin recursive control search
	if catalog.Controls != nil {
		newReqs, err := searchControls(catalog.Controls, controlMap, targetRemarks)
		if err != nil {
			return componentDefinition, err
		}
		implementedRequirements = append(implementedRequirements, newReqs...)
	}

	// TODO: rework this - a catalog does not require groups
	if catalog.Groups == nil {
		return componentDefinition, fmt.Errorf("catalog Groups is nil - no catalog provided")
	}

	for id, found := range controlMap {
		if !found {
			message.Debugf("Control %s not found", id)
		}
	}

	props := []oscalTypes_1_1_2.Property{
		{
			Name:  "generation",
			Ns:    LULA_NAMESPACE,
			Value: command,
		},
	}

	if framework != "" {
		prop := oscalTypes_1_1_2.Property{
			Name:  "framework",
			Ns:    LULA_NAMESPACE,
			Value: framework,
		}
		props = append(props, prop)
	}

	if len(implementedRequirements) == 0 {
		return componentDefinition, fmt.Errorf("no controls were identified in the catalog from the requirements list: %v\n", targetControls)
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
					ImplementedRequirements: implementedRequirements,
					Description:             "Control Implementation Description",
					Props:                   &props,
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

func searchGroups(groups *[]oscalTypes_1_1_2.Group, controlMap map[string]bool, remarks []string) ([]oscalTypes_1_1_2.ImplementedRequirementControlImplementation, error) {

	implementedRequirements := make([]oscalTypes_1_1_2.ImplementedRequirementControlImplementation, 0)

	for _, group := range *groups {
		if group.Groups != nil {
			newReqs, err := searchGroups(group.Groups, controlMap, remarks)
			if err != nil {
				return implementedRequirements, err
			}
			implementedRequirements = append(implementedRequirements, newReqs...)
		}
		if group.Controls != nil {
			newReqs, err := searchControls(group.Controls, controlMap, remarks)
			if err != nil {
				return implementedRequirements, err
			}
			implementedRequirements = append(implementedRequirements, newReqs...)
		}
	}
	return implementedRequirements, nil
}

func searchControls(controls *[]oscalTypes_1_1_2.Control, controlMap map[string]bool, remarks []string) ([]oscalTypes_1_1_2.ImplementedRequirementControlImplementation, error) {

	implementedRequirements := make([]oscalTypes_1_1_2.ImplementedRequirementControlImplementation, 0)

	for _, control := range *controls {
		if _, ok := controlMap[control.ID]; ok {
			newRequirement, err := ControlToImplementedRequirement(&control, remarks)
			if err != nil {
				return implementedRequirements, err
			}
			implementedRequirements = append(implementedRequirements, newRequirement)
			controlMap[control.ID] = true
		}

		if control.Controls != nil {
			newReqs, err := searchControls(control.Controls, controlMap, remarks)
			if err != nil {
				return implementedRequirements, err
			}
			implementedRequirements = append(implementedRequirements, newReqs...)
		}
	}
	return implementedRequirements, nil
}

// Consume a control - Identify statements - iterate through parts in order to create a description
func ControlToImplementedRequirement(control *oscalTypes_1_1_2.Control, targetRemarks []string) (implementedRequirement oscalTypes_1_1_2.ImplementedRequirementControlImplementation, err error) {
	var controlDescription string
	paramMap := make(map[string]parameter)

	if control == nil {
		return implementedRequirement, fmt.Errorf("control is nil")
	}

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
		message.Debugf("No parameters (control.Params) found for %s", control.ID)
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
	resourceMap = make(map[string]string)
	if backMatter.Resources == nil {
		return resourceMap
	}

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

func ControlImplementationstToRequirementsMap(controlImplementations *[]oscalTypes_1_1_2.ControlImplementationSet) (requirementMap map[string]Requirement) {
	requirementMap = make(map[string]Requirement)

	if controlImplementations != nil {
		for _, controlImplementation := range *controlImplementations {
			for _, requirement := range controlImplementation.ImplementedRequirements {
				requirementMap[requirement.UUID] = Requirement{
					ImplementedRequirement: &requirement,
					ControlImplementation:  &controlImplementation,
				}
			}
		}
	}
	return requirementMap
}

func FilterControlImplementations(componentDefinition *oscalTypes_1_1_2.ComponentDefinition) (controlMap map[string][]oscalTypes_1_1_2.ControlImplementationSet) {
	controlMap = make(map[string][]oscalTypes_1_1_2.ControlImplementationSet)

	if componentDefinition.Components != nil {
		// Build a map[source/framework][]control-implementations
		for _, component := range *componentDefinition.Components {
			if component.ControlImplementations != nil {
				for _, controlImplementation := range *component.ControlImplementations {
					// Using UUID here as the key -> could also be string -> what would we rather the user pass in?
					controlMap[controlImplementation.Source] = append(controlMap[controlImplementation.Source], controlImplementation)
					status, value := GetProp("framework", LULA_NAMESPACE, controlImplementation.Props)
					if status {
						controlMap[value] = append(controlMap[value], controlImplementation)
					}
				}
			}
		}
	}

	return controlMap
}

// Need to get components + frameworks + control implementations -> new struct?
type ComponentFrameworks struct {
	Component  oscalTypes_1_1_2.DefinedComponent
	Frameworks map[string][]oscalTypes_1_1_2.ControlImplementationSet
}

func NewComponentFrameworks(componentDefinition *oscalTypes_1_1_2.ComponentDefinition) map[string]ComponentFrameworks {
	componentTargets := make(map[string]ComponentFrameworks)

	if componentDefinition.Components != nil {
		// Build a map[source/framework][]control-implementations
		for _, component := range *componentDefinition.Components {
			controlImplementationsMap := make(map[string][]oscalTypes_1_1_2.ControlImplementationSet)
			if component.ControlImplementations != nil {
				for _, controlImplementation := range *component.ControlImplementations {
					// Using UUID here as the key -> could also be string -> what would we rather the user pass in?
					controlImplementationsMap[controlImplementation.Source] = append(controlImplementationsMap[controlImplementation.Source], controlImplementation)
					status, value := GetProp("framework", LULA_NAMESPACE, controlImplementation.Props)
					if status {
						controlImplementationsMap[value] = append(controlImplementationsMap[value], controlImplementation)
					}
				}
			}
			componentTargets[component.UUID] = ComponentFrameworks{
				Component:  component,
				Frameworks: controlImplementationsMap,
			}
		}
	}

	return componentTargets
}

func MakeComponentDeterminstic(component *oscalTypes_1_1_2.ComponentDefinition) {
	// sort components by title

	if component.Components != nil {
		components := *component.Components
		sort.Slice(components, func(i, j int) bool {
			return components[i].Title < components[j].Title
		})

		// sort control-implementations per component by source
		for _, component := range components {
			if component.ControlImplementations != nil {
				controlImplementations := *component.ControlImplementations
				sort.Slice(controlImplementations, func(i, j int) bool {
					return controlImplementations[i].Source < controlImplementations[j].Source
				})
				// sort implemented-requirements per control-implementation by control-id
				for _, controlImplementation := range controlImplementations {
					implementedRequirements := controlImplementation.ImplementedRequirements
					sort.Slice(implementedRequirements, func(i, j int) bool {
						return implementedRequirements[i].ControlId < implementedRequirements[j].ControlId
					})
				}
			}

		}
		component.Components = &components
	}

	// sort capabilities

	if component.Capabilities != nil {
		capabilities := *component.Capabilities
		sort.Slice(capabilities, func(i, j int) bool {
			return capabilities[i].Name < capabilities[j].Name
		})

		for _, capability := range capabilities {

			if capability.ControlImplementations != nil {
				controlImplementations := *capability.ControlImplementations
				sort.Slice(controlImplementations, func(i, j int) bool {
					return controlImplementations[i].Source < controlImplementations[j].Source
				})
				// sort implemented-requirements per control-implementation by control-id
				for _, controlImplementation := range controlImplementations {
					implementedRequirements := controlImplementation.ImplementedRequirements
					sort.Slice(implementedRequirements, func(i, j int) bool {
						return implementedRequirements[i].ControlId < implementedRequirements[j].ControlId
					})
				}
			}

		}
		component.Capabilities = &capabilities
	}

	// sort backmatter
	if component.BackMatter != nil {
		backmatter := *component.BackMatter
		if backmatter.Resources != nil {
			resources := *backmatter.Resources
			sort.Slice(resources, func(i, j int) bool {
				return resources[i].Title < resources[j].Title
			})
			backmatter.Resources = &resources
		}
		component.BackMatter = &backmatter
	}

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

	if part != nil {
		for _, part := range *part {
			// need to get the label first - unsure if there will ever be more than one?
			if part.Props != nil {
				for _, prop := range *part.Props {
					if prop.Name == "label" {
						label = prop.Value
					}
				}
			}

			var tabs string
			// Indents based on labels
			for i := 0; i < level; i++ {
				tabs += "\t"
			}
			// Trims the whitespace
			prose := strings.TrimSpace(part.Prose)
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
