package generate

import (
	"bytes"
	"fmt"
	"os"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"

	"gopkg.in/yaml.v3"
)

// The generate command is expected to have multiple different use-cases based upon
// the artifact to be generated. Some opinionation will be involved - such as OSCAL
// for the supported OSCAL models. This leaves generate open to generating other
// artifacts in the future that support assessment and accreditation processes.

type flags struct {
	InputFile  string // -f --input-file
	OutputFile string // -o --output-file
}

type componentFlags struct {
	flags
	CatalogSource string   // -c --catalog
	Profile       string   // -p --profile
	Requirements  []string // -r --requirements
}

var opts = &flags{}
var componentOpts = &componentFlags{}

// Base generate command will handle a large E2E focused generation that is driven from various artifacts.
// This will include some prerequisites such as component-definitions - but will ultimately lead to generation
// of the SSP / SAP / POAM / and maybe SAR in a maintainable way.
var generateCmd = &cobra.Command{
	Use:     "generate",
	Hidden:  false, // Hidden for now until fully implemented
	Aliases: []string{"g", "gen"},
	Short:   "Generate a specified compliance artifact template",
}

// Component-Definition generation will generate an OSCAL file that can be used both as the basis for Lula validations
// as well as required components for SSP/SAP/SAR/POAM.
var generateComponentCmd = &cobra.Command{
	Use:     "component",
	Aliases: []string{"c"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Generate a component definition OSCAL template",
	Run: func(_ *cobra.Command, args []string) {
		message.Info("generate component executed")

		// check for inputFile flag content
		if componentOpts.CatalogSource == "" {
			message.Fatal(fmt.Errorf("no catalog source provided"), "generate component requires a catalog input source")
		}

		// var existingComponent oscalTypes_1_1_2.ComponentDefinition
		// if componentOpts.OutputFile != "" {
		// 	// We meed tp check if the file exists
		// 	if _, err := os.Stat(componentOpts.OutputFile); err == nil {
		// 		// if the file exists, we need to read it into bytes
		// 		existingFileBytes, err := os.ReadFile(componentOpts.OutputFile)
		// 		if err != nil {
		// 			message.Fatalf(fmt.Errorf("error reading existing file"), "error reading existing file")
		// 		}
		// 		existingComponent, err = oscal.NewOscalComponentDefinition(componentOpts.OutputFile, existingFileBytes)
		// 	}
		// }

		// Existing component has now potentially been identified - do something with it.
		// TODO: start here

		source := componentOpts.CatalogSource

		data, err := network.Fetch(source)
		if err != nil {
			message.Fatalf(fmt.Errorf("error fetching catalog source"), "error fetching catalog source")
		}

		catalog, err := oscal.NewCatalog(source, data)
		if err != nil {
			message.Fatalf(fmt.Errorf("error creating catalog"), "error creating catalog")
		}

		comp, err := oscal.ComponentFromCatalog(source, catalog, componentOpts.Requirements)
		if err != nil {
			message.Fatalf(fmt.Errorf("error creating component"), "error creating component")
		}

		var fileName string
		if opts.OutputFile != "" {
			fileName = opts.OutputFile
		} else {
			fileName = "new-component.yaml"
		}

		var b bytes.Buffer

		var component = oscalTypes_1_1_2.OscalModels{
			ComponentDefinition: &comp,
		}

		yamlEncoder := yaml.NewEncoder(&b)
		yamlEncoder.SetIndent(2)
		yamlEncoder.Encode(component)

		message.Infof("Writing Security Assessment Results to: %s", fileName)

		err = os.WriteFile(fileName, b.Bytes(), 0644)
		if err != nil {
			message.Fatalf(fmt.Errorf("error writing Security Assessment Results to: %s", fileName), "error writing Security Assessment Results to: %s", fileName)
		}

	},
}

var generateAssessmentPlanCmd = &cobra.Command{
	Use:     "assessment-plan",
	Aliases: []string{"ap"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Generate an assessment plan OSCAL template",
	Run: func(_ *cobra.Command, args []string) {
		message.Info("generate assessment-plan executed")

		// For each component-definition in array of component-definitions
		// Read component-definition
		// Collect all implemented-requirements
		// Collect all items from the backmatter
		// Create new assessment-plan object
		// Transfer to assessment-plan.reviewed-controls?
	},
}

var generateSystemSecurityPlanCmd = &cobra.Command{
	Use:     "system-security-plan",
	Aliases: []string{"ssp"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Generate a system security plan OSCAL template",
	Run: func(_ *cobra.Command, args []string) {
		message.Info("generate system-security-plan executed")

		// For each component-definition in array of component-definitions
		// Read component-definition
		// Collect all implemented-requirements
		// aggregate by control-id
		// export to system-security-plan.control-implementation.implemented-requirements
	},
}

var generatePOAMCmd = &cobra.Command{
	Use:     "plan-of-action-and-milestones",
	Aliases: []string{"poam"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Generate a plan of actions and milestones OSCAL template",
	Run: func(_ *cobra.Command, args []string) {
		message.Info("generate plan-of-action-and-milestones executed")

		// Locate an assessment-results artifact
		// Create an assessment-results object
		// Transfer 'not-satisfied' findings and observations to poam.findings and poam.observations as appropriate
	},
}

func GenerateCommand() *cobra.Command {

	generateCmd.AddCommand(generateComponentCmd)
	generateCmd.AddCommand(generateAssessmentPlanCmd)
	generateCmd.AddCommand(generateSystemSecurityPlanCmd)
	generateCmd.AddCommand(generatePOAMCmd)

	generateFlags()
	generateComponentFlags()

	return generateCmd
}

func generateFlags() {
	generateFlags := generateCmd.PersistentFlags()

	generateFlags.StringVarP(&opts.InputFile, "input-file", "f", "", "Path to a manifest file")
	generateFlags.StringVarP(&opts.OutputFile, "output-file", "o", "", "Path and Name to an output file")

}

func generateComponentFlags() {
	componentFlags := generateComponentCmd.Flags()

	componentFlags.StringVarP(&componentOpts.CatalogSource, "catalog-source", "c", "", "Catalog source location (local or remote)")
	componentFlags.StringVarP(&componentOpts.Profile, "profile", "p", "", "Profile source location (local or remote)")
	componentFlags.StringSliceVarP(&componentOpts.Requirements, "requirements", "r", []string{}, "List of requirements to capture")
}
