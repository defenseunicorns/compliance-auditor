package evaluate

import (
	"fmt"
	"slices"

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

	if len(fileArray) == 0 {
		return fmt.Errorf("no files provided for evaluation")
	}

	// Potentially write changes back to multiple files requires some storage
	resultMap := make(map[string]*oscalTypes_1_1_2.AssessmentResults)
	for _, fileString := range fileArray {
		err := files.IsJsonOrYaml(fileString)
		if err != nil {
			return fmt.Errorf("invalid file extension: %s, requires .json or .yaml", fileString)
		}

		data, err := common.ReadFileToBytes(fileString)
		if err != nil {
			return err
		}
		assessment, err := oscal.NewAssessmentResults(data)
		if err != nil {
			return err
		}
		resultMap[fileString] = assessment
	}

	// Now that we have the map of assessment results - we need to identify the threshold(s)
	// Also sort the results -> if we maintain pointers, can we update and write all artifacts in one go?

	thresholds, sortedResults, err := findAndSortResults(resultMap)
	if err != nil {
		return err
	}

	if len(sortedResults) <= 1 {
		// Should this implicitly pass? If so then a workflow can operate on the assumption that it will pass from 0 -> N results
		message.Infof("%v result object identified - unable to evaluate", len(sortedResults))
		return nil
	}

	if len(thresholds) == 0 {
		// No thresholds identified but we have > 1 results - compare the latest and the preceding
		threshold = sortedResults[len(sortedResults)-2]
		latest = sortedResults[len(sortedResults)-1]

		status, findings, err = EvaluateResults(threshold, latest)
		if err != nil {
			return err
		}
	} else {
		// Constraint - Always evaluate the latest threshold against the latest result
		threshold = thresholds[len(thresholds)-1]
		latest = sortedResults[len(sortedResults)-1]

		if threshold.UUID == latest.UUID {
			// They are the same - return error
			return fmt.Errorf("unable to evaluate - threshold and latest result are the same result - nothing to compare")
		}
		status, findings, err = EvaluateResults(threshold, latest)
		if err != nil {
			return err
		}
	}

	if status {
		if len(findings["new-passing-findings"]) > 0 {
			message.Info("New passing finding Target-Ids:")
			for _, finding := range findings["new-passing-findings"] {
				message.Infof("%s", finding.Target.TargetId)
			}

			message.Info("New threshold identified - threshold will be updated to latest result")

			updateProp("threshold", "false", threshold.Props)
			updateProp("threshold", "true", latest.Props)

			// Props are updated - now write back to all files
			// if we create the model and write it - the merge will need to de-duplicate instead of merge results
			for filePath, assessment := range resultMap {
				model := oscalTypes_1_1_2.OscalCompleteSchema{
					AssessmentResults: assessment,
				}

				oscal.WriteOscalModel(filePath, &model)
			}

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
		if finding.Target.Status.State == "satisfied" {
			findings["new-passing-findings"] = append(findings["new-passing-findings"], finding)
		} else {
			findings["new-failing-findings"] = append(findings["new-failing-findings"], finding)
		}

	}

	spinner.Success()
	return result, findings, nil
}

// findAndSortResults takes a map of results and returns a list of thresholds and a sorted list of results in order of time
func findAndSortResults(resultMap map[string]*oscalTypes_1_1_2.AssessmentResults) ([]*oscalTypes_1_1_2.Result, []*oscalTypes_1_1_2.Result, error) {

	thresholds := make([]*oscalTypes_1_1_2.Result, 0)
	sortedResults := make([]*oscalTypes_1_1_2.Result, 0)

	for _, assessment := range resultMap {
		for _, result := range assessment.Results {
			if result.Props != nil {
				for _, prop := range *result.Props {
					if prop.Name == "threshold" && prop.Value == "true" {
						thresholds = append(thresholds, &result)
					}
				}
			}
			// Store all results in a non-sorted list
			sortedResults = append(sortedResults, &result)
		}
	}

	// Sort the results by start time
	slices.SortFunc(sortedResults, func(a, b *oscalTypes_1_1_2.Result) int { return a.Start.Compare(b.Start) })
	slices.SortFunc(thresholds, func(a, b *oscalTypes_1_1_2.Result) int { return a.Start.Compare(b.Start) })

	return thresholds, sortedResults, nil
}

func updateProp(name string, value string, props *[]oscalTypes_1_1_2.Property) error {

	for index, prop := range *props {
		if prop.Name == name {
			prop.Value = value
			(*props)[index] = prop
			return nil
		}
	}
	return fmt.Errorf("property not found")
}
