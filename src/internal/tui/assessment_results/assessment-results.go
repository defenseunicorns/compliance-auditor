package assessmentresults

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

const (
	height           = 20
	width            = 12
	pickerHeight     = 20
	pickerWidth      = 80
	dialogFixedWidth = 40
)

func NewAssessmentResultsModel(assessmentResults *oscalTypes_1_1_2.AssessmentResults) Model {
	results := make([]result, 0)
	findings := make([]blist.Item, 0)
	var selectedResult result

	if assessmentResults != nil {
		for _, r := range assessmentResults.Results {
			results = append(results, result{
				uuid:         r.UUID,
				title:        r.Title,
				findings:     r.Findings,
				observations: r.Observations,
			})
		}
	}
	if len(results) != 0 {
		selectedResult = results[0]
		observationMap := makeObservationMap(selectedResult.observations)
		if selectedResult.findings != nil {
			for _, f := range *selectedResult.findings {
				// get the related observations
				observations := make([]observation, 0)
				if f.RelatedObservations != nil {
					for _, o := range *f.RelatedObservations {
						observationUuid := o.ObservationUuid
						if _, ok := observationMap[observationUuid]; ok {
							observations = append(observations, observationMap[observationUuid])
						}
					}
				}
				findings = append(findings, finding{
					title:        f.Title,
					uuid:         f.UUID,
					controlId:    f.Target.TargetId,
					state:        f.Target.Status.State,
					observations: observations,
				})
			}
		}
	}

	resultsPicker := viewport.New(pickerWidth, pickerHeight)
	resultsPicker.Style = common.OverlayStyle

	f := blist.New(findings, common.NewUnfocusedDelegate(), width, height)
	findingPicker := viewport.New(width, height)
	findingPicker.Style = common.PanelStyle

	findingSummary := viewport.New(width, height)
	findingSummary.Style = common.PanelStyle
	observationSummary := viewport.New(width, height)
	observationSummary.Style = common.PanelStyle

	help := common.NewHelpModel(false)
	help.OneLine = true
	help.ShortHelp = []key.Binding{assessmentHotkeys.Help}

	return Model{
		keys:               assessmentHotkeys,
		help:               help,
		results:            results,
		resultsPicker:      resultsPicker,
		selectedResult:     selectedResult,
		findings:           f,
		findingPicker:      findingPicker,
		findingSummary:     findingSummary,
		observationSummary: observationSummary,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateSizing(msg.Height-common.TabOffset, msg.Width)

	case tea.KeyMsg:
		if m.open {
			k := msg.String()
			switch k {

			case common.ContainsKey(k, m.keys.Quit.Keys()):
				return m, tea.Quit

			case common.ContainsKey(k, m.keys.Help.Keys()):
				m.help.ShowAll = !m.help.ShowAll

			case common.ContainsKey(k, m.keys.NavigateLeft.Keys()):
				if m.focus == 0 {
					m.focus = maxFocus
				} else {
					m.focus--
				}
				m.updateKeyBindings()

			case common.ContainsKey(k, m.keys.NavigateRight.Keys()):
				m.focus = (m.focus + 1) % (maxFocus + 1)
				m.updateKeyBindings()

			case common.ContainsKey(k, m.keys.Up.Keys()):
				if m.inResultOverlay && m.selectedResultIndex > 0 {
					m.selectedResultIndex--
					m.resultsPicker.SetContent(m.updateViewportContent("view"))
				}

			case common.ContainsKey(k, m.keys.Down.Keys()):
				if m.inResultOverlay && m.selectedResultIndex < len(m.results)-1 {
					m.selectedResultIndex++
					m.resultsPicker.SetContent(m.updateViewportContent("view"))
				}

			case common.ContainsKey(k, m.keys.Confirm.Keys()):
				if m.focus == focusResultSelection {
					if m.inResultOverlay {
						if len(m.results) > 1 {
							m.selectedResult = m.results[m.selectedResultIndex]
						}
						m.inResultOverlay = false
					} else {
						m.inResultOverlay = true
						m.resultsPicker.SetContent(m.updateViewportContent("view"))
					}
				} else if m.focus == focusCompareSelection {
					if m.inResultOverlay {
						if len(m.results) > 1 {
							m.compareResult = m.results[m.selectedResultIndex]
						}
						m.inResultOverlay = false
					} else {
						m.inResultOverlay = true
						m.resultsPicker.SetContent(m.updateViewportContent("compare"))
					}
				} else if m.focus == focusFindings {
					m.findingSummary.SetContent(m.renderSummary())
				}

			case common.ContainsKey(k, m.keys.Cancel.Keys()):
				if m.inResultOverlay {
					m.inResultOverlay = false
				}
			}
		}
	}
	m.findings, cmd = m.findings.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.inResultOverlay {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.resultsPicker.View(), lipgloss.WithWhitespaceChars(" "))
	}
	return m.mainView()
}

