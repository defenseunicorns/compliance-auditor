package composition

import (
	"fmt"
	"path/filepath"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

type Option func(*CompositionContext) error

// TODO: add remote option?
func WithModelFromPath(path string) Option {
	return func(ctx *CompositionContext) error {
		if err := files.IsJsonOrYaml(path); err != nil {
			return fmt.Errorf("invalid file extension: %s, requires .json or .yaml", path)
		}
		ctx.modelDir = filepath.Dir(path)
		return nil
	}
}

func WithTemplateRenderer(renderTypeString string, renderRemote bool, setOpts []string) Option {
	return func(ctx *CompositionContext) error {
		if renderTypeString == "" {
			if len(setOpts) > 0 {
				message.Warn("`render` not specified, the --set options will be ignored")
			}
			if renderRemote {
				message.Warn("`render` not specified, `render-remote` will be ignored")
			}
			return nil
		}

		ctx.renderTemplate = true
		ctx.renderRemote = renderRemote

		// Get the template render type
		renderType, err := template.ParseRenderType(renderTypeString)
		if err != nil {
			message.Warnf("invalid render type, defaulting to non-sensitive: %v", err)
			renderType = template.NONSENSITIVE
		}
		ctx.renderType = renderType

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
		tr := template.NewTemplateRenderer(templateData)

		ctx.templateRenderer = tr

		return nil
	}
}
