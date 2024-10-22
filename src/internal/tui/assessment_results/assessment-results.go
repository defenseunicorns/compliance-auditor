package assessmentresults

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	pkgResult "github.com/defenseunicorns/lula/src/pkg/common/result"
	"github.com/evertras/bubble-table/table"
)

var (
	satisfiedColors = map[string]lipgloss.Style{
		"satisfied":     lipgloss.NewStyle().Foreground(lipgloss.Color("#3ad33c")),
		"not-satisfied": lipgloss.NewStyle().Foreground(lipgloss.Color("#e36750")),
		"other":         lipgloss.NewStyle().Foreground(lipgloss.Color("#f3f3f3")),
	}
)

type result struct {
	Uuid, Title      string
	Timestamp        string
	OscalResult      *oscalTypes_1_1_2.Result
	Findings         *[]oscalTypes_1_1_2.Finding
	Observations     *[]oscalTypes_1_1_2.Observation
	FindingsRows     []table.Row
	ObservationsRows []table.Row
	FindingsMap      map[string]table.Row
	ObservationsMap  map[string]table.Row
	SummaryData      summaryData
}

type summaryData struct {
	NumFindings, NumObservations int
	NumFindingsSatisfied         int
	NumObservationsSatisfied     int
}

func GetResults(assessmentResults *oscalTypes_1_1_2.AssessmentResults) []result {
	results := make([]result, 0)

	if assessmentResults != nil {
		for _, r := range assessmentResults.Results {
			numFindings := len(*r.Findings)
			numObservations := len(*r.Observations)
			numFindingsSatisfied := 0
			numObservationsSatisfied := 0
			findingsRows := make([]table.Row, 0)
			observationsRows := make([]table.Row, 0)
			observationsMap := make(map[string]table.Row)
			findingsMap := make(map[string]table.Row)
			observationsControlMap := make(map[string][]string, 0)

			for _, f := range *r.Findings {
				findingString, err := common.ToYamlString(f)
				if err != nil {
					common.PrintToLog("error converting finding to yaml: %v", err)
					findingString = ""
				}
				relatedObs := make([]string, 0)
				if f.RelatedObservations != nil {
					for _, o := range *f.RelatedObservations {
						relatedObs = append(relatedObs, o.ObservationUuid)
						if _, ok := observationsControlMap[o.ObservationUuid]; !ok {
							observationsControlMap[o.ObservationUuid] = []string{f.Target.TargetId}
						} else {
							observationsControlMap[o.ObservationUuid] = append(observationsControlMap[o.ObservationUuid], f.Target.TargetId)
						}
					}
				}
				if f.Target.Status.State == "satisfied" {
					numFindingsSatisfied++
				}

				style, exists := satisfiedColors[f.Target.Status.State]
				if !exists {
					style = satisfiedColors["other"]
				}

				findingRow := table.NewRow(table.RowData{
					ColumnKeyName:        f.Target.TargetId,
					ColumnKeyStatus:      table.NewStyledCell(f.Target.Status.State, style),
					ColumnKeyDescription: strings.ReplaceAll(f.Description, "\n", " "),
					// Hidden columns
					ColumnKeyFinding:    findingString,
					ColumnKeyRelatedObs: relatedObs,
				})
				findingsRows = append(findingsRows, findingRow)
				findingsMap[f.Target.TargetId] = findingRow
			}

			for _, o := range *r.Observations {
				state := "undefined"
				var remarks strings.Builder
				if o.RelevantEvidence != nil {
					for _, e := range *o.RelevantEvidence {
						if e.Description == "Result: satisfied\n" {
							state = "satisfied"
						} else if e.Description == "Result: not-satisfied\n" {
							state = "not-satisfied"
						}
						if e.Remarks != "" {
							remarks.WriteString(strings.ReplaceAll(e.Remarks, "\n", " "))
						}
					}
					if state == "satisfied" {
						numObservationsSatisfied++
					}
				}

				style, exists := satisfiedColors[state]
				if !exists {
					style = satisfiedColors["other"]
				}

				obsString, err := common.ToYamlString(o)
				if err != nil {
					common.PrintToLog("error converting observation to yaml: %v", err)
					obsString = ""
				}

				var controlIds []string
				if ids, ok := observationsControlMap[o.UUID]; ok {
					controlIds = ids
				}

				obsRow := table.NewRow(table.RowData{
					ColumnKeyName:        GetReadableObservationName(o.Description),
					ColumnKeyStatus:      table.NewStyledCell(state, style),
					ColumnKeyControlIds:  strings.Join(controlIds, ", "),
					ColumnKeyDescription: remarks.String(),
					// Hidden columns
					ColumnKeyObservation:  obsString,
					ColumnKeyValidationId: findUuid(o.Description),
				})
				observationsRows = append(observationsRows, obsRow)
				observationsMap[o.UUID] = obsRow
			}

			results = append(results, result{
				Uuid:             r.UUID,
				Title:            r.Title,
				OscalResult:      &r,
				Timestamp:        r.Start.Format(time.RFC3339),
				Findings:         r.Findings,
				Observations:     r.Observations,
				FindingsRows:     findingsRows,
				ObservationsRows: observationsRows,
				FindingsMap:      findingsMap,
				ObservationsMap:  observationsMap,
				SummaryData: summaryData{
					NumFindings:              numFindings,
					NumObservations:          numObservations,
					NumFindingsSatisfied:     numFindingsSatisfied,
					NumObservationsSatisfied: numObservationsSatisfied,
				},
			})
		}
	}

	return results
}

