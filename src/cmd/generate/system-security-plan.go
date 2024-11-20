package generate

import (
	"fmt"
	"os"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/spf13/cobra"

	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

var sspExample = `
To generate a system security plan from source and component definition:
	lula generate system-security-plan -s <catalog/profile source> -c <path/to/component-definition>

To specify the name and filetype of the generated artifact:
	lula generate system-security-plan -s <catalog/profile source> -c <path/to/component-definition> -o my_ssp.yaml
`

var sspLong = `Generation of a System Security Plan OSCAL artifact from a source catalog/profile and a component definition.`

func GenerateSSPCommand() *cobra.Command {
	var (
		component  []string
		source     string
		outputFile string
	)

	sspCmd := &cobra.Command{
		Use:     "system-security-plan",
		Aliases: []string{"ssp"},
		Short:   "Generate a system security plan OSCAL template",
		Long:    sspLong,
		Example: sspExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			message.Info("generate system-security-plan executed")

			if outputFile == "" {
				outputFile = "system-security-plan.yaml"
			}

			// Check if output file contains a valid OSCAL model
			_, err := oscal.ValidOSCALModelAtPath(outputFile)
			if err != nil {
				return fmt.Errorf("invalid OSCAL model at output: %v", err)
			}

			command := fmt.Sprintf("%s --source %s", cmd.CommandPath(), source)

			// Get component definitions from file(s)
			var componentDefs []*oscalTypes.ComponentDefinition
			for _, componentPath := range component {
				componentBytes, err := os.ReadFile(componentPath)
				if err != nil {
					return err
				}
				componentDef, err := oscal.NewOscalComponentDefinition(componentBytes)
				if err != nil {
					return err
				}
				if componentDef == nil {
					return fmt.Errorf("component definition at %s is nil", componentPath)
				}

				// TODO: Compose component definition
				// Add option: compose to path (remap validation links)
				// Add option: partial compose(?) just imported component-definitions

				componentDefs = append(componentDefs, componentDef)
				command += fmt.Sprintf(" --component %s", componentPath)
			}

			// Generate the system security plan
			ssp, err := oscal.GenerateSystemSecurityPlan(command, source, componentDefs...)
			if err != nil {
				return err
			}

			// Write the system security plan to file
			err = oscal.WriteOscalModelNew(outputFile, ssp)
			if err != nil {
				message.Fatalf(err, "error writing component to file")
			}

			return nil
		},
	}

	sspCmd.Flags().StringVarP(&source, "source", "s", "", "the path to the source catalog/profile")
	err := sspCmd.MarkFlagRequired("source")
	if err != nil {
		message.Fatal(err, "error initializing source command flags")
	}
	sspCmd.Flags().StringSliceVarP(&component, "component", "c", []string{}, "comma delimited list the paths to the component definitions to include for the SSP")
	err = sspCmd.MarkFlagRequired("component")
	if err != nil {
		message.Fatal(err, "error initializing component command flags")
	}
	sspCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to the output file. If not specified, the output file will default to `system-security-plan.yaml`")

	return sspCmd
}
