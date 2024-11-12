package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/spf13/cobra"

	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

var printHelp = `
To print resources from lula validation manifest:
	lula tools print --resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>

To print resources from lula validation manifest to output file:
	lula tools print --resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid> --output-file /path/to/output.json

To print the lula validation that generated a given observation:
	lula tools print --validation --component /path/to/component.yaml --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>
`

var printCmdLong = `
Prints out data about an OSCAL Observation from the OSCAL Assessment Results model. 
Given "--resources", the command will print the JSON resources input that were provided to a Lula Validation, as identified by a given observation and assessment results file. 
Given "--validation", the command will print the Lula Validation that generated a given observation, as identified by a given observation, assessment results file, and component definition file.
`

func PrintCommand() *cobra.Command {
	var (
		resources       bool   // -r --resources
		validation      bool   // -v --validation
		assessment      string // -a --assessment
		observationUuid string // -u --observation-uuid
		outputFile      string // -o --output-file
		component       string // -c --component
	)

	printCmd := &cobra.Command{
		Use:     "print",
		Short:   "Print Resources or Lula Validation from an Assessment Observation",
		Long:    printCmdLong,
		Example: printHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			assessmentData, err := common.ReadFileToBytes(assessment)
			if err != nil {
				return fmt.Errorf("invalid assessment file: %v", err)
			}

			assessmentDir, err := filepath.Abs(filepath.Dir(assessment))
			if err != nil {
				return fmt.Errorf("error getting assessment directory: %v", err)
			}

			oscalAssessment, err := oscal.NewAssessmentResults(assessmentData)
			if err != nil {
				return fmt.Errorf("error creating oscal assessment results model: %v", err)
			}

			// Print the resources or validation
			if resources {
				err = PrintResources(oscalAssessment, observationUuid, assessmentDir, outputFile)
				if err != nil {
					return fmt.Errorf("error printing resources: %v", err)
				}
			} else if validation {
				// Compose the component definition
				composer, err := composition.New(composition.WithModelFromLocalPath(component))
				if err != nil {
					return fmt.Errorf("error creating new composer: %v", err)
				}
				oscalModel, err := composer.ComposeFromPath(cmd.Context(), component)
				if err != nil {
					return fmt.Errorf("error composing model: %v", err)
				}

				// Print the validation
				err = PrintValidation(oscalModel.ComponentDefinition, oscalAssessment, observationUuid, outputFile)
				if err != nil {
					return fmt.Errorf("error printing validation: %v", err)
				}
			}
			return nil
		},
	}

	// Add flags, set logic for flag behavior
	printCmd.Flags().BoolVarP(&resources, "resources", "r", false, "true if the user is printing resources")
	printCmd.Flags().BoolVarP(&validation, "validation", "v", false, "true if the user is printing validation")
	printCmd.MarkFlagsMutuallyExclusive("resources", "validation")

	printCmd.Flags().StringVarP(&assessment, "assessment", "a", "", "the path to an assessment-results file")
	err := printCmd.MarkFlagRequired("assessment")
	if err != nil {
		message.Fatal(err, "error initializing print-resources command flag: assessment")
	}

	printCmd.Flags().StringVarP(&observationUuid, "observation-uuid", "u", "", "the observation uuid")
	err = printCmd.MarkFlagRequired("observation-uuid")
	if err != nil {
		message.Fatal(err, "error initializing required command flag: observation-uuid")
	}

	printCmd.Flags().StringVarP(&component, "component", "c", "", "the path to a validation manifest file")
	printCmd.MarkFlagsRequiredTogether("validation", "component")

	printCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to write the resources json")

	return printCmd
}

func init() {
	toolsCmd.AddCommand(PrintCommand())
}

func PrintResources(assessment *oscalTypes_1_1_2.AssessmentResults, observationUuid, assessmentDir, outputFile string) error {
	if assessment == nil {
		return fmt.Errorf("assessment is nil")
	}

	observation, err := oscal.GetObservationByUuid(assessment, observationUuid)
	if err != nil {
		return err
	}

	// Get the resources from the remote reference
	// TODO: will an observation ever have multiple resource links?
	resourceCount := 0
	var resource types.DomainResources

	if observation.Links == nil {
		return fmt.Errorf("observation does not contain a remote reference")
	}

	for _, link := range *observation.Links {
		if link.Rel == "lula.resources" {
			resourceCount++
			if resourceCount > 1 {
				return fmt.Errorf("observation contains multiple remote references, only the first printed")
			}

			resourceData, err := network.Fetch(link.Href, network.WithBaseDir(assessmentDir))
			if err != nil {
				return fmt.Errorf("error fetching resource: %v", err)
			}

			err = json.Unmarshal(resourceData, &resource)
			if err != nil {
				return fmt.Errorf("error unmarshalling resource: %v", err)
			}
		}
	}

	// Write the resources to a file if found
	err = types.WriteResources(resource, outputFile)
	if err != nil {
		return err
	}

	return nil
}

func PrintValidation(component *oscalTypes_1_1_2.ComponentDefinition, assessment *oscalTypes_1_1_2.AssessmentResults, observationUuid, outputFile string) error {
	if component == nil {
		return fmt.Errorf("component definition is nil")
	}

	if assessment == nil {
		return fmt.Errorf("assessment results is nil")
	}

	// Get the observation
	observation, err := oscal.GetObservationByUuid(assessment, observationUuid)
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
		message.Printf("%s", validation)
	} else {
		err = os.WriteFile(outputFile, []byte(validation), 0600)
		if err != nil {
			return fmt.Errorf("error writing validation to file: %v", err)
		}
	}
	return nil
}
