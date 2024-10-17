package tools

import (
	"fmt"
	"os"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

var templateHelp = `
To template an OSCAL Model, defaults to masking sensitive variables:
	lula tools template -f ./oscal-component.yaml

To indicate a specific output file:
	lula tools template -f ./oscal-component.yaml -o templated-oscal-component.yaml

To perform overrides on the template data:
	lula tools template -f ./oscal-component.yaml --set .var.key1=value1 --set .const.key2=value2

To perform the full template operation, including sensitive data:
	lula tools template -f ./oscal-component.yaml --render all

Data for templating should be stored under 'constants' or 'variables' configuration items in a lula-config.yaml file
See documentation for more detail on configuration schema
`

func TemplateCommand() *cobra.Command {
	var (
		inputFile        string
		outputFile       string
		setOpts          []string
		renderTypeString string
	)

	cmd := &cobra.Command{
		Use:     "template",
		Short:   "Template an artifact",
		Long:    "Resolving templated artifacts with configuration data",
		Args:    cobra.NoArgs,
		Example: templateHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read file
			data, err := network.Fetch(inputFile)
			if err != nil {
				return fmt.Errorf("error reading file: %v", err)
			}

			// Validate render type
			renderType, err := template.ParseRenderType(renderTypeString)
			if err != nil {
				message.Warnf("invalid render type, defaulting to masked: %v", err)
				renderType = template.MASKED
			}

			// Get overrides from --set flag
			overrides, err := common.ParseTemplateOverrides(setOpts)
			if err != nil {
				return fmt.Errorf("error parsing template overrides: %v", err)
			}

			// Handles merging viper config file data + environment variables
			// Throws an error if config keys are invalid for templating
			templateData, err := template.CollectTemplatingData(common.TemplateConstants, common.TemplateVariables, overrides)
			if err != nil {
				return fmt.Errorf("error collecting templating data: %v", err)
			}

			templateRenderer := template.NewTemplateRenderer(templateData)
			output, err := templateRenderer.Render(string(data), renderType)
			if err != nil {
				return fmt.Errorf("error rendering template: %v", err)
			}

			if outputFile == "" {
				_, err := cmd.OutOrStdout().Write(output)
				if err != nil {
					return fmt.Errorf("failed to write to stdout: %v", err)
				}
			} else {
				err = files.CreateFileDirs(outputFile)
				if err != nil {
					return fmt.Errorf("failed to create output file path: %v", err)
				}
				err = os.WriteFile(outputFile, output, 0644)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input-file", "f", "", "the path to the target artifact")
	err := cmd.MarkFlagRequired("input-file")
	if err != nil {
		message.Fatal(err, "error initializing template command flags")
	}
	cmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will be directed to stdout")
	cmd.Flags().StringSliceVarP(&setOpts, "set", "s", []string{}, "set a value in the template data")
	cmd.Flags().StringVarP(&renderTypeString, "render", "r", "masked", "values to render the template with, options are: masked, constants, non-sensitive, all")

	return cmd
}

func init() {
	common.InitViper()
	toolsCmd.AddCommand(TemplateCommand())
}
