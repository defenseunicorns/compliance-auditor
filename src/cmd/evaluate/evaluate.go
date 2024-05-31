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
	var threshold, latest *oscalTypes_1_1_2.Result

	// Items for updating the threshold automatically
	var thresholdFile string
	var thresholdAssessment *oscalTypes_1_1_2.AssessmentResults

	// Read in files - establish the results to
	if len(fileArray) == 0 {
		// TODO: Determine if we will handle a default location/name for assessment files
		return fmt.Errorf("no files provided for evaluation")
	}

	for _, f := range fileArray {
		err := files.IsJsonOrYaml(f)
		if err != nil {
			return fmt.Errorf("invalid file extension: %s, requires .json or .yaml", f)
		}
	}

	if len(fileArray) == 1 {
		thresholdFile = fileArray[0]
		data, err := common.ReadFileToBytes(fileArray[0])
		if err != nil {
			return err
		}
		assessment, err := oscal.NewAssessmentResults(data)
		if err != nil {
			return err
		}
		if len(assessment.Results) < 2 {
			message.Infof("%v result object identified - unable to evaluate", len(assessment.Results))
			return nil
		}

		// Identify the threshold
		threshold, err = findThreshold(&assessment.Results)
		if err != nil {
			return err
		}

		latest = &assessment.Results[0]

		status, findings, err = EvaluateResults(threshold, latest)
		if err != nil {
			return err
		}

	} else if len(fileArray) == 2 {
		thresholdFile = fileArray[1]
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
		return fmt.Errorf("exceeded maximum of 2 files for evaluation")
	}

	if status {
		if len(findings["new-passing-findings"]) > 0 {
			message.Info("New passing finding Target-Ids:")
			for _, finding := range findings["new-passing-findings"] {
				message.Infof("%s", finding.Target.TargetId)
			}
			// TODO: If there are new passing Findings -> update the threshold in the assessment
			updateProp("threshold", "false", threshold)
			updateProp("threshold", "true", latest)

			// Props are updated - now write the thresholdAssessment to the existing assessment?
			// if we create the model and write it - the merge will need to de-duplicate instead of merge results
			model := oscalTypes_1_1_2.OscalCompleteSchema{
				AssessmentResults: thresholdAssessment,
			}

			oscal.WriteOscalModel(thresholdFile, &model)

		}

		if len(findings["new-failing-findings"]) > 0 {
			message.Info("New failing finding Target-Ids:")
			for _, finding := range findings["new-failing-findings"] {
				message.Infof("%s", finding.Target.TargetId)
			}
		}

		return nil
	} else {
		message.Warn("Evaluation Failed against the following findings:")
		for _, finding := range findings["no-longer-satisfied"] {
			message.Warnf("%s", finding.Target.TargetId)
		}
		return fmt.Errorf("failed to meet established threshold")
	}
}

func EvaluateResults(thresholdResult *oscalTypes_1_1_2.Result, newResult *oscalTypes_1_1_2.Result) (bool, map[string][]oscalTypes_1_1_2.Finding, error) {
	if thresholdResult.Findings == nil || newResult.Findings == nil {
		return false, nil, fmt.Errorf("results must contain findings to evaluate")
	}

	spinner := message.NewProgressSpinner("Evaluating Assessment Results %s against %s", newResult.UUID, thresholdResult.UUID)
	defer spinner.Stop()

	// Store unique findings for review here
	findings := make(map[string][]oscalTypes_1_1_2.Finding, 0)
	result := true

	findingMapThreshold := oscal.GenerateFindingsMap(*thresholdResult.Findings)
	message.Debug(findingMapThreshold)
	findingMapNew := oscal.GenerateFindingsMap(*newResult.Findings)
	message.Debug(findingMapNew)

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

	message.Debug(findingMapNew)

	// All remaining findings in the new map are new findings
	for _, finding := range findingMapNew {
		if finding.Target.Status.State == "satisfied" {
			message.Debugf("New finding to append: %s", finding.Target.TargetId)
			findings["new-passing-findings"] = append(findings["new-passing-findings"], finding)
		} else {
			findings["new-failing-findings"] = append(findings["new-failing-findings"], finding)
		}

	}

	spinner.Success()
	return result, findings, nil
}

func findThreshold(results *[]oscalTypes_1_1_2.Result) (*oscalTypes_1_1_2.Result, error) {
	for _, result := range *results {
		for _, prop := range *result.Props {
			if prop.Name == "threshold" {
				if prop.Value == "true" {
					return &result, nil
				}
			}
		}
	}
	return &oscalTypes_1_1_2.Result{}, fmt.Errorf("threshold not found")
}

func updateProp(name string, value string, result *oscalTypes_1_1_2.Result) error {
	for index, prop := range *result.Props {
		if prop.Name == name {
			prop.Value = value
			(*result.Props)[index] = prop
			message.Debug(*result)
			return nil
		}
	}
	return fmt.Errorf("property not found")
}