func (m Model) mainView() string {
	// Add help panel at the top left
	helpStyle := common.HelpStyle(m.width)
	helpView := helpStyle.Render(m.help.View())

	// Add viewport styles
	focusedViewport := common.PanelStyle.BorderForeground(common.Focused)
	focusedViewportHeaderColor := common.Focused
	focusedDialogBox := common.DialogBoxStyle.BorderForeground(common.Focused)

	selectedResultDialogBox := common.DialogBoxStyle
	compareResultDialogBox := common.DialogBoxStyle
	findingsViewport := common.PanelStyle
	findingsViewportHeader := common.Highlight
	summaryViewport := common.PanelStyle
	summaryViewportHeader := common.Highlight
	observationsViewport := common.PanelStyle
	observationsViewportHeader := common.Highlight

	switch m.focus {
	case focusResultSelection:
		selectedResultDialogBox = focusedDialogBox
	case focusCompareSelection:
		compareResultDialogBox = focusedDialogBox
	case focusFindings:
		findingsViewport = focusedViewport
		findingsViewportHeader = focusedViewportHeaderColor
	case focusSummary:
		summaryViewport = focusedViewport
		summaryViewportHeader = focusedViewportHeaderColor
	case focusObservations:
		observationsViewport = focusedViewport
		observationsViewportHeader = focusedViewportHeaderColor
	}

	// add panels at the top for selecting a result, selecting a comparison result
	const dialogFixedWidth = 40

	selectedResultLabel := common.LabelStyle.Render("Selected Result")
	selectedResultText := common.TruncateText(getResultText(m.selectedResult), dialogFixedWidth)
	selectedResultContent := selectedResultDialogBox.Width(dialogFixedWidth).Render(selectedResultText)
	selectedResult := lipgloss.JoinHorizontal(lipgloss.Top, selectedResultLabel, selectedResultContent)

	compareResultLabel := common.LabelStyle.Render("Compare Result")
	compareResultText := common.TruncateText(getResultText(m.compareResult), dialogFixedWidth)
	compareResultContent := compareResultDialogBox.Width(dialogFixedWidth).Render(compareResultText)
	compareResult := lipgloss.JoinHorizontal(lipgloss.Top, compareResultLabel, compareResultContent)

	resultSelectionContent := lipgloss.JoinHorizontal(lipgloss.Top, selectedResult, compareResult)

	// Add Controls panel + Results Tables
	m.findings.SetShowTitle(false)

	m.findingPicker.Style = findingsViewport
	m.findingPicker.SetContent(m.findings.View())
	bottomLeftView := fmt.Sprintf("%s\n%s", common.HeaderView("Findings List", m.findingPicker.Width-common.PanelStyle.GetMarginRight(), findingsViewportHeader), m.findingPicker.View())

	m.findingSummary.Style = summaryViewport
	m.findingSummary.SetContent(m.renderSummary())
	summaryPanel := fmt.Sprintf("%s\n%s", common.HeaderView("Summary", m.findingSummary.Width-common.PanelStyle.GetPaddingRight(), summaryViewportHeader), m.findingSummary.View())

	m.observationSummary.Style = observationsViewport
	m.observationSummary.SetContent(m.renderObservations())
	observationsPanel := fmt.Sprintf("%s\n%s", common.HeaderView("Observations", m.observationSummary.Width-common.PanelStyle.GetPaddingRight(), observationsViewportHeader), m.observationSummary.View())

	bottomRightView := lipgloss.JoinVertical(lipgloss.Top, summaryPanel, observationsPanel)
	bottomContent := lipgloss.JoinHorizontal(lipgloss.Top, bottomLeftView, bottomRightView)

	return lipgloss.JoinVertical(lipgloss.Top, helpView, resultSelectionContent, bottomContent)
}

func (m Model) updateViewportContent(resultType string) string {
	// TODO: refactor this to use the PiickerModel
	help := common.NewHelpModel(true)
	help.ShortHelp = common.ShortHelpPicker
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Select a result to %s:\n\n", resultType))

	for i, result := range m.results {
		if m.selectedResultIndex == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(getResultText(result))
		s.WriteString("\n")
	}

	return lipgloss.JoinVertical(lipgloss.Top, s.String(), help.View())
}

func (m Model) renderSummary() string {
	return "⚠️ Summary Under Construction ⚠️"
}

func (m Model) renderObservations() string {
	return "⚠️ Observations Under Construction ⚠️"
}

func getResultText(result result) string {
	if result.uuid == "" {
		return "No Result Selected"
	}
	return fmt.Sprintf("%s - %s", result.title, result.uuid)
}

func makeObservationMap(observations *[]oscalTypes_1_1_2.Observation) map[string]observation {
	observationMap := make(map[string]observation)

	for _, o := range *observations {
		validationId := findUuid(o.Description)
		state := "not-satisfied"
		remarks := strings.Builder{}
		if o.RelevantEvidence != nil {
			for _, re := range *o.RelevantEvidence {
				if re.Description == "Result: satisfied\n" {
					state = "satisfied"
				} else if re.Description == "Result: not-satisfied\n" {
					state = "not-satisfied"
				}
				remarks.WriteString(re.Remarks)
			}
		}
		observationMap[o.UUID] = observation{
			uuid:         o.UUID,
			description:  o.Description,
			remarks:      remarks.String(),
			state:        state,
			validationId: validationId,
		}
	}
	return observationMap
}

func findUuid(input string) string {
	uuidPattern := `[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`

	re := regexp.MustCompile(uuidPattern)

	return re.FindString(input)
}
