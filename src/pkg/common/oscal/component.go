package oscal

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"sigs.k8s.io/yaml"

	"github.com/defenseunicorns/lula/src/pkg/message"
)

type Requirement struct {
	ImplementedRequirement *oscalTypes.ImplementedRequirementControlImplementation
	ControlImplementation  *oscalTypes.ControlImplementationSet
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

type ComponentDefinition struct {
	Model *oscalTypes.ComponentDefinition
}

func NewComponentDefinition() *ComponentDefinition {
	var compdef ComponentDefinition
	compdef.Model = nil
	return &compdef
}

// Create a new ComponentDefinition model
func (c *ComponentDefinition) NewModel(data []byte) error {

	var oscalModels oscalTypes.OscalModels

	err := multiModelValidate(data)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &oscalModels)
	if err != nil {
		return err
	}

	c.Model = oscalModels.ComponentDefinition
	if c.Model == nil {
		return fmt.Errorf("unable to find component definition model")
	}

	return nil
}

func (*ComponentDefinition) GetType() string {
	return "component-definition"
}

func (c *ComponentDefinition) GetCompleteModel() *oscalTypes.OscalModels {
	return &oscalTypes.OscalModels{
		ComponentDefinition: c.Model,
	}
}

func (c *ComponentDefinition) ImportComponentDefinitions() error {
	// TODO: get all imported component definitions and re-write any relative links?

	return nil
}

// TODO: Add other interface methods(?)

