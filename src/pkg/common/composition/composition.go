package composition

import (
	"bytes"
	"fmt"
	"io"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	"github.com/defenseunicorns/go-oscal/src/pkg/versioning"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func ComposeComponentDefinitions(compDef *oscalTypes_1_1_2.ComponentDefinition) error {
	if compDef == nil {
		return fmt.Errorf("component definition is nil")
	}

	// Compose the component validations
	err := ComposeComponentValidations(compDef)
	if err != nil {
		return err
	}

	// If there are no components, create an empty array
	// Components aren't required by oscal but are by merge?
	// TODO: fix merge to match required OSCAL fields
	if compDef.Components == nil {
		compDef.Components = &[]oscalTypes_1_1_2.DefinedComponent{}
	}

	// Same as above
	if compDef.BackMatter == nil {
		compDef.BackMatter = &oscalTypes_1_1_2.BackMatter{}
	}

	if compDef.ImportComponentDefinitions != nil {
		for _, importComponentDef := range *compDef.ImportComponentDefinitions {
			// Fetch the response
			response, err := network.Fetch(importComponentDef.Href)
			if err != nil {
				return err
			}

			// Handle multi-docs
			componentDefs, err := readComponentDefinitionsFromYaml(response)
			if err != nil {
				return err
			}
			// Unmarshal the component definition
			for _, importDef := range componentDefs {
				err = ComposeComponentDefinitions(importDef)
				if err != nil {
					return err
				}

				// Merge the component definitions
				compDef, err = oscal.MergeComponentDefinitions(compDef, importDef)
				if err != nil {
					return err
				}
			}
		}
	}

	compDef.Metadata.LastModified = versioning.GetTimestamp()
	compDef.ImportComponentDefinitions = nil

	return nil
}

// ComposeComponentValidations compiles the component validations by adding the remote resources to the back matter and updating with back matter links.
func ComposeComponentValidations(compDef *oscalTypes_1_1_2.ComponentDefinition) error {

	if compDef == nil {
		return fmt.Errorf("component definition is nil")
	}

	resourceMap := NewResourceStoreFromBackMatter(compDef.BackMatter)

	// If there are no components, there is nothing to do
	if compDef.Components == nil {
		return nil
	}

	for componentIndex, component := range *compDef.Components {
		// If there are no control-implementations, skip to the next component
		controlImplementations := *component.ControlImplementations
		if controlImplementations == nil {
			continue
		}
		for controlImplementationIndex, controlImplementation := range controlImplementations {
			for implementedRequirementIndex, implementedRequirement := range controlImplementation.ImplementedRequirements {
				if implementedRequirement.Links != nil {
					compiledLinks := []oscalTypes_1_1_2.Link{}

					for _, link := range *implementedRequirement.Links {
						if common.IsLulaLink(link) {
							ids, err := resourceMap.AddFromLink(&link)
							if err != nil {
								// return err
								newId := uuid.NewUUID()
								message.Debugf("Error adding validation %s from link %s: %v", newId, link.Href, err)
								ids = []string{newId}
							}
							for _, id := range ids {
								link := oscalTypes_1_1_2.Link{
									Rel:  link.Rel,
									Href: common.AddIdPrefix(id),
									Text: link.Text,
								}
								compiledLinks = append(compiledLinks, link)
							}
						} else {
							compiledLinks = append(compiledLinks, link)
						}
					}
					(*component.ControlImplementations)[controlImplementationIndex].ImplementedRequirements[implementedRequirementIndex].Links = &compiledLinks
					(*compDef.Components)[componentIndex] = component
				}
			}
		}
	}
	allFetched := resourceMap.AllFetched()
	if compDef.BackMatter != nil && compDef.BackMatter.Resources != nil {
		existingResources := *compDef.BackMatter.Resources
		existingResources = append(existingResources, allFetched...)
		compDef.BackMatter.Resources = &existingResources
	} else {
		compDef.BackMatter = &oscalTypes_1_1_2.BackMatter{
			Resources: &allFetched,
		}
	}

	compDef.Metadata.LastModified = versioning.GetTimestamp()

	return nil
}

// ReadComponentDefinitionsFromYaml reads a yaml file of validations to an array of validations
func readComponentDefinitionsFromYaml(componentDefinitionBytes []byte) (componentDefinitionsArray []*oscalTypes_1_1_2.ComponentDefinition, err error) {
	decoder := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(componentDefinitionBytes), 4096)

	for {
		oscalComplete := &oscalTypes_1_1_2.OscalCompleteSchema{}
		if err := decoder.Decode(oscalComplete); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		componentDefinitionsArray = append(componentDefinitionsArray, oscalComplete.ComponentDefinition)
	}

	return componentDefinitionsArray, nil
}
