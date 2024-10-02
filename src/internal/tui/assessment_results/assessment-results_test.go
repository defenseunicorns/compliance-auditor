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
)

const (
	timeout    = time.Second * 20
	maxRetries = 3
	height     = common.DefaultHeight
	width      = common.DefaultWidth

	validAssessmentResults = "../../../test/unit/common/oscal/valid-assessment-results.yaml"
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

	err := testhelpers.RunTestModelView(t, model, msgs, timeout, maxRetries, height, width)
	if err != nil {
		t.Fatal(err)
	}
}
