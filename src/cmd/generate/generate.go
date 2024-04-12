package generate

import (
	"bytes"
	"fmt"
	"os"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"

	"gopkg.in/yaml.v3"
)

// The generate command is expected to have multiple different use-cases based upon
// the artifact to be generated. Some opinionation will be involved - such as OSCAL
// for the supported OSCAL models. This leaves generate open to generating other
// artifacts in the future that support assessment and accreditation processes.

var inputFile, outputFile string // Used to capture a future/potential manifest file for input
type flags struct {
	AssessmentFile string   // -a --assessment-file
	InputFile      string   // -f --input-file
	CatalogSource  string   // -c --catalog
	Profile        string   // -p --profile
	Requirements   []string // -r --requirements
}

var opts = &flags{}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Hidden:  false, // Hidden for now until fully implemented
	Aliases: []string{"g", "gen"},
	Short:   "Generate a specified compliance artifact template",
}

var generateComponentCmd = &cobra.Command{
	Use:     "component",
	Aliases: []string{"c"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Generate a component definition OSCAL template",
	Run: func(_ *cobra.Command, args []string) {
		message.Info("generate component executed")

		// check for inputFile flag content
		if opts.CatalogSource == "" {
			message.Fatal(fmt.Errorf("no catalog source provided"), "generate component requires a catalog input source")
		}

		path := "test/NIST_SP-800-53_rev5_catalog.json"

		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			message.Fatalf(fmt.Errorf("catalog source file not found"), "generate component requires a valid catalog path")
		}
		data, err := os.ReadFile(path)
		if err != nil {
			message.Fatalf(fmt.Errorf("error reading file"), "error reading file")
		}

		catalog, err := oscal.NewCatalog(data)

		comp, err := oscal.ComponentFromCatalog(opts.CatalogSource, catalog, opts.Requirements)

		// marshall to yaml
		if err != nil {
			message.Fatalf(fmt.Errorf("error creating component"), "error creating component")
		}

		fileName := "new-catalog.yaml"

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

	generateFlags.StringVarP(&inputFile, "input-file", "f", "", "Path to a manifest file")
	generateFlags.StringVarP(&outputFile, "output-file", "o", "", "Path and Name to an output file")

}

func generateComponentFlags() {
	componentFlags := generateComponentCmd.Flags()

	componentFlags.StringVarP(&opts.CatalogSource, "catalog-source", "c", "", "Catalog source location (local or remote)")
	componentFlags.StringVarP(&opts.Profile, "profile", "p", "", "Profile source location (local or remote)")
	componentFlags.StringSliceVarP(&opts.Requirements, "requirements", "r", []string{}, "List of requirements to capture")
}
