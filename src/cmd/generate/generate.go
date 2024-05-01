package generate

import (
	"bytes"
	"fmt"
	"os"

	gooscalUtils "github.com/defenseunicorns/go-oscal/src/pkg/utils"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"

	"gopkg.in/yaml.v3"
)

type flags struct {
	InputFile  string // -f --input-file
	OutputFile string // -o --output-file
}

type componentFlags struct {
	flags
	CatalogSource string   // -c --catalog
	Profile       string   // -p --profile
	Component     string   // --component
	Requirements  []string // -r --requirements
	Remarks       []string // --remarks
}

var opts = &flags{}
var componentOpts = &componentFlags{}

// Base generate command could handle a large E2E focused generation that is driven from various artifacts.
// This will include some prerequisites such as component-definitions - but will ultimately lead to generation
// of the SSP / SAP / POAM / and maybe SAR in a maintainable way.
var generateCmd = &cobra.Command{
	Use:     "generate",
	Hidden:  false, // Hidden for now until fully implemented
	Aliases: []string{"g", "gen"},
	Short:   "Generate a specified compliance artifact template",
}

var componentHelp = `
To generate a new component-definition template:
lula generate component -c <catalog source url> -r control-a,control-b,control-c
- IE lula generate component -c https://raw.githubusercontent.com/usnistgov/oscal-content/master/nist.gov/SP800-53/rev5/json/NIST_SP-800-53_rev5_catalog.json -r ac-1,ac-2,au-5

To Generate and merge with an existing Component Definition:
lula generate component -c <catalog source url> -r control-a,control-b,control-c -o existing-component.yaml

To Generate a component definition with a specific "named" component:
lula generate component -c <catalog source url> -r control-a --component "Software X"

To Generate a component definition with remarks populated from specific control "parts":
lula generate component -c <catalog source url> -r control-a --remarks guidance,assessment-objective
`

