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

To evaluate two results (threshold and latest) in a single OSCAL file:
	lula evaluate -f assessment-results.yaml
`

type flags struct {
	InputFile []string // -f --input-file
	Target    string   // -t --target
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

		EvaluateAssessments(assessmentMap, opts.Target)
	},
}

func EvaluateCommand() *cobra.Command {

	evaluateCmd.Flags().StringSliceVarP(&opts.InputFile, "input-file", "f", []string{}, "Path to the file to be evaluated")
	evaluateCmd.MarkFlagRequired("input-file")
	evaluateCmd.Flags().StringVarP(&opts.Target, "target", "t", "", "the specific control implementations or framework to validate against")
	// insert flag options here
	return evaluateCmd
}

func EvaluateAssessments(assessmentMap map[string]*oscalTypes_1_1_2.AssessmentResults, target string) {
	// Identify the threshold & latest for comparison
	resultMap := oscal.FilterResults(assessmentMap)

	if target != "" {
		if result, ok := resultMap[target]; ok {
			err := evaluateTarget(result)
			if err != nil {
				message.Warn(err.Error())
			}
		}
	} else {
		for _, result := range resultMap {
			err := evaluateTarget(result)
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

func evaluateTarget(target oscal.EvalResult) error {

	if len(target.Results) == 0 {
		return fmt.Errorf("no results found")
	}

	if len(target.Results) == 1 {
		// Only one result identified - update to make it the threshold
		oscal.UpdateProps("threshold", "https://docs.lula.dev/ns", "true", target.Results[0].Props)
		message.Warn("less than 2 results found - no comparison possible")
		return nil
	}

	if target.Threshold != nil && target.Latest != nil {
		// Compare the assessment results
		spinner := message.NewProgressSpinner("Evaluating Assessment Results %s against %s", target.Threshold.UUID, target.Latest.UUID)
		defer spinner.Stop()

		message.Debugf("threshold UUID: %s / latest UUID: %s", target.Threshold.UUID, target.Latest.UUID)

		status, findings, err := oscal.EvaluateResults(target.Threshold, target.Latest)
		if err != nil {
			message.Fatal(err, err.Error())
		}

		if status {
			if len(findings["new-passing-findings"]) > 0 {
				message.Info("New passing finding Target-Ids:")
				for _, finding := range findings["new-passing-findings"] {
					message.Infof("%s", finding.Target.TargetId)
				}

				message.Infof("New threshold identified - threshold will be updated to result %s", target.Latest.UUID)

				// Update latest threshold prop
				oscal.UpdateProps("threshold", "https://docs.lula.dev/ns", "true", target.Latest.Props)
				oscal.UpdateProps("threshold", "https://docs.lula.dev/ns", "false", target.Threshold.Props)
			} else {
				// retain result as threshold
				oscal.UpdateProps("threshold", "https://docs.lula.dev/ns", "true", target.Threshold.Props)
			}

			if len(findings["new-failing-findings"]) > 0 {
				message.Info("New failing finding Target-Ids:")
				for _, finding := range findings["new-failing-findings"] {
					message.Infof("%s", finding.Target.TargetId)
				}
			}
			message.Info("Evaluation Passed Successfully")

		} else {
			message.Warn("Evaluation Failed against the following findings:")
			for _, finding := range findings["no-longer-satisfied"] {
				message.Warnf("%s", finding.Target.TargetId)
			}
			message.Fatalf(fmt.Errorf("failed to meet established threshold"), "failed to meet established threshold")

			// retain result as threshold
			oscal.UpdateProps("threshold", "https://docs.lula.dev/ns", "true", target.Threshold.Props)
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
