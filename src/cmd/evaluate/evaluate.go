package evaluate

import (
	"fmt"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/files"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/common/result"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

var evaluateHelp = `
To evaluate the latest results in two assessment results files:
	lula evaluate -f assessment-results-threshold.yaml -f assessment-results-new.yaml

To evaluate two results (threshold and latest) in a single OSCAL file:
	lula evaluate -f assessment-results.yaml

To target a specific framework for validation:
	lula evaluate -f assessment-results.yaml --target critical

`

type flags struct {
	InputFile []string // -f --input-file
	Target    string   // -t --target
	summary   bool     // -s --summary
}

var opts = &flags{}

var evaluateCmd = &cobra.Command{
	Use:     "evaluate",
	Short:   "evaluate two results of a Security Assessment Results",
	Long:    "Lula evaluation of Security Assessment Results",
	Example: evaluateHelp,
	Aliases: []string{"eval"},
	Run: func(cmd *cobra.Command, args []string) {

		// Build map of filepath -> assessment results
		assessmentMap, err := readManyAssessmentResults(opts.InputFile)
		if err != nil {
			message.Fatal(err, err.Error())
		}

		EvaluateAssessments(assessmentMap, opts.Target, opts.summary)
	},
}

func EvaluateCommand() *cobra.Command {

	evaluateCmd.Flags().StringSliceVarP(&opts.InputFile, "input-file", "f", []string{}, "Path to the file to be evaluated")
	evaluateCmd.MarkFlagRequired("input-file")
	evaluateCmd.Flags().StringVarP(&opts.Target, "target", "t", "", "the specific control implementations or framework to validate against")
	evaluateCmd.Flags().BoolVarP(&opts.summary, "summary", "s", false, "Print a summary of the evaluation")
	// insert flag options here
	return evaluateCmd
}

func EvaluateAssessments(assessmentMap map[string]*oscalTypes_1_1_2.AssessmentResults, target string, summary bool) {
	// Identify the threshold & latest for comparison
	resultMap := oscal.FilterResults(assessmentMap)

	if target != "" {
		if result, ok := resultMap[target]; ok {
			err := evaluateTarget(result, target, summary)
			if err != nil {
				message.Warn(err.Error())
			}
		}
	} else {
		for source, result := range resultMap {
			err := evaluateTarget(result, source, summary)
			if err != nil {
				message.Warn(err.Error())
			}
		}
	}

	// Write each file back in the case of modification
	for filePath, assessment := range assessmentMap {
		model := oscalTypes_1_1_2.OscalCompleteSchema{
			AssessmentResults: assessment,
		}

		oscal.WriteOscalModel(filePath, &model)
	}
}

func evaluateTarget(target oscal.EvalResult, source string, summary bool) error {
	message.Debugf("Length of results: %d", len(target.Results))
	if len(target.Results) == 0 {
		return fmt.Errorf("no results found")
	}

	if len(target.Results) == 1 {
		// Only one result identified - update to make it the threshold
		oscal.UpdateProps("threshold", oscal.LULA_NAMESPACE, "true", target.Results[0].Props)
		message.Warnf("less than 2 results found for target: %s - no comparison possible", source)
		return nil
	}

	if target.Threshold != nil && target.Latest != nil {
		if target.Threshold.UUID == target.Latest.UUID {
			message.Fatal(fmt.Errorf("cannot compare the same assessment result against itself"), "cannot compare the same assessment result against itself")
		}
		var findingsWithoutObservations []string
		// Compare the assessment results
		spinner := message.NewProgressSpinner("Evaluating Assessment Results %s against %s\n", target.Threshold.UUID, target.Latest.UUID)
		defer spinner.Stop()

		message.Debugf("threshold UUID: %s / latest UUID: %s", target.Threshold.UUID, target.Latest.UUID)

		status, resultComparison, err := oscal.EvaluateResults(target.Threshold, target.Latest)
		if err != nil {
			message.Fatal(err, err.Error())
		}

		// Print summary
		if summary {
			message.Info("Summary of All Observations:")
			findingsWithoutObservations = result.Collapse(resultComparison).PrintObservationComparisonTable(false, true, false)
			if len(findingsWithoutObservations) > 0 {
				message.Warnf("%d Finding(s) Without Observations", len(findingsWithoutObservations))
				message.Info(strings.Join(findingsWithoutObservations, ", "))
			}
		}

		// Check 'status' - Result if evaluation is passing or failing
		// Fails if anything went from satisfied -> not-satisfied OR if any old findings are removed (doesn't matter whether they were satisfied or not)
		if status {
			// Print new-passing-findings
			newSatisfied := resultComparison["new-satisfied"]
			nowSatisfied := resultComparison["now-satisfied"]
			if len(newSatisfied) > 0 || len(nowSatisfied) > 0 {
				message.Info("New passing finding Target-Ids:")
				for id := range newSatisfied {
					message.Infof("%s", id)
				}
				for id := range nowSatisfied {
					message.Infof("%s", id)
				}

				message.Infof("New threshold identified - threshold will be updated to result %s", target.Latest.UUID)

				// Update latest threshold prop
				oscal.UpdateProps("threshold", oscal.LULA_NAMESPACE, "true", target.Latest.Props)
				oscal.UpdateProps("threshold", oscal.LULA_NAMESPACE, "false", target.Threshold.Props)
			} else {
				// retain result as threshold
				oscal.UpdateProps("threshold", oscal.LULA_NAMESPACE, "true", target.Threshold.Props)
			}

			// Print new-not-satisfied
			newFailing := resultComparison["new-not-satisfied"]
			if len(newFailing) > 0 {
				message.Info("New failing finding Target-Ids:")
				for id := range newFailing {
					message.Infof("%s", id)
				}
			}
			spinner.Success()
			message.Info("Evaluation Passed Successfully")
		} else {
			// Print no-longer-satisfied
			message.Warn("Evaluation Failed against the following:")

			// Alternative printing in a single table
			failedFindings := map[string]result.ResultComparisonMap{
				"no-longer-satisfied":   resultComparison["no-longer-satisfied"],
				"removed-satisfied":     resultComparison["removed-satisfied"],
				"removed-not-satisfied": resultComparison["removed-not-satisfied"],
			}
			findingsWithoutObservations = result.Collapse(failedFindings).PrintObservationComparisonTable(true, false, true)
			// handle controls that failed but didn't have observations
			if len(findingsWithoutObservations) > 0 {
				message.Warnf("%d Failed Finding(s) Without Observations", len(findingsWithoutObservations))
				message.Info(strings.Join(findingsWithoutObservations, ", "))
			}

			message.Fatalf(fmt.Errorf("failed to meet established threshold"), "failed to meet established threshold")

			// retain result as threshold
			oscal.UpdateProps("threshold", oscal.LULA_NAMESPACE, "true", target.Threshold.Props)
		}

		spinner.Success()

	} else if target.Threshold == nil {
		message.Fatal(fmt.Errorf("no threshold assessment results could be identified"), "no threshold assessment results could be identified")
	}

	return nil
}

// Read many filepaths into a map[filepath]*AssessmentResults
// Placing here until otherwise decided on value elsewhere
func readManyAssessmentResults(fileArray []string) (map[string]*oscalTypes_1_1_2.AssessmentResults, error) {
	if len(fileArray) == 0 {
		return nil, fmt.Errorf("no files provided for evaluation")
	}

	assessmentMap := make(map[string]*oscalTypes_1_1_2.AssessmentResults)
	for _, fileString := range fileArray {
		err := files.IsJsonOrYaml(fileString)
		if err != nil {
			return nil, fmt.Errorf("invalid file extension: %s, requires .json or .yaml", fileString)
		}

		data, err := common.ReadFileToBytes(fileString)
		if err != nil {
			return nil, err
		}
		assessment, err := oscal.NewAssessmentResults(data)
		if err != nil {
			return nil, err
		}
		assessmentMap[fileString] = assessment
	}

	return assessmentMap, nil
}
