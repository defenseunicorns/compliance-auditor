package tui_test

import (
	"os"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/testhelpers"
	"github.com/defenseunicorns/lula/src/internal/tui"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	timeout    = time.Second * 20
	maxRetries = 3
	height     = common.DefaultHeight
	width      = common.DefaultWidth
)

func init() {
	lipgloss.SetColorProfile(termenv.Ascii)
}

// TestNewOSCALModel tests that the NewOSCALModel creates the expected model from the provided map of OSCAL models
func TestNewOSCALModel(t *testing.T) {
	tempLog := testhelpers.CreateTempFile(t, "log")
	defer os.Remove(tempLog.Name())

	oscalComponent := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-component.yaml")
	oscalAssessment := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-assessment-results.yaml")
	model := tui.NewOSCALModel(
		map[string]*oscalTypes_1_1_2.OscalCompleteSchema{
			"component":          oscalComponent,
			"assessment-results": oscalAssessment,
		},
		map[string]string{
			"component":          "component.yaml",
			"assessment-results": "assessment-results.yaml",
		}, tempLog)

	msgs := []tea.Msg{}

	err := testhelpers.RunTestModelView(t, model, nil, msgs, timeout, maxRetries, height, width)
	require.NoError(t, err)
}

func TestNewOSCALModelWithSave(t *testing.T) {
	tempLog := testhelpers.CreateTempFile(t, "log")
	defer os.Remove(tempLog.Name())

	tempOscalFile := testhelpers.CreateTempFile(t, "yaml")
	defer os.Remove(tempOscalFile.Name())

	oscalComponent := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-generated-component.yaml")
	model := tui.NewOSCALModel(
		map[string]*oscalTypes_1_1_2.OscalCompleteSchema{
			"component": oscalComponent,
		},
		map[string]string{
			"component": tempOscalFile.Name(),
		}, tempLog)

	msgs := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRight},                                    // Select component
		tea.KeyMsg{Type: tea.KeyRight},                                    // Select framework
		tea.KeyMsg{Type: tea.KeyRight},                                    // Select control
		tea.KeyMsg{Type: tea.KeyEnter},                                    // Open control
		tea.KeyMsg{Type: tea.KeyRight},                                    // Navigate to remarks
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}},                // Edit remarks
		tea.KeyMsg{Type: tea.KeyCtrlE},                                    // Newline
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'e', 's', 't'}}, // Add "test" to remarks
		tea.KeyMsg{Type: tea.KeyEnter},                                    // Confirm edit
		tea.KeyMsg{Type: tea.KeyCtrlS},                                    // Save
		tea.KeyMsg{Type: tea.KeyEnter},                                    // Confirm Save
	}

	// Execute model to read output file
	tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(width, height))

	for _, msg := range msgs {
		tm.Send(msg)
		time.Sleep(time.Millisecond * time.Duration(50))
	}

	time.Sleep(time.Second * 2) // Give it a couple seconds to write the file

	err := tm.Quit()
	assert.NoError(t, err)

	// Check output file contents
	modifiedOscalModel := testhelpers.OscalFromPath(t, tempOscalFile.Name())
	compDefn := modifiedOscalModel.ComponentDefinition
	require.NotNil(t, compDefn)

	for _, c := range *compDefn.Components {
		require.NotNil(t, c.ControlImplementations)
		for _, f := range *c.ControlImplementations {
			for _, r := range f.ImplementedRequirements {
				if r.ControlId == "ac-1" {
					if !strings.HasSuffix(r.Remarks, "test") {
						t.Errorf("Expected remarks in ac-1 to contain 'test' at the end, got %s", r.Remarks)
					}
				}
			}
		}
	}
}

func TestNewOSCALWithValidate(t *testing.T) {
	tempLog := testhelpers.CreateTempFile(t, "log")
	defer os.Remove(tempLog.Name())
	message.UseLogFile(tempLog)
	message.NoProgress = true

	tempOscalFile := testhelpers.CreateTempFile(t, "yaml")
	defer os.Remove(tempOscalFile.Name())

	oscalComponent := testhelpers.OscalFromPath(t, "../../test/unit/common/oscal/valid-multi-component-validations.yaml")
	model := tui.NewOSCALModel(
		map[string]*oscalTypes_1_1_2.OscalCompleteSchema{
			"component": oscalComponent,
		},
		map[string]string{
			"component":          "none.yaml",
			"assessment-results": tempOscalFile.Name(),
		}, tempLog)

	msgs := []tea.Msg{
		tea.KeyMsg{Type: tea.KeyRight}, // Select component
		tea.KeyMsg{Type: tea.KeyRight}, // Select framework
		tea.KeyMsg{Type: tea.KeyEnter}, // Open framework
		tea.KeyMsg{Type: tea.KeyDown},  // Navigate to rev4
		tea.KeyMsg{Type: tea.KeyDown},  // Navigate to rev4
		tea.KeyMsg{Type: tea.KeyEnter}, // Select rev4
		tea.KeyMsg{Type: tea.KeyCtrlV}, // Open validation
		tea.KeyMsg{Type: tea.KeyEnter}, // Run validation
	}

	tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(width, height))

	for _, msg := range msgs {
		tm.Send(msg)
		time.Sleep(time.Millisecond * time.Duration(50))
	}

	time.Sleep(time.Second * 5) // Give it a few seconds to execute and write the file

	err := tm.Quit()
	assert.NoError(t, err)

	// Check output file contents
	oscalModel := testhelpers.OscalFromPath(t, tempOscalFile.Name())
	assessmentResults := oscalModel.AssessmentResults
	require.NotNil(t, assessmentResults)

	for _, result := range assessmentResults.Results {
		require.NotNil(t, result.Findings)
		assert.Equal(t, 6, len(*result.Findings))

		require.NotNil(t, result.Observations)
		assert.Equal(t, 2, len(*result.Observations))
	}
}
