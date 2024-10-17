package composition

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	"github.com/defenseunicorns/go-oscal/src/pkg/versioning"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

type RenderedContent string

type Composer struct {
	modelDir          string
	templateRenderer  *template.TemplateRenderer
	renderTemplate    bool
	renderValidations bool
	renderType        template.RenderType
}

func New(opts ...Option) (*Composer, error) {
	var composer Composer

	for _, opt := range opts {
		if err := opt(&composer); err != nil {
			return nil, err
		}
	}

	return &composer, nil
}

// ComposeFromPath composes an OSCAL model from a file path
func (c *Composer) ComposeFromPath(ctx context.Context, path string) (model *oscalTypes_1_1_2.OscalCompleteSchema, err error) {
	path = filepath.Clean(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Template if renderTemplate is true -> Only renders the local data (e.g., what is in the file)
	if c.renderTemplate {
		data, err = c.templateRenderer.Render(string(data), c.renderType)
		if err != nil {
			return nil, err
		}
	}

	model, err = oscal.NewOscalModel(data)
	if err != nil {
		return nil, err
	}

	err = c.ComposeComponentDefinitions(ctx, model.ComponentDefinition, c.modelDir)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// ComposeComponentDefinitions composes an OSCAL component definition by adding the remote resources to the back matter and updating with back matter links.
func (c *Composer) ComposeComponentDefinitions(ctx context.Context, compDef *oscalTypes_1_1_2.ComponentDefinition, baseDir string) error {
	if compDef == nil {
		return fmt.Errorf("component definition is nil")
	}

	// Compose the component validations
	err := c.ComposeComponentValidations(ctx, compDef, baseDir)
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
			response, err := network.Fetch(importComponentDef.Href, network.WithBaseDir(baseDir))
			if err != nil {
				return err
			}

			// template here if renderTemplate is true
			if c.renderTemplate {
				response, err = c.templateRenderer.Render(string(response), c.renderType)
				if err != nil {
					return err
				}
			}

			// Handle multi-docs
			componentDefs, err := readComponentDefinitionsFromYaml(response)
			if err != nil {
				return err
			}
			// Unmarshal the component definition
			for _, importDef := range componentDefs {
				// Reconcile the base directory from the import component definition href
				importDir := network.GetLocalFileDir(importComponentDef.Href, baseDir)
				err = c.ComposeComponentDefinitions(ctx, importDef, importDir)
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
func (c *Composer) ComposeComponentValidations(ctx context.Context, compDef *oscalTypes_1_1_2.ComponentDefinition, baseDir string) error {

	if compDef == nil {
		return fmt.Errorf("component definition is nil")
	}

	resourceMap := NewResourceStoreFromBackMatter(c, compDef.BackMatter)

	// If there are no components, there is nothing to do
	if compDef.Components == nil {
		return nil
	}

	for componentIndex, component := range *compDef.Components {
		// If there are no control-implementations, skip to the next component
		if component.ControlImplementations != nil {
			controlImplementations := *component.ControlImplementations
			for controlImplementationIndex, controlImplementation := range controlImplementations {
				for implementedRequirementIndex, implementedRequirement := range controlImplementation.ImplementedRequirements {
					if implementedRequirement.Links != nil {
						compiledLinks := []oscalTypes_1_1_2.Link{}

						for _, link := range *implementedRequirement.Links {
							if common.IsLulaLink(link) {
								ids, err := resourceMap.AddFromLink(&link, baseDir)
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

// CreateTempDir creates a temporary directory to store the composed OSCAL models
func CreateTempDir() (string, error) {
	return os.MkdirTemp("", "lula-composed-*")
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
