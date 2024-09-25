package validation

import (
	"fmt"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

type Option func(*ValidationContext) error

// TODO: multiple paths? comments in validate cmd indicate multiple paths should be supported, but aren't currently
func WithComponentDefinitionFromPath(path string) Option {
	return func(ctx *ValidationContext) error {
		if err := files.IsJsonOrYaml(path); err != nil {
			return fmt.Errorf("invalid file extension: %s, requires .json or .yaml", path)
		}
		ctx.componentPath = path
		return nil
	}
}

func WithComponentDefinition(componentDefinition *oscalTypes_1_1_2.ComponentDefinition) Option {
	return func(ctx *ValidationContext) error {
		if componentDefinition == nil {
			return fmt.Errorf("component definition is nil")
		}
		ctx.componentDefinition = componentDefinition
		return nil
	}
}

func WithTemplateRenderer(renderTemplate bool, setOpts []string) Option {
	return func(ctx *ValidationContext) error {
		tr := new(template.TemplateRenderer)

		if renderTemplate {
			// Get constants and variables for templating from viper config
			constants, variables, err := common.GetTemplateConfig()
			if err != nil {
				return fmt.Errorf("error getting template config: %v", err)
			}

			// Get overrides from setOpts flag
			overrides, err := common.ParseTemplateOverrides(setOpts)
			if err != nil {
				return fmt.Errorf("error parsing template overrides: %v", err)
			}

			// Handles merging viper config file data + environment variables
			// Throws an error if config keys are invalid for templating
			templateData, err := template.CollectTemplatingData(constants, variables, overrides)
			if err != nil {
				return fmt.Errorf("error collecting templating data: %v", err)
			}

			// need to update the template with the templateString...
			tr = template.NewTemplateRenderer(templateData)
		}

		ctx.renderTemplate = renderTemplate
		ctx.templateRenderer = tr

		return nil
	}
}

func WithAllowExecution(confirmExecution, runNonInteractively bool) Option {
	return func(ctx *ValidationContext) error {
		if !confirmExecution {
			if !runNonInteractively {
				ctx.requestExecutionConfirmation = true
			} else {
				message.Infof("Validations requiring execution will NOT be run")
			}
		} else {
			ctx.runExecutableValidations = true
		}
		return nil
	}
}

func WithResourcesDir(saveResources bool, rootDir string) Option {
	return func(ctx *ValidationContext) error {
		if saveResources {
			ctx.resourcesDir = rootDir
		}
		return nil
	}
}
