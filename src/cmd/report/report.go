package report

import (
	"errors"
	"os"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/defenseunicorns/lula/src/pkg/message"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"gopkg.in/yaml.v3"
)

type flags struct {
	InputFile  string // -f --input-file
	OutputFile string // -o --output-file
}

var opts = &flags{}

var reportHelp = `

`

var reportCmd = &cobra.Command{
	Use:     "report",
	Hidden:  false,
	Aliases: []string{"r"},
	Short:   "Build a compliance report",
	Example: reportHelp,
	Run: func(_ *cobra.Command, args []string) {
		if opts.InputFile == "" {
			message.Fatal(errors.New("flag input-file is not set"),
				"Please specify an input file with the -f flag")
		}

		oscalModels, err := readOSCALFile(opts.InputFile)
		if err != nil {
			message.Fatal(err, "Failed to unmarshall OSCAL from file")
		}

		modelType, err := determineOSCALModel(oscalModels)
		if err != nil {
			message.Fatal(err, "Unable to determine OSCAL model type")
		}

		// REMOVE: Using for testing if I am able to determine models correctly
		fmt.Println(modelType)

		switch modelType {
		case "catalog":
			fmt.Println("reporting does not create reports for catalogs at this time")
		case "profile":
			fmt.Println("reporting does not create reports for profile at this time")
		case "component-definition":

		case "assessment-plan":
			fmt.Println("reporting does not create reports for assessment plans at this time")
		case "assessment-results":
			fmt.Println("reporting does not create reports for assessment results at this time")
		case "system-security-plan":
			fmt.Println("reporting does not create reports for system security plans at this time")
		case "plan-of-action-and-milestones":
			fmt.Println("reporting does not create reports for plan of action and milestones at this time")
		default:
			message.Fatal(fmt.Errorf("unknown OSCAL model type: %s", modelType), "Failed to process OSCAL file")
		}


	},
}

func ReportCommand() *cobra.Command {

	reportFlags()

	return reportCmd
}

func reportFlags() {
	reportFlags := reportCmd.PersistentFlags()

	reportFlags.StringVarP(&opts.InputFile, "input-file", "f", "", "Path to a manifest file")
	reportFlags.StringVarP(&opts.OutputFile, "output-file", "o", "", "Path and Name to an output file")

}

func readOSCALFile(fileName string ) (oscalTypes_1_1_2.OscalModels, error) {
	var oscalModels oscalTypes_1_1_2.OscalModels

	// Read the YAML file using os.ReadFile
	yamlData, err := os.ReadFile(fileName)
	if err != nil {
		return oscalModels, err
	}

	err = yaml.Unmarshal(yamlData, &oscalModels)
    if err != nil {
        return oscalModels, err
    }

	return oscalModels, nil
}

func determineOSCALModel(oscalModels oscalTypes_1_1_2.OscalModels) (string, error) {
	// Determine which OSCAL model is present by checking non-nil fields
	switch {
	case oscalModels.AssessmentPlan != nil:
		return "assessment-plan", nil
	case oscalModels.AssessmentResults != nil:
		return "assessment-results", nil
	case oscalModels.Catalog != nil:
		return "catalog", nil
	case oscalModels.ComponentDefinition != nil:
		return "component-definition", nil
	case oscalModels.PlanOfActionAndMilestones != nil:
		return "plan-of-action-and-milestones", nil
	case oscalModels.Profile != nil:
		return "profile", nil
	case oscalModels.SystemSecurityPlan != nil:
		return "system-security-plan", nil
	default:
		return "", fmt.Errorf("unable to determine OSCAL model type")
	}
}
