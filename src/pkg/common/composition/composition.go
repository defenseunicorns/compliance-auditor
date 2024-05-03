package composition

import (
	"fmt"

	gooscalUtils "github.com/defenseunicorns/go-oscal/src/pkg/utils"
	"github.com/defenseunicorns/go-oscal/src/pkg/validation"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
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

	if compDef.Components == nil {
		compDef.Components = &[]oscalTypes_1_1_2.DefinedComponent{}
	}

	if compDef.BackMatter == nil {
		compDef.BackMatter = &oscalTypes_1_1_2.BackMatter{}
	}

	if compDef.ImportComponentDefinitions != nil {
		for _, importComponentDef := range *compDef.ImportComponentDefinitions {
			// Fetch the file
			file, err := network.Fetch(importComponentDef.Href)
			if err != nil {
				return err
			}
			// Unmarshal the component definition
			importDef, err := oscal.NewOscalComponentDefinitionFromBytes(file)
			if err != nil {
				return err
			}

			validator, err := validation.NewValidator(file)
			if err != nil {
				return err
			}

			err = validator.Validate()
			if err != nil {
				return err
			}

			err = ComposeComponentDefinitions(&importDef)
			if err != nil {
				return err
			}

			// Merge the component definitions
			*compDef, err = oscal.MergeComponentDefinitions(*compDef, importDef)
			if err != nil {
				return err
			}
		}
	}

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
							ids, err := resourceMap.AddFromLink(link)
							if err != nil {
								return err
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

	compDef.Metadata.LastModified = gooscalUtils.GetTimestamp()

	return nil
}