// NewOscalComponentDefinition consumes a byte array and returns a new single OscalComponentDefinitionModel object
// Standard use is to read a file from the filesystem and pass the []byte to this function
func NewOscalComponentDefinition(data []byte) (componentDefinition *oscalTypes.ComponentDefinition, err error) {
	var oscalModels oscalTypes.OscalModels

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

// MergeVariadicComponentDefinition merges multiple variadic component definitions into a single component definition
func MergeVariadicComponentDefinition(compdefs ...*oscalTypes.ComponentDefinition) (mergedCompDef *oscalTypes.ComponentDefinition, err error) {
	for _, compdef := range compdefs {
		if mergedCompDef == nil {
			mergedCompDef = compdef
		} else {
			mergedCompDef, err = MergeComponentDefinitions(mergedCompDef, compdef)
			if err != nil {
				return nil, err
			}
		}
	}
	return mergedCompDef, nil
}

// This function should perform a merge of two component-definitions where maintaining the original component-definition is the primary concern.
func MergeComponentDefinitions(original *oscalTypes.ComponentDefinition, latest *oscalTypes.ComponentDefinition) (*oscalTypes.ComponentDefinition, error) {

	originalMap := make(map[string]oscalTypes.DefinedComponent)

	if original.Components == nil {
		return original, fmt.Errorf("original component-definition is nil")
	}

	if latest.Components == nil {
		return original, fmt.Errorf("latest component-definition is nil")
	}

	for _, component := range *original.Components {
		originalMap[component.Title] = component
	}

	latestMap := make(map[string]oscalTypes.DefinedComponent)

	for _, component := range *latest.Components {
		latestMap[component.Title] = component
	}

	tempItems := make([]oscalTypes.DefinedComponent, 0)
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
		original.BackMatter = &oscalTypes.BackMatter{
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

func mergeComponents(original *oscalTypes.DefinedComponent, latest *oscalTypes.DefinedComponent) *oscalTypes.DefinedComponent {
	originalMap := make(map[string]oscalTypes.ControlImplementationSet)

	if original.ControlImplementations != nil {
		for _, item := range *original.ControlImplementations {
			originalMap[item.Source] = item
		}
	}

	latestMap := make(map[string]oscalTypes.ControlImplementationSet)

	if latest.ControlImplementations != nil {
		for _, item := range *latest.ControlImplementations {
			latestMap[item.Source] = item
		}
	}

	tempItems := make([]oscalTypes.ControlImplementationSet, 0)
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

func mergeControlImplementations(original *oscalTypes.ControlImplementationSet, latest *oscalTypes.ControlImplementationSet) *oscalTypes.ControlImplementationSet {
	originalMap := make(map[string]oscalTypes.ImplementedRequirementControlImplementation)

	if original.ImplementedRequirements != nil {
		for _, item := range original.ImplementedRequirements {
			originalMap[item.ControlId] = item
		}
	}
	latestMap := make(map[string]oscalTypes.ImplementedRequirementControlImplementation)

	if latest.ImplementedRequirements != nil {
		for _, item := range latest.ImplementedRequirements {
			latestMap[item.ControlId] = item
		}
	}

	tempItems := make([]oscalTypes.ImplementedRequirementControlImplementation, 0)
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

// Merges two arrays of links into a single array
// TODO: account for overriding validations
func mergeLinks(orig []oscalTypes.Link, latest []oscalTypes.Link) *[]oscalTypes.Link {
	result := make([]oscalTypes.Link, 0)

	tempLinks := make(map[string]oscalTypes.Link)
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
func ComponentFromCatalog(command string, source string, catalog *oscalTypes.Catalog, componentTitle string, targetControls []string, targetRemarks []string, framework string) (*oscalTypes.ComponentDefinition, error) {
	// store all of the implemented requirements
	implementedRequirements := make([]oscalTypes.ImplementedRequirementControlImplementation, 0)
	var componentDefinition = &oscalTypes.ComponentDefinition{}

	if len(targetControls) == 0 {
		return componentDefinition, fmt.Errorf("no controls identified for generation")
	}

	controlsToImplement, err := ResolveCatalogControls(catalog, targetControls, nil)
	if err != nil {
		return componentDefinition, err
	}

	if len(controlsToImplement) == 0 {
		return componentDefinition, fmt.Errorf("no controls were identified in the catalog from the requirements list: %v\n", targetControls)
	}

	for _, control := range controlsToImplement {
		ir, err := ControlToImplementedRequirement(&control, targetRemarks)
		if err != nil {
			return componentDefinition, fmt.Errorf("error creating implemented requirement: %v", err)
		}
		implementedRequirements = append(implementedRequirements, ir)
	}

	for _, id := range targetControls {
		if _, ok := controlsToImplement[id]; !ok {
			message.Debugf("Control %s not found", id)
		}
	}

	props := []oscalTypes.Property{
		{
			Name:  "generation",
			Ns:    LULA_NAMESPACE,
			Value: command,
		},
	}

	if framework != "" {
		prop := oscalTypes.Property{
			Name:  "framework",
			Ns:    LULA_NAMESPACE,
			Value: framework,
		}
		props = append(props, prop)
	}

	componentDefinition.Components = &[]oscalTypes.DefinedComponent{
		{
			UUID:        uuid.NewUUID(),
			Type:        "software",
			Title:       componentTitle,
			Description: "Component Description",
			ControlImplementations: &[]oscalTypes.ControlImplementationSet{
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

	componentDefinition.Metadata = oscalTypes.Metadata{
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
func ControlToImplementedRequirement(control *oscalTypes.Control, targetRemarks []string) (implementedRequirement oscalTypes.ImplementedRequirementControlImplementation, err error) {
	remarks, err := getControlRemarks(control, targetRemarks)
	if err != nil {
		return implementedRequirement, err
	}

	// assemble implemented-requirements object
	implementedRequirement.Remarks = remarks
	implementedRequirement.Description = "<how the specified control may be implemented if the containing component or capability is instantiated in a system security plan>"
	implementedRequirement.ControlId = control.ID
	implementedRequirement.UUID = uuid.NewUUID()

	return implementedRequirement, nil
}

// Returns a map of the uuid - description of the back-matter resources
func BackMatterToMap(backMatter oscalTypes.BackMatter) (resourceMap map[string]string) {
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

func ControlImplementationstToRequirementsMap(controlImplementations *[]oscalTypes.ControlImplementationSet) (requirementMap map[string]Requirement) {
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

func FilterControlImplementations(componentDefinition *oscalTypes.ComponentDefinition) (controlMap map[string][]oscalTypes.ControlImplementationSet) {
	controlMap = make(map[string][]oscalTypes.ControlImplementationSet)

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

func ComponentsToMap(componentDefinition *oscalTypes.ComponentDefinition) map[string]*oscalTypes.DefinedComponent {
	components := make(map[string]*oscalTypes.DefinedComponent)

	if componentDefinition != nil && componentDefinition.Components != nil {
		for _, component := range *componentDefinition.Components {
			components[component.UUID] = &component
		}
	}
	return components
}

func MakeComponentDeterminstic(component *oscalTypes.ComponentDefinition) {
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
		sortBackMatter(&backmatter)
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
func addPart(part *[]oscalTypes.Part, paramMap map[string]parameter, level int) string {

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
				result += fmt.Sprintf("%s%s %s\n", tabs, label, replaceParams(prose, paramMap, false))
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

func replaceParams(input string, params map[string]parameter, nested bool) string {
	re := regexp.MustCompile(`{{\s*insert:\s*param,\s*([^}\s]+)\s*}}`)
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		paramName := strings.TrimSpace(re.FindStringSubmatch(match)[1])
		if param, ok := params[paramName]; ok {
			if nested {
				// If we know there is no information that requires prepending
				return param.Label
			} else if param.Select == nil {
				// Standard assignment given a label
				return fmt.Sprintf("[Assignment: organization-defined %s]", param.Label)
			} else {
				// Join many choices into a single item of prose
				prose := fmt.Sprintf("[Selection: (%s) organization-defined", param.Select.HowMany)
				for _, choice := range param.Select.Choice {
					// Handle potential nested parameter use
					if strings.Contains(choice, "{{ insert: param,") {
						replaced := replaceParams(choice, params, true)
						prose += fmt.Sprintf(" %s;", replaced)
					} else {
						prose += fmt.Sprintf(" %s;", choice)
					}
				}
				return fmt.Sprintf("%s]", prose)
			}
		}
		return match
	})
	return result
}