// Component-Definition generation will generate an OSCAL file that can be used both as the basis for Lula validations
// as well as required components for SSP/SAP/SAR/POAM.
var generateComponentCmd = &cobra.Command{
	Use:     "component",
	Aliases: []string{"c"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Generate a component definition OSCAL template",
	Example: componentHelp,
	Run: func(_ *cobra.Command, args []string) {
		var remarks []string
		var title = "Component Title"
		var outputFile = "oscal-component.yaml"
		// check for Catalog Source - this field is required
		if componentOpts.CatalogSource == "" {
			message.Fatal(fmt.Errorf("no catalog source provided"), "generate component requires a catalog input source")
		}

		// Assign remarks from flag or default to "statement"
		if len(componentOpts.Remarks) == 0 {
			remarks = []string{"statement"}
		} else {
			remarks = componentOpts.Remarks
		}

		// Assign the component title from flag or default to "Component Title"
		if componentOpts.Component != "" {
			title = componentOpts.Component
		}

		source := componentOpts.CatalogSource

		// Fetch the catalog source
		data, err := network.Fetch(source)
		if err != nil {
			message.Fatalf(fmt.Errorf("error fetching catalog source"), "error fetching catalog source")
		}

		// Create new catalog object
		catalog, err := oscal.NewCatalog(source, data)
		if err != nil {
			message.Fatalf(fmt.Errorf("error creating catalog"), "error creating catalog")
		}

		// Create a component definition from the catalog given required context
		comp, err := oscal.ComponentFromCatalog(source, catalog, title, componentOpts.Requirements, remarks)
		if err != nil {
			message.Fatalf(fmt.Errorf("error creating component"), "error creating component")
		}

		var fileName string
		if opts.OutputFile != "" {
			fileName = opts.OutputFile
		} else {
			fileName = "oscal-component.yaml"
		}

		if componentOpts.OutputFile != "" {
			outputFile = componentOpts.OutputFile
		}
		// Check for an existing file
		if _, err := os.Stat(outputFile); err == nil {
			existingFileBytes, err := os.ReadFile(outputFile)
			if err != nil {
				message.Fatalf(fmt.Errorf("error reading existing file"), "error reading existing file")
			}
			// Create new component definition object
			existingComponent, err := oscal.NewOscalComponentDefinition(outputFile, existingFileBytes)
			if err != nil {
				message.Fatalf(fmt.Errorf("error creating new component definition"), "error creating new component definition")
			}
			// Merge the newly generated document into the existing document
			comp, err = oscal.MergeComponentDefinitions(existingComponent, comp)
			if err != nil {
				message.Fatalf(fmt.Errorf("error merging component definition on component"), "error merging component definition on component")
			}
			comp.Metadata.LastModified = gooscalUtils.GetTimestamp()
		}

		var b bytes.Buffer

		var component = oscalTypes_1_1_2.OscalModels{
			ComponentDefinition: &comp,
		}

		yamlEncoder := yaml.NewEncoder(&b)
		yamlEncoder.SetIndent(2)
		yamlEncoder.Encode(component)

		message.Infof("Writing Component Definition to: %s", fileName)

		err = os.WriteFile(fileName, b.Bytes(), 0644)
		if err != nil {
			message.Fatalf(fmt.Errorf("error writing Component Definition to: %s", fileName), "error writing Component Definition to: %s", fileName)
		}

	},
}

// var generateAssessmentPlanCmd = &cobra.Command{
// 	Use:     "assessment-plan",
// 	Aliases: []string{"ap"},
// 	Args:    cobra.MaximumNArgs(1),
// 	Short:   "Generate an assessment plan OSCAL template",
// 	Run: func(_ *cobra.Command, args []string) {
// 		message.Info("generate assessment-plan executed")

// 		// For each component-definition in array of component-definitions
// 		// Read component-definition
// 		// Collect all implemented-requirements
// 		// Collect all items from the backmatter
// 		// Create new assessment-plan object
// 		// Transfer to assessment-plan.reviewed-controls?
// 	},
// }

// var generateSystemSecurityPlanCmd = &cobra.Command{
// 	Use:     "system-security-plan",
// 	Aliases: []string{"ssp"},
// 	Args:    cobra.MaximumNArgs(1),
// 	Short:   "Generate a system security plan OSCAL template",
// 	Run: func(_ *cobra.Command, args []string) {
// 		message.Info("generate system-security-plan executed")

// 		// For each component-definition in array of component-definitions
// 		// Read component-definition
// 		// Collect all implemented-requirements
// 		// aggregate by control-id
// 		// export to system-security-plan.control-implementation.implemented-requirements
// 	},
// }

// var generatePOAMCmd = &cobra.Command{
// 	Use:     "plan-of-action-and-milestones",
// 	Aliases: []string{"poam"},
// 	Args:    cobra.MaximumNArgs(1),
// 	Short:   "Generate a plan of actions and milestones OSCAL template",
// 	Run: func(_ *cobra.Command, args []string) {
// 		message.Info("generate plan-of-action-and-milestones executed")

// 		// Locate an assessment-results artifact
// 		// Create an assessment-results object
// 		// Transfer 'not-satisfied' findings and observations to poam.findings and poam.observations as appropriate
// 	},
// }

func GenerateCommand() *cobra.Command {

	generateCmd.AddCommand(generateComponentCmd)
	// generateCmd.AddCommand(generateAssessmentPlanCmd)
	// generateCmd.AddCommand(generateSystemSecurityPlanCmd)
	// generateCmd.AddCommand(generatePOAMCmd)

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
	componentFlags.StringVar(&componentOpts.Component, "component", "", "Component Title")
	componentFlags.StringSliceVarP(&componentOpts.Requirements, "requirements", "r", []string{}, "List of requirements to capture")
	componentFlags.StringSliceVar(&componentOpts.Remarks, "remarks", []string{}, "Target for remarks population (default = statement)")
}
