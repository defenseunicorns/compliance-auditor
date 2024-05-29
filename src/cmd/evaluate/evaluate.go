package evaluate

import (
	"fmt"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

var evaluateHelp = `
To evaluate the latest results in two assessment results files:
	lula evaluate -f assessment-results-threshold.yaml -f assessment-results-new.yaml

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
	Aliases: []string{"eval"},
	Run: func(cmd *cobra.Command, args []string) {

		// Access the files and evaluate them
		err := EvaluateAssessmentResults(opts.files)
		if err != nil {
			message.Fatal(err, err.Error())
		}
	},
}

func EvaluateCommand() *cobra.Command {

	evaluateCmd.Flags().StringArrayVarP(&opts.files, "file", "f", []string{}, "Path to the file to be evaluated")
	// insert flag options here
	return evaluateCmd
}

func EvaluateAssessmentResults(fileArray []string) error {
	var status bool
	var findings map[string][]oscalTypes_1_1_2.Finding

	// Read in files - establish the results to
	if len(fileArray) == 0 {
		// TODO: Determine if we will handle a default location/name for assessment files
		return fmt.Errorf("No files provided for evaluation")
	}

	for _, f := range fileArray {
		err := files.IsJsonOrYaml(f)
		if err != nil {
			return fmt.Errorf("invalid file extension: %s, requires .json or .yaml", f)
		}
	}

	if len(fileArray) == 1 {
		data, err := common.ReadFileToBytes(fileArray[0])
		if err != nil {
			return err
		}
		assessment, err := oscal.NewAssessmentResults(data)
		if err != nil {
			return err
		}
		if len(assessment.Results) < 2 {
			return fmt.Errorf("2 or more result objects must be present for evaluation\n")
		}
		// We write results to the assessment-results report in newest -> oldest
		// Older being our threshold here
		status, findings, err = EvaluateResults(&assessment.Results[1], &assessment.Results[0])
		if err != nil {
			return err
		}

	} else if len(fileArray) == 2 {
		data, err := common.ReadFileToBytes(fileArray[0])
		if err != nil {
			return err
		}
		assessmentOne, err := oscal.NewAssessmentResults(data)
		if err != nil {
			return err
		}
		data, err = common.ReadFileToBytes(fileArray[1])
		if err != nil {
			return err
		}
		assessmentTwo, err := oscal.NewAssessmentResults(data)
		if err != nil {
			return err
		}

		// Consider parsing the timestamps for comparison
		// Older timestamp being the threshold

		status, findings, err = EvaluateResults(&assessmentOne.Results[0], &assessmentTwo.Results[0])
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Exceeded maximum of 2 files for evaluation\n")
	}

	if status {
		message.Info("Evaluation Passing the established threshold")
		if len(findings["new-findings"]) > 0 {
			message.Info("New finding Target-Ids:")
			for _, finding := range findings["new-findings"] {
				message.Infof("%s", finding.Target.TargetId)
			}
		}
		return nil
	} else {
		message.Warn("Evaluation Failed against the following findings:")
		for _, finding := range findings["no-longer-satisfied"] {
			message.Warnf("%s", finding.Target.TargetId)
		}
		return fmt.Errorf("Failed to meet established threshold")
	}
}

func EvaluateResults(thresholdResult *oscalTypes_1_1_2.Result, newResult *oscalTypes_1_1_2.Result) (bool, map[string][]oscalTypes_1_1_2.Finding, error) {
	if thresholdResult == nil || thresholdResult.Findings == nil || newResult == nil || newResult.Findings == nil {
		return false, nil, fmt.Errorf("results must contain findings to evaluate")
	}

	spinner := message.NewProgressSpinner("Evaluating Assessment Results %s against %s", newResult.UUID, thresholdResult.UUID)
	defer spinner.Stop()

	// Store unique findings for review here
	findings := make(map[string][]oscalTypes_1_1_2.Finding, 0)
	result := true

	findingMapThreshold := oscal.GenerateFindingsMap(*thresholdResult.Findings)
	findingMapNew := oscal.GenerateFindingsMap(*newResult.Findings)

	// For a given oldResult - we need to prove that the newResult implements all of the oldResult findings/controls
	// We are explicitly iterating through the findings in order to collect a delta to display

	for targetId, finding := range findingMapThreshold {
		if _, ok := findingMapNew[targetId]; !ok {
			// If the new result does not contain the finding of the old result
			// set result to fail, add finding to the findings map and continue
			result = false
			findings[targetId] = append(findings["no-longer-satisfied"], finding)
		} else {
			// If the finding is present in each map - we need to check if the state has changed from "not-satisfied" to "satisfied"
			if finding.Target.Status.State == "satisfied" {
				// Was previously satisfied - compare state
				if findingMapNew[targetId].Target.Status.State == "not-satisfied" {
					// If the new finding is now not-satisfied - set result to false and add to findings
					result = false
					findings["no-longer-satisfied"] = append(findings["no-longer-satisfied"], finding)
				}
			}
			delete(findingMapNew, targetId)
		}
	}

	// All remaining findings in the new map are new findings
	for _, finding := range findingMapNew {
		findings["new-findings"] = append(findings["new-findings"], finding)
	}

	spinner.Success()
	return result, findings, nil
}
