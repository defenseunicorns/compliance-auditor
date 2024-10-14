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
	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
)

const (
	validateWidthScale  = 0.5
	validateHeightScale = 0.6
	defaultPopupHeight  = 15
	defaultPopupWidth   = 50
)

type ValidateStartMsg struct{}
type ValidationCompleteMsg struct {
	Err error
}
type ValidationDataMsg struct {
	AssessmentResults *oscalTypes_1_1_2.AssessmentResults
}
type ValidateCloseAndResetMsg struct{}
type ValidateModelMsg struct {
	RunExecutables bool
}

type ValidateModel struct {
	IsOpen            bool
	runExecutable     bool
	validating        bool
	validatable       bool
	content           string
	filePath          string
	componentModel    *Model
	assessmentResults *oscalTypes_1_1_2.AssessmentResults
	help              common.HelpModel
	height            int
	width             int
}

func NewValidateModel(componentModel *Model, filepath string) ValidateModel {
	help := common.NewHelpModel(true)
	help.ShortHelp = []key.Binding{common.CommonKeys.Confirm, common.CommonKeys.Cancel}

	return ValidateModel{
		help:           help,
		componentModel: componentModel,
		filePath:       filepath,
		height:         defaultPopupHeight,
		width:          defaultPopupWidth,
	}
}

func (m ValidateModel) Init() tea.Cmd {
	return nil
}

func (m ValidateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var err error

	if m.IsOpen {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.updateSizing(int(float64(msg.Height)*validateHeightScale), int(float64(msg.Width)*validateWidthScale))

		case tea.KeyMsg:
			k := msg.String()

			switch k {
			case common.ContainsKey(k, common.CommonKeys.Confirm.Keys()):
				if m.validatable {
					m.validating = true
					cmds = append(cmds, func() tea.Msg {
						return ValidateStartMsg{}
					})
				} else {
					m.IsOpen = false
				}

			case common.ContainsKey(k, common.CommonKeys.Cancel.Keys()):
				m.IsOpen = false
			}

		case ValidateStartMsg:
			m.assessmentResults, err = m.RunValidations(m.runExecutable)

			cmds = append(cmds, func() tea.Msg {
				return ValidationCompleteMsg{
					Err: err,
				}
			})

		case ValidationCompleteMsg:
			m.validating = false

			cmds = append(cmds, func() tea.Msg {
				time.Sleep(time.Second * 2)
				return ValidateCloseAndResetMsg{}
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

		case ValidateCloseAndResetMsg:
			m.IsOpen = false
			m.validatable = false
			m.validating = false

		}
	}
	return m, tea.Sequence(cmds...)
}

func (m ValidateModel) View() string {
	common.PrintToLog("in validate view")
	common.DumpToLog(m)
	var content string
	popupStyle := common.OverlayStyle.
		Width(m.width).
		Height(m.height)

	if m.validating {
		// TODO: Add progress spinner/feedback
		content = "Validating..."
		m.help.ShortHelp = []key.Binding{}
	} else {
		content = m.content
	}
	validationContent := lipgloss.JoinVertical(lipgloss.Top, content, "\n", m.help.View())
	return popupStyle.Render(validationContent)
}

func (m *ValidateModel) Open(height, width int, oscalComponent *oscalTypes_1_1_2.ComponentDefinition, target string) {
	m.IsOpen = true
	m.updateSizing(int(float64(height)*validateHeightScale), int(float64(width)*validateWidthScale))
	var content strings.Builder

	// update the model with component data
	var validationStore *validationstore.ValidationStore
	if oscalComponent != nil {
		if oscalComponent.BackMatter != nil {
			validationStore = validationstore.NewValidationStoreFromBackMatter(*oscalComponent.BackMatter)
		} else {
			validationStore = validationstore.NewValidationStore()
		}

		hasExecutables, msg := validationStore.DryRun()

		content.WriteString(fmt.Sprintf("Validate Component Definition on Target: %s\n\n%s", target, msg))
		if hasExecutables {
			content.WriteString("\n⚠️ Includes Executable Validations ⚠️\n")
		}
		m.validatable = true
	} else {
		content.WriteString("Nothing to Validate")
	}
	m.content = content.String()
}

func (m *ValidateModel) updateSizing(height, width int) {
	if height < defaultPopupHeight {
		height = defaultPopupHeight
	}
	if width < defaultPopupWidth {
		width = defaultPopupWidth
	}
}

func (m *ValidateModel) RunValidations(runExecutable bool) (*oscalTypes_1_1_2.AssessmentResults, error) {
	validator, err := validation.New(
		validation.WithAllowExecution(runExecutable, true),
	)
	if err != nil {
		return nil, err
	}
	framework := m.componentModel.GetSelectedFramework()
	if framework.OscalFramework == nil {
		return nil, fmt.Errorf("framework is nil")
	}

	results, err := validator.ValidateOnCompDef(context.Background(), m.componentModel.GetComponentDefinition(), framework.Name)
	if err != nil {
		return nil, err
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
