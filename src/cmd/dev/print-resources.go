package dev

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/spf13/cobra"

	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

var printResourcesHelp = `
To print resources from lula validation manifest:
	lula dev print-resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid>

To print resources from lula validation manifest to output file:
	lula dev print-resources --assessment /path/to/assessment.yaml --observation-uuid <observation-uuid> --output-file /path/to/output.json
`

var printResourcesCmdLong = `
Print out the JSON resources input that were provided to a Lula Validation, as identified by a given observation and assessment results file.
`

func PrintResourcesCommand() *cobra.Command {
	var (
		assessment      string // -a --assessment
		observationUuid string // -u --observation-uuid
		outputFile      string // -o --output-file
	)

	printResourcesCmd := &cobra.Command{
		Use:     "print-resources",
		Short:   "Print Resources from a Lula Validation evaluation",
		Long:    printResourcesCmdLong,
		Example: printResourcesHelp,
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

			err = PrintResources(oscalAssessment, observationUuid, assessmentDir, outputFile)
			if err != nil {
				return fmt.Errorf("error printing resources: %v", err)
			}
			return nil
		},
	}

	printResourcesCmd.Flags().StringVarP(&assessment, "assessment", "a", "", "the path to an assessment-results file")
	printResourcesCmd.Flags().StringVarP(&observationUuid, "observation-uuid", "u", "", "the observation uuid")
	printResourcesCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the path to write the resources json")
	err := printResourcesCmd.MarkFlagRequired("assessment")
	if err != nil {
		message.Fatal(err, "error initializing print-resources command flag: assessment")
	}
	err = printResourcesCmd.MarkFlagRequired("observation-uuid")
	if err != nil {
		message.Fatal(err, "error initializing required command flag: observation-uuid")
	}

	return printResourcesCmd
}

func init() {
	devCmd.AddCommand(PrintResourcesCommand())
}

func PrintResources(assessment *oscalTypes_1_1_2.AssessmentResults, observationUuid, assessmentDir, outputFile string) error {
	if assessment == nil {
		return fmt.Errorf("assessment is nil")
	}

	observation, err := GetObservationByUuid(assessment, observationUuid)
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
	err = writeResources(resource, outputFile)
	if err != nil {
		return err
	}

	return nil
}
