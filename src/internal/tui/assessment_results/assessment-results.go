package assessmentresults

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
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
	uuid, title      string
	timestamp        string
	oscalResult      *oscalTypes_1_1_2.Result
	findings         *[]oscalTypes_1_1_2.Finding
	observations     *[]oscalTypes_1_1_2.Observation
	findingsRows     []table.Row
	observationsRows []table.Row
	observationsMap  map[string]table.Row
	summaryData      summaryData
}

type summaryData struct {
	numFindings, numObservations int
	numFindingsSatisfied         int
	numObservationsSatisfied     int
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
					}
				}
				if f.Target.Status.State == "satisfied" {
					numFindingsSatisfied++
				}

				style, exists := satisfiedColors[f.Target.Status.State]
				if !exists {
					style = satisfiedColors["other"]
				}

				findingsRows = append(findingsRows, table.NewRow(table.RowData{
					columnKeyName:        f.Target.TargetId,
					columnKeyStatus:      table.NewStyledCell(f.Target.Status.State, style),
					columnKeyDescription: strings.ReplaceAll(f.Description, "\n", " "),
					// Hidden columns
					columnKeyFinding:    findingString,
					columnKeyRelatedObs: relatedObs,
				}))
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
				obsRow := table.NewRow(table.RowData{
					columnKeyName:        GetReadableObservationName(o.Description),
					columnKeyStatus:      table.NewStyledCell(state, style),
					columnKeyDescription: remarks.String(),
					// Hidden columns
					columnKeyObservation:  obsString,
					columnKeyValidationId: findUuid(o.Description),
				})
				observationsRows = append(observationsRows, obsRow)
				observationsMap[o.UUID] = obsRow
			}

			results = append(results, result{
				uuid:             r.UUID,
				title:            r.Title,
				oscalResult:      &r,
				timestamp:        r.Start.Format(time.RFC3339),
				findings:         r.Findings,
				observations:     r.Observations,
				findingsRows:     findingsRows,
				observationsRows: observationsRows,
				observationsMap:  observationsMap,
				summaryData: summaryData{
					numFindings:              numFindings,
					numObservations:          numObservations,
					numFindingsSatisfied:     numFindingsSatisfied,
					numObservationsSatisfied: numObservationsSatisfied,
				},
			})
		}
	}

	return results
}

func getComparedResults(results []result, selectedResult result) []string {
	comparedResults := []string{"None"}
	for _, r := range results {
		if r.uuid != selectedResult.uuid {
			comparedResults = append(comparedResults, getResultText(r))
		}
	}
	return comparedResults
}

func getResultText(result result) string {
	var resultText strings.Builder
	if result.uuid == "" {
		return "No Result Selected"
	}
	resultText.WriteString(result.title)
	if result.oscalResult != nil {
		thresholdFound, threshold := oscal.GetProp("threshold", oscal.LULA_NAMESPACE, result.oscalResult.Props)
		if thresholdFound && threshold == "true" {
			resultText.WriteString(", Threshold")
		}
		targetFound, target := oscal.GetProp("target", oscal.LULA_NAMESPACE, result.oscalResult.Props)
		if targetFound {
			resultText.WriteString(fmt.Sprintf(", %s", target))
		}
	}
	resultText.WriteString(fmt.Sprintf(" - %s", result.timestamp))

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
