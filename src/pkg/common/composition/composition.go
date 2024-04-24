package composition

import (
	"fmt"

	gooscalUtils "github.com/defenseunicorns/go-oscal/src/pkg/utils"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
)

// ComposeComponentValidations compiles the component validations by adding the remote resources to the back matter and updating with back matter links.
func ComposeComponentValidations(compDef *oscalTypes_1_1_2.ComponentDefinition) error {

	if compDef == nil {
		return fmt.Errorf("component definition is nil")
	}

	resourceMap := NewResourceStoreFromBackMatter(compDef.BackMatter)

	if *compDef.Components == nil {
		return fmt.Errorf("no components found in component definition")
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
