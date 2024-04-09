package oscal

import (
	"fmt"
	"strings"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/types"
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
func NewOscalComponentDefinition(data []byte) (componentDefinition oscalTypes_1_1_2.ComponentDefinition, err error) {
	var oscalModels oscalTypes_1_1_2.OscalModels

	err = yaml.Unmarshal(data, &oscalModels)
	if err != nil {
		return componentDefinition, err
	}

	if oscalModels.ComponentDefinition == nil {
		return componentDefinition, fmt.Errorf("No Component Definition found in the provided data")
	}

	return *oscalModels.ComponentDefinition, nil
}

func ComponentFromCatalog(catalog oscalTypes_1_1_2.Catalog, targetControls []string) (componentDefinition oscalTypes_1_1_2.ComponentDefinition, err error) {

	// How will we identify a control based on the array of controls passed?
	controlMap := make(map[string][]string)
	for _, control := range targetControls {
		id := strings.Split(control, "-")
		controlMap[id[0]] = append(controlMap[id[0]], id[1])
	}

	// We then want to iterate through the catalog, identify the controls and map control information to the implemented-requirements

	for _, group := range *catalog.Groups {
		// Is this a group of controls that we are targeting
		if controlArray, ok := controlMap[group.ID]; !ok {
			for _, control := range *group.Controls {
				id := strings.Split(control.ID, "-")
				if contains(controlArray, id[1]) {

				}

				// Identify if this is a control we are targeting

				// We need to assemble the implemented requirement description from the

			}
		}

	}

	return componentDefinition, nil

}

func ControlToImplementedRequirement(control oscalTypes_1_1_2.Control) (implementedRequirement oscalTypes_1_1_2.ImplementedRequirement, err error) {

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
				controlDescription = addPart(part.Parts, paramMap)
			}

		}

		// Check for parts != nil
		// if not nil, recurse
		// if nil, return the above string
		// add a tab as necessary

	}
	return implementedRequirement, nil
}

// Map an array of resources to a map of UUID to validation object
func BackMatterToMap(backMatter oscalTypes_1_1_2.BackMatter) map[string]types.Validation {
	resourceMap := make(map[string]types.Validation)

	if backMatter.Resources == nil {
		return nil
	}

	for _, resource := range *backMatter.Resources {
		if resource.Title == "Lula Validation" {
			var validation types.Validation

			err := yaml.Unmarshal([]byte(resource.Description), &validation)
			if err != nil {
				fmt.Printf("Error marshalling yaml: %s\n", err.Error())
				return nil
			}

			// Do version checking here to establish if the version is correct/acceptable
			var result types.Result
			var evaluated bool
			currentVersion := strings.Split(config.CLIVersion, "-")[0]

			versionConstraint := currentVersion
			if validation.LulaVersion != "" {
				versionConstraint = validation.LulaVersion
			}

			validVersion, versionErr := common.IsVersionValid(versionConstraint, currentVersion)
			if versionErr != nil {
				result.Failing = 1
				result.Observations = map[string]string{"Lula Version Error": versionErr.Error()}
				evaluated = true
			} else if !validVersion {
				result.Failing = 1
				result.Observations = map[string]string{"Version Constraint Incompatible": "Lula Version does not meet the constraint for this validation."}
				evaluated = true
			}

			validation.Title = resource.Title
			validation.Evaluated = evaluated
			validation.Result = result

			resourceMap[resource.UUID] = validation
		}

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
func addPart(part *[]oscalTypes_1_1_2.Part, paramMap map[string]parameter) string {

	prose := part.Prose

	// Check for parameter insertion

	// Perform a replace on the prose

}