func GetResultComparison(selectedResult, comparedResult result) ([]table.Row, []table.Row) {
	findingsRows := make([]table.Row, 0)
	observationsRows := make([]table.Row, 0)
	observations := make([]string, 0)

	if selectedResult.OscalResult != nil && comparedResult.OscalResult != nil {
		resultComparison := pkgResult.NewResultComparisonMap(*selectedResult.OscalResult, *comparedResult.OscalResult)
		for k, v := range resultComparison {
			// Make compared finding row
			var comparedFindingRow table.Row
			var ok bool
			if comparedFindingRow, ok = selectedResult.FindingsMap[k]; ok {
				comparedFindingRow.Data[ColumnKeyStatusChange] = v.StateChange
				if r, ok := comparedResult.FindingsMap[k]; ok {
					// Finding exists in both results
					comparedFindingRow.Data[ColumnKeyComparedFinding] = r.Data[ColumnKeyFinding]
				} else {
					// Finding is new
					comparedFindingRow.Data[ColumnKeyComparedFinding] = ""
				}
			} else {
				if comparedFindingRow, ok = comparedResult.FindingsMap[k]; ok {
					// Finding was removed
					comparedFindingRow.Data[ColumnKeyComparedFinding] = comparedFindingRow.Data[ColumnKeyFinding]
					comparedFindingRow.Data[ColumnKeyFinding] = ""
					comparedFindingRow.Data[ColumnKeyStatusChange] = v.StateChange
				}
			}
			findingsRows = append(findingsRows, comparedFindingRow)

			// Make compared observation row
			for _, op := range v.ObservationPairs {
				if op != nil {
					obsUuid := ""
					var comparedObservationRow table.Row
					if comparedObservationRow, ok = selectedResult.ObservationsMap[op.ObservationUuid]; ok {
						obsUuid = op.ObservationUuid
						comparedObservationRow.Data[ColumnKeyStatusChange] = op.StateChange
						if r, ok := comparedResult.ObservationsMap[op.ComparedObservationUuid]; ok {
							comparedObservationRow.Data[ColumnKeyComparedObservation] = r.Data[ColumnKeyObservation]
						} else {
							// Observation is new
							comparedObservationRow.Data[ColumnKeyComparedObservation] = ""
						}
					} else {
						if comparedObservationRow, ok = comparedResult.ObservationsMap[op.ComparedObservationUuid]; ok {
							// Observation was removed
							obsUuid = op.ComparedObservationUuid
							comparedObservationRow.Data[ColumnKeyStatusChange] = op.StateChange
							comparedObservationRow.Data[ColumnKeyComparedObservation] = comparedObservationRow.Data[ColumnKeyObservation]
							comparedObservationRow.Data[ColumnKeyObservation] = ""
						}
					}
					// Check if observation has already been added
					if obsUuid != "" && !slices.Contains(observations, obsUuid) {
						observations = append(observations, obsUuid)
						observationsRows = append(observationsRows, comparedObservationRow)
					}
				}
			}
		}
	}

	return findingsRows, observationsRows
}

func getComparedResults(results []result, selectedResult result) []string {
	comparedResults := []string{"None"}
	for _, r := range results {
		if r.Uuid != selectedResult.Uuid {
			comparedResults = append(comparedResults, getResultText(r))
		}
	}
	return comparedResults
}

func getResultText(result result) string {
	var resultText strings.Builder
	if result.Uuid == "" {
		return "No Result Selected"
	}
	resultText.WriteString(result.Title)
	if result.OscalResult != nil {
		thresholdFound, threshold := oscal.GetProp("threshold", oscal.LULA_NAMESPACE, result.OscalResult.Props)
		if thresholdFound && threshold == "true" {
			resultText.WriteString(", Threshold")
		}
		targetFound, target := oscal.GetProp("target", oscal.LULA_NAMESPACE, result.OscalResult.Props)
		if targetFound {
			resultText.WriteString(fmt.Sprintf(", %s", target))
		}
	}
	resultText.WriteString(fmt.Sprintf(" - %s", result.Timestamp))

	return resultText.String()
}

func findUuid(input string) string {
	uuidPattern := `[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`

	re := regexp.MustCompile(uuidPattern)

	return re.FindString(input)
}

func GetReadableObservationName(desc string) string {
	// Define the regular expression pattern
	pattern := `\[TEST\]: ([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}) - (.+)`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find the matches
	matches := re.FindStringSubmatch(desc)

	if len(matches) == 3 {
		message := matches[2]

		return message
	} else {
		return desc
	}
}
