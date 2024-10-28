package component

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	requirementstore "github.com/defenseunicorns/lula/src/pkg/common/requirement-store"
	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
)

const (
	validateWidthScale  = 0.5
	validateHeightScale = 0.3
	minimumHeight       = 15
	minimumWidth        = 50
)

type ValidateOpenMsg struct {
	Height int
	Width  int
	Target string
}
type ValidateStartMsg struct{}
type ValidationCompleteMsg struct {
	Err error
}
type ValidationDataMsg struct {
	AssessmentResults *oscalTypes_1_1_2.AssessmentResults
}
type ValidateCloseMsg struct{}

type ValidateModel struct {
	IsOpen            bool
	runExecutable     bool
	validating        bool
	validatable       bool
	target            string
	content           string
	oscalComponent    *oscalTypes_1_1_2.ComponentDefinition
	controlImplSet    []oscalTypes_1_1_2.ControlImplementationSet
	validationStore   *validationstore.ValidationStore
	assessmentResults *oscalTypes_1_1_2.AssessmentResults
	help              common.HelpModel
	height            int
	width             int
}

func NewValidateModel(oscalComponent *oscalTypes_1_1_2.ComponentDefinition) ValidateModel {
	help := common.NewHelpModel(true)
	help.ShortHelp = []key.Binding{common.CommonKeys.Confirm, common.CommonKeys.Cancel}

	return ValidateModel{
		help:           help,
		oscalComponent: oscalComponent,
		runExecutable:  true, // Hardcoding for now, tbd on optional input from user
		height:         minimumHeight,
		width:          minimumWidth,
	}
}

func (m ValidateModel) Init() tea.Cmd {
	return nil
}

func (m ValidateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var err error

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateSizing(int(float64(msg.Height)*validateHeightScale), int(float64(msg.Width)*validateWidthScale))

	case tea.KeyMsg:
		k := msg.String()

		switch k {
		case common.ContainsKey(k, common.CommonKeys.Confirm.Keys()):
			if m.validatable {
				m.validating = true
				m.content = "Validating..."
				cmds = append(cmds, func() tea.Msg {
					return ValidateStartMsg{}
				})
			} else {
				m.IsOpen = false
			}

		case common.ContainsKey(k, common.CommonKeys.Cancel.Keys()):
			m.IsOpen = false
		}

	case ValidateOpenMsg:
		m.IsOpen = true
		m.target = msg.Target
		m.updateSizing(int(float64(msg.Height)*validateHeightScale), int(float64(msg.Width)*validateWidthScale))
		m.setInitialContent()

	case ValidateStartMsg:
		validationStart := time.Now()
		m.assessmentResults, err = m.RunValidations(m.runExecutable, m.target)
		if err != nil {
			common.PrintToLog("error running validations: %v", err)
		}
		validationDuration := time.Since(validationStart)
		// just adding a minimum of 2 seconds to the "validating" popup
		if validationDuration < time.Second*2 {
			time.Sleep(time.Second*2 - validationDuration)
		}

		cmds = append(cmds, func() tea.Msg {
			return ValidationCompleteMsg{
				Err: err,
			}
		})

	case ValidationCompleteMsg:
		cmds = append(cmds, func() tea.Msg {
			time.Sleep(time.Second * 2)
			return ValidateCloseMsg{}
		})
		if msg.Err != nil {
			m.content = fmt.Sprintf("Error running validation: %v", msg.Err)
		} else {
			m.content = "Validation Complete"
			cmds = append(cmds, func() tea.Msg {
				return ValidationDataMsg{
					AssessmentResults: m.assessmentResults,
				}
			})
		}
		return m, tea.Sequence(cmds...)

	case ValidateCloseMsg:
		m.IsOpen = false
		m.validatable = false
		m.validating = false
		m.target = ""

	}

	return m, tea.Batch(cmds...)
}

func (m ValidateModel) View() string {
	popupStyle := common.OverlayStyle.
		Width(m.width).
		Height(m.height)

	if m.validating {
		// Add progress spinner/feedback?
		m.help.ShortHelp = []key.Binding{}
	}
	m.fillContentPane()
	validationContent := lipgloss.JoinVertical(lipgloss.Top, m.content, "\n", m.help.View())
	return popupStyle.Render(validationContent)
}

func (m *ValidateModel) setInitialContent() {
	var content strings.Builder

	// update the model with component data
	if m.oscalComponent != nil {
		if m.oscalComponent.BackMatter != nil {
			m.validationStore = validationstore.NewValidationStoreFromBackMatter(*m.oscalComponent.BackMatter)
		} else {
			m.validationStore = validationstore.NewValidationStore()
		}

		controlImplementationMap := oscal.FilterControlImplementations(m.oscalComponent)

		if controlImplSet, ok := controlImplementationMap[m.target]; ok {
			m.controlImplSet = controlImplSet
			requirementStore := requirementstore.NewRequirementStore(&controlImplSet)

			content.WriteString("Run Validations?\n\n")
			content.WriteString("🔍 Validate Component Definition on Target: ")
			content.WriteString(m.target)
			requirementStore.ResolveLulaValidations(m.validationStore)
			reqtStats := requirementStore.GetStats(m.validationStore)
			content.WriteString(fmt.Sprintf("\n\n• Found %d Implemented Requirements", reqtStats.TotalRequirements))
			content.WriteString(fmt.Sprintf("\n• Found %d runnable Lula Validations", reqtStats.TotalValidations))

			hasExecutables, _ := m.validationStore.DryRun()
			if hasExecutables {
				content.WriteString("\n⚠️ Includes Executable Validations ⚠️\n")
			}
			m.validatable = true
		} else {
			content.WriteString("No Framework selected")
		}
	} else {
		content.WriteString("Nothing to Validate")
	}
	m.content = content.String()
}

func (m *ValidateModel) updateSizing(height, width int) {
	m.height = common.Max(height, minimumHeight)
	m.width = common.Max(width, minimumWidth)
}

func (m *ValidateModel) RunValidations(runExecutable bool, target string) (*oscalTypes_1_1_2.AssessmentResults, error) {
	validator, err := validation.New(
		validation.WithAllowExecution(runExecutable, true),
	)
	if err != nil {
		return nil, err
	}

	results := make([]oscalTypes_1_1_2.Result, 0)
	if len(m.controlImplSet) > 0 {
		findings, observations, _ := validator.ValidateOnControlImplementations(context.Background(), &m.controlImplSet, m.validationStore, target)
		result, err := oscal.CreateResult(findings, observations)
		if err != nil {
			return nil, err
		}
		// add/update the source to the result props - make source = framework or omit?
		oscal.UpdateProps("target", oscal.LULA_NAMESPACE, target, result.Props)
		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results produced")
	}

	assessmentresults, err := oscal.GenerateAssessmentResults(results)
	if err != nil {
		return nil, err
	}

	return assessmentresults, nil
}

// Helper method to fill up the validate model content pane with newlines so help is pushed down
func (m *ValidateModel) fillContentPane() {
	availableContentHeight := m.height - common.HelpStyle(m.width).GetHeight() - 2 // 2 for the border
	contentHeight := strings.Count(m.content, "\n")
	for range availableContentHeight - contentHeight {
		m.content += "\n"
	}
}
