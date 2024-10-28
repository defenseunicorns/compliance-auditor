package assessmentresults_test

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	assessmentresults "github.com/defenseunicorns/lula/src/internal/tui/assessment_results"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/require"
)

const (
	timeout    = time.Second * 20
	maxRetries = 3
	height     = common.DefaultHeight
	width      = common.DefaultWidth

	validAssessmentResults               = "../../../test/unit/common/oscal/valid-assessment-results.yaml"
	validAssessmentResultsMulti          = "../../../test/unit/common/oscal/valid-assessment-results-multi.yaml"
	validAssessmentResultsRemovedFinding = "../../../test/unit/common/oscal/valid-assessment-results-removed-finding.yaml"
	validAssessmentResultsAddedFinding   = "../../../test/unit/common/oscal/valid-assessment-results-added-finding.yaml"
	validAssessmentResultsRemovedObs     = "../../../test/unit/common/oscal/valid-assessment-results-removed-observation.yaml"
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

// TestAssessmentResultsBasicView tests that the model is created correctly from an assessment results model
func TestAssessmentResultsBasicView(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validAssessmentResults)
	model := assessmentresults.NewAssessmentResultsModel(oscalModel.AssessmentResults)
	model.Open(height, width)

	msgs := []tea.Msg{}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}

// TestAssessmentResultsWithResultsSwitch tests that the model can switch between results
func TestAssessmentResultsWithResultsSwitch(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validAssessmentResultsMulti)
	model := assessmentresults.NewAssessmentResultsModel(oscalModel.AssessmentResults)
	model.Open(height, width)
	mdl, _ := model.Update(tea.WindowSizeMsg{Width: width, Height: height})
	model = mdl.(assessmentresults.Model)

	msgs := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRight}, // Select result
		tea.KeyMsg{Type: tea.KeyEnter}, // Open result detail
		tea.KeyMsg{Type: tea.KeyDown},  // Navigate to next result
		tea.KeyMsg{Type: tea.KeyEnter}, // Select result
		tea.KeyMsg{Type: tea.KeyRight}, // Navigate to compared result
		tea.KeyMsg{Type: tea.KeyRight}, // Navigate to findings table
		tea.KeyMsg{Type: tea.KeyEnter}, // Select finding to filter results
	}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}

// TestAssessmentResultsWithFindingsDetail tests that the model can open the findings detail view
func TestAssessmentResultsWithFindingsDetail(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validAssessmentResults)
	model := assessmentresults.NewAssessmentResultsModel(oscalModel.AssessmentResults)
	model.Open(height, width)

	msgs := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRight},                     // Select result
		tea.KeyMsg{Type: tea.KeyRight},                     // Navigate to compared result
		tea.KeyMsg{Type: tea.KeyRight},                     // Navigate to findings table
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}, // Detail finding
	}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}

// TestAssessmentResultsWithComparison tests that the model can show the comparison between two results
func TestAssessmentResultsWithComparison(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validAssessmentResultsMulti)
	model := assessmentresults.NewAssessmentResultsModel(oscalModel.AssessmentResults)
	model.Open(height, width)

	msgs := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRight}, // Select result
		tea.KeyMsg{Type: tea.KeyRight}, // Navigate to compared result
		tea.KeyMsg{Type: tea.KeyEnter}, // Open compared result detail
		tea.KeyMsg{Type: tea.KeyDown},  // Navigate to next result
		tea.KeyMsg{Type: tea.KeyEnter}, // Select result
		tea.KeyMsg{Type: tea.KeyRight}, // Navigate to findings table
	}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}

// TestAssessmentResultsWithComparison tests that the model can show the comparison between two results
func TestAssessmentResultsWithObservationComparisonDetail(t *testing.T) {
	oscalModel := testhelpers.OscalFromPath(t, validAssessmentResultsMulti)
	model := assessmentresults.NewAssessmentResultsModel(oscalModel.AssessmentResults)
	model.Open(height, width)

	msgs := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRight},                     // Select result
		tea.KeyMsg{Type: tea.KeyRight},                     // Navigate to compared result
		tea.KeyMsg{Type: tea.KeyEnter},                     // Open compared result detail
		tea.KeyMsg{Type: tea.KeyDown},                      // Navigate to next result
		tea.KeyMsg{Type: tea.KeyEnter},                     // Select result
		tea.KeyMsg{Type: tea.KeyRight},                     // Navigate to findings table
		tea.KeyMsg{Type: tea.KeyRight},                     // Navigate to observations table
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}, // Detail observation
	}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}
