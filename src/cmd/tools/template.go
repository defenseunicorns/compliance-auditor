package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/internal/template"
	pkgCommon "github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

var templateHelp = `
To template an OSCAL Model, defaults to masking sensitive variables:
	lula tools template -f ./oscal-component.yaml

To indicate a specific output file:
	lula tools template -f ./oscal-component.yaml -o templated-oscal-component.yaml

To perform overrides on the template data:
	lula tools template -f ./oscal-component.yaml --set var.key1=value1 --set const.key2=value2

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
			// Get current viper pointer
			v := common.GetViper()

			// Read file
			data, err := pkgCommon.ReadFileToBytes(inputFile)
			if err != nil {
				return fmt.Errorf("error reading file: %v", err)
			}

			// Validate render type
			renderType, err := getRenderType(renderTypeString)
			if err != nil {
				message.Warn("invalid render type, defaulting to masked")
			}

			// Get constants and variables for templating from viper config
			var constants map[string]interface{}
			var variables []template.VariableConfig

			err = v.UnmarshalKey(common.VConstants, &constants)
			if err != nil {
				return fmt.Errorf("unable to unmarshal constants into map: %v", err)
			}

			err = v.UnmarshalKey(common.VVariables, &variables)
			if err != nil {
				return fmt.Errorf("unable to unmarshal variables into slice: %v", err)
			}

			// Get overrides from --set flag
			overrides := getOverrides(setOpts)

			// Handles merging viper config file data + environment variables
			// Throws an error if config keys are invalid for templating
			templateData, err := template.CollectTemplatingData(constants, variables, overrides)
			if err != nil {
				return fmt.Errorf("error collecting templating data: %v", err)
			}

			templateRenderer := template.NewTemplateRenderer(string(data), renderType, templateData)
			output, err := templateRenderer.Render()
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
	cmd.MarkFlagRequired("input-file")
	cmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will be directed to stdout")
	cmd.Flags().StringSliceVarP(&setOpts, "set", "s", []string{}, "set a value in the template data")
	cmd.Flags().StringVarP(&renderTypeString, "render", "r", "masked", "values to render the template with, options are: masked, constants, non-sensitive, all")

	return cmd
}

func init() {
	common.InitViper()
	toolsCmd.AddCommand(TemplateCommand())
}

func getRenderType(item string) (template.RenderType, error) {
	switch strings.ToLower(item) {
	case "masked":
		return template.MASKED, nil
	case "constants":
		return template.CONSTANTS, nil
	case "non-sensitive":
		return template.NONSENSITIVE, nil
	case "all":
		return template.ALL, nil
	}
	return template.MASKED, fmt.Errorf("invalid render type: %s", item)
}

func getOverrides(setFlags []string) map[string]string {
	overrides := make(map[string]string)
	for _, flag := range setFlags {
		parts := strings.SplitN(flag, "=", 2)
		if len(parts) != 2 {
			message.Fatalf(fmt.Errorf("invalid --set flag format, should be .root.key=value"), "invalid --set flag format, should be .root.key=value")
		}

		if !strings.HasPrefix(parts[0], "."+template.CONST+".") && !strings.HasPrefix(parts[0], "."+template.VAR+".") {
			message.Fatalf(fmt.Errorf("invalid --set flag format, path should start with .const or .var"), "invalid --set flag format, path should start with .const or .var")
		}

		path, value := parts[0], parts[1]
		overrides[path] = value
	}
	return overrides
}
