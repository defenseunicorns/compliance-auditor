package dev

import (
	"fmt"
	"os"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

type printResourcesFlags struct {
	assessment      string // -a --assessment
	observationUuid string // --observation-uuid
	outputFile      string // -o --output-file
}

var printResourcesOpts = &printResourcesFlags{}

var printResourcesHelp = `
To print resources from lula validation manifest:
	lula dev print-resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid> --result-type <result-type=threshold|latest> --output-file <output-file>
`

var printResourcesCmdLong = `
Print out the the JSON resources used by a observation, which provided the input to the Lula Validation
`

func init() {
	printResourcesCmd := &cobra.Command{
		Use:     "print-resources",
		Short:   "Print Resources from a Lula Validation execution",
		Long:    printResourcesCmdLong,
		Example: printResourcesHelp,
		Run: func(cmd *cobra.Command, args []string) {
			assessmentData, err := os.ReadFile(printResourcesOpts.assessment)
			if err != nil {
				message.Fatalf(nil, "Invalid assessment file: %v", err)
			}

			oscalAssessment, err := oscal.NewAssessmentResults(assessmentData)
			if err != nil {
				message.Fatalf(nil, "Invalid assessment results model: %v", err)
			}

			printResources(oscalAssessment, printResourcesOpts.observationUuid, printResourcesOpts.outputFile)
		},
	}

	devCmd.AddCommand(printResourcesCmd)

	printResourcesCmd.Flags().StringVar(&printResourcesOpts.assessment, "assessment", "", "the path to an assessment-results file")
	printResourcesCmd.Flags().StringVar(&printResourcesOpts.observationUuid, "observation-uuid", "", "the observation uuid")
	printResourcesCmd.Flags().StringVarP(&printResourcesOpts.outputFile, "output-file", "o", "", "the path to write the resources json")
	printResourcesCmd.MarkFlagRequired("assessment")
	printResourcesCmd.MarkFlagRequired("observation-uuid")
}

func printResources(assessment *oscalTypes_1_1_2.AssessmentResults, observationUuid, outputFile string) {
	if assessment == nil {
		message.Fatalf(nil, "Assessment results is nil")
	}

	// Get all observations from the assessment results
	observationMap := make(map[string]oscalTypes_1_1_2.Observation)
	for _, result := range assessment.Results {
		if result.Observations != nil {
			for _, observation := range *result.Observations {
				observationMap[observation.UUID] = observation
			}
		}
	}

	observation, found := observationMap[observationUuid]
	if !found {
		message.Fatalf(nil, "Observation not found in assessment results")
	}

	// Get the resources from the assessment backmatter
	resourceMap := make(map[string]string)
	if assessment.BackMatter != nil {
		resourceMap = oscal.BackMatterToMap(*assessment.BackMatter)
	}

	if observation.Links == nil {
		message.Fatalf(nil, "Observation does not contain resources")
	}
	var resourceId string
	for _, link := range *observation.Links {
		if link.Rel == "lula.resources" {
			resourceId = common.TrimIdPrefix(link.Href)
			break
		}
	}

	resources, found := resourceMap[resourceId]
	if !found {
		message.Fatalf(nil, "Resource not found in assessment backmatter")
	}

	// Write the resources to a file
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(resources), 0644)
		if err != nil {
			message.Fatalf(err, "Error writing resource JSON to file: %v", err)
		}
	} else {
		// Else print to stdout
		fmt.Println(resources)
	}
}
