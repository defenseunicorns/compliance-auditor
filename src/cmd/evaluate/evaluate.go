package evaluate

import (
	"errors"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-1"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/spf13/cobra"
)

var evaluateHelp = `
To evaluate the latest results in two assessment results files:
	lula evaluate -f assessment-results-1.yaml -f assessment-results-2.yaml

To evaluate two results (latest and preceding) in a single assessment results file:
	lula evaluate -f assessment-results.yaml
`

type flags struct {
	files []string
}

var opts = &flags{}

var evaluateCmd = &cobra.Command{
	Use:     "evaluate",
	Short:   "evaluate two results of a Security Assessment Results",
	Long:    "Lula evaluation of Security Assessment Results",
	Example: evaluateHelp,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Access the files and evaluate them
		err := EvaluateAssessmentResults(opts.files)
		if err != nil {
			return err
		}
		return nil
	},
}

func EvaluateCommand() *cobra.Command {

	evaluateCmd.Flags().StringArrayVarP(&opts.files, "file", "f", []string{}, "Path to the file to be evaluated")
	// insert flag options here
	return evaluateCmd
}

func EvaluateAssessmentResults(files []string) error {
	// Read in files - establish the results to
	if len(files) == 0 {
		// TODO: Determine if we will handle a default location/name for assessment files
		return errors.New("No files provided")
	} else if len(files) == 1 {
		data, err := common.ReadFileToBytes(files[0])
		if err != nil {
			return err
		}
		assessment, err := oscal.NewAssessmentResults(data)
		if err != nil {
			return err
		}
		status, findings, err := EvaluateResults(assessment.Results[0], assessment.Results[1])
	} else if len(files) == 2 {

	} else {
		return errors.New("Exceeded maximum of 2 files for evaluation")
	}
	// Create assessmentResults objects
	// Identify the results objects to evaluate
	// Evaluate the results objects
	// Print the results
	return nil
}

func EvaluateResults(oldResult oscalTypes.Result, newResult oscalTypes.Result) (bool, []oscalTypes.Finding, error) {
	// Store findings for review here
	findings := make([]oscalTypes.Finding, 0)

	return false, nil
}
