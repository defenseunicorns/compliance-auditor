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

type printValidationFlags struct {
	component       string // -c --component
	assessment      string // -a --assessment
	observationUuid string // -o --observation-uuid
}

var printValidationOpts = &printValidationFlags{}

var printValidationHelp = `
To print a specific lula validation that generated a given observation:
	lula dev print-validation --component /path/to/component.yaml --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>
`

var printValidationCmdLong = `
Print out the the Lula Validation that yielded the provided observation
`

func init() {
	printValidationCmd := &cobra.Command{
		Use:     "print-validation",
		Short:   "Print Lula Validation",
		Long:    printResourcesCmdLong,
		Example: printResourcesHelp,
		Run: func(cmd *cobra.Command, args []string) {
			// Get the component, assessment, validation ID, and result type to lookup the resources
			componentData, err := os.ReadFile(printValidationOpts.component)
			if err != nil {
				message.Fatalf(nil, "Invalid component file: %v", err)
			}

			assessmentData, err := os.ReadFile(printValidationOpts.assessment)
			if err != nil {
				message.Fatalf(nil, "Invalid assessment file: %v", err)
			}

			oscalComponent, err := oscal.NewOscalComponentDefinition(componentData)
			if err != nil {
				message.Fatalf(nil, "Invalid component definition model: %v", err)
			}

			oscalAssessment, err := oscal.NewAssessmentResults(assessmentData)
			if err != nil {
				message.Fatalf(nil, "Invalid assessment results model: %v", err)
			}

			printValidation(oscalComponent, oscalAssessment, printValidationOpts.observationUuid)

		},
	}

	devCmd.AddCommand(printValidationCmd)
	printValidationCmd.Flags().StringVar(&printValidationOpts.component, "component", "", "the path to a validation manifest file")
	printValidationCmd.Flags().StringVar(&printValidationOpts.assessment, "assessment", "", "the path to an assessment-results file")
	printValidationCmd.Flags().StringVar(&printValidationOpts.observationUuid, "observation-uuid", "", "the observation uuid")
	printValidationCmd.MarkFlagRequired("component")
	printValidationCmd.MarkFlagRequired("assessment")
	printValidationCmd.MarkFlagRequired("observation-uuid")
}

func printValidation(component *oscalTypes_1_1_2.ComponentDefinition, assessment *oscalTypes_1_1_2.AssessmentResults, observationUuid string) {
	if component == nil {
		message.Fatalf(nil, "Component definition is nil")
	}

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

	found, validationUuid := oscal.GetProp("validation", oscal.LULA_NAMESPACE, observation.Props)
	if !found {
		message.Fatalf(nil, "No validation linked to observation")
	}

	// Get all validations from the component definition
	resourceMap := make(map[string]string)
	if component.BackMatter != nil {
		resourceMap = oscal.BackMatterToMap(*component.BackMatter)
	}

	trimmedId := common.TrimIdPrefix(validationUuid)

	// Find the validation in the map
	validation, found := resourceMap[trimmedId]
	if !found {
		message.Fatalf(nil, "Validation not found in component definition")
	}

	// Print the validation
	fmt.Println(validation)
}
