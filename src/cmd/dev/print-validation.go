package dev

import (
	"fmt"
	"os"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/spf13/cobra"

	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

var printValidationHelp = `
To print a specific lula validation that generated a given observation:
	lula dev print-validation --component /path/to/component.yaml --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>
`

var printValidationCmdLong = `
Prints the Lula Validation from a specified observation. Assumes that the validation is in the back matter of the provided component definition.
`

func PrintValidationCommand() *cobra.Command {
	var (
		component       string // -c --component
		assessment      string // -a --assessment
		observationUuid string // -u --observation-uuid
		outputFile      string // -o --output-file
	)

	printValidationCmd := &cobra.Command{
		Use:     "print-validation",
		Short:   "Print Lula Validation",
		Long:    printValidationCmdLong,
		Example: printValidationHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the component, assessment, validation ID, and result type to lookup the resources
			componentData, err := os.ReadFile(component)
			if err != nil {
				return fmt.Errorf("invalid component file: %v", err)
			}

			oscalComponent, err := oscal.NewOscalComponentDefinition(componentData)
			if err != nil {
				return fmt.Errorf("error creating oscal component definition model: %v", err)
			}

			assessmentData, err := os.ReadFile(assessment)
			if err != nil {
				return fmt.Errorf("invalid assessment file: %v", err)
			}

			oscalAssessment, err := oscal.NewAssessmentResults(assessmentData)
			if err != nil {
				return fmt.Errorf("error creating oscal assessment results model: %v", err)
			}

			err = PrintValidation(oscalComponent, oscalAssessment, observationUuid, outputFile)
			if err != nil {
				return fmt.Errorf("error printing validation: %v", err)
			}

			return nil
		},
	}

	printValidationCmd.Flags().StringVarP(&component, "component", "c", "", "the path to a validation manifest file")
	printValidationCmd.Flags().StringVarP(&assessment, "assessment", "a", "", "the path to an assessment-results file")
	printValidationCmd.Flags().StringVarP(&observationUuid, "observation-uuid", "u", "", "the observation uuid")
	printValidationCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to write the resources json")

	err := printValidationCmd.MarkFlagRequired("component")
	if err != nil {
		message.Fatal(err, "error initializing print-resources command flag: assessment")
	}
	err = printValidationCmd.MarkFlagRequired("assessment")
	if err != nil {
		message.Fatal(err, "error initializing print-resources command flag: assessment")
	}
	err = printValidationCmd.MarkFlagRequired("observation-uuid")
	if err != nil {
		message.Fatal(err, "error initializing required command flag: observation-uuid")
	}

	return printValidationCmd
}

func init() {
	devCmd.AddCommand(PrintValidationCommand())
}

func PrintValidation(component *oscalTypes_1_1_2.ComponentDefinition, assessment *oscalTypes_1_1_2.AssessmentResults, observationUuid, outputFile string) error {
	if component == nil {
		return fmt.Errorf("component definition is nil")
	}

	if assessment == nil {
		return fmt.Errorf("assessment results is nil")
	}

	// Get the observation
	observation, err := GetObservationByUuid(assessment, observationUuid)
	if err != nil {
		return err
	}

	// Get the validation
	found, validationUuid := oscal.GetProp("validation", oscal.LULA_NAMESPACE, observation.Props)
	if !found {
		return fmt.Errorf("no validation linked to observation")
	}

	// Find validation ID in the component definition back matter
	resourceMap := make(map[string]string)
	if component.BackMatter != nil {
		resourceMap = oscal.BackMatterToMap(*component.BackMatter)
	}

	trimmedId := common.TrimIdPrefix(validationUuid)

	// Find the validation in the map
	validation, found := resourceMap[trimmedId]
	if !found {
		return fmt.Errorf("validation not found in component definition")
	}

	// Print the validation
	if outputFile == "" {
		message.Printf(validation)
	} else {
		err = os.WriteFile(outputFile, []byte(validation), 0600)
		if err != nil {
			return fmt.Errorf("error writing validation to file: %v", err)
		}
	}
	return nil
}
