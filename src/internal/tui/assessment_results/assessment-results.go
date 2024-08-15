package assessmentresults

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	tui "github.com/defenseunicorns/lula/src/internal/tui/common"
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
	resultsPicker.Style = tui.OverlayStyle

	f := blist.New(findings, tui.NewUnfocusedDelegate(), width, height)
	findingPicker := viewport.New(width, height)
	findingPicker.Style = tui.PanelStyle

	findingSummary := viewport.New(width, height)
	findingSummary.Style = tui.PanelStyle
	observationSummary := viewport.New(width, height)
	observationSummary.Style = tui.PanelStyle

	return Model{
		content:            "Assessment Results Content",
		results:            results,
		resultsPicker:      resultsPicker,
		selectedResult:     selectedResult,
		findings:           f,
		findingPicker:      findingPicker,
		findingSummary:     findingSummary,
		observationSummary: observationSummary,
		keys:               assessmentHotkeys,
		help:               help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.open {
			switch msg.String() {
			case "r":
				if !m.inResultOverlay {
					m.inResultOverlay = true
					m.resultsPicker.SetContent(m.updateViewportContent("view"))
				} else {
					m.inResultOverlay = false
				}
			case "enter":
				if m.focus == focusResultSelection {
					if m.inResultOverlay {
						m.selectedResult = m.results[m.selectedResultIndex]
						m.inResultOverlay = false
					} else {
						m.inResultOverlay = true
						m.resultsPicker.SetContent(m.updateViewportContent("view"))
					}
				} else if m.focus == focusCompareSelection {
					if m.inResultOverlay {
						m.compareResult = m.results[m.selectedResultIndex]
						m.inResultOverlay = false
					} else {
						m.inResultOverlay = true
						m.resultsPicker.SetContent(m.updateViewportContent("compare"))
					}
				} else if m.focus == focusFindings {
					m.findingSummary.SetContent(m.renderSummary())
				}
			case "up":
				if m.inResultOverlay && m.selectedResultIndex > 0 {
					m.selectedResultIndex--
					m.resultsPicker.SetContent(m.updateViewportContent("view"))
				}
			case "down":
				if m.inResultOverlay && m.selectedResultIndex < len(m.results)-1 {
					m.selectedResultIndex++
					m.resultsPicker.SetContent(m.updateViewportContent("view"))
				}
			case "left":
				if m.focus != 0 {
					m.focus--
					m.updateKeyBindings()
				}
			case "right":
				if m.focus <= focusObservations {
					m.focus++
					m.updateKeyBindings()
				}
			case "esc", "q":
				if m.inResultOverlay {
					m.inResultOverlay = false
				}
			case "?":
				m.help.ShowAll = !m.help.ShowAll
			case "ctrl+c":
				return m, tea.Quit
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
	totalHeight := m.height
	leftWidth := m.width / 4
	rightWidth := m.width - leftWidth - tui.PanelStyle.GetHorizontalPadding() - tui.PanelStyle.GetHorizontalMargins()

	// Add help panel at the top left
	helpStyle := lipgloss.NewStyle().Align(lipgloss.Right).Width(m.width - tui.PanelStyle.GetHorizontalPadding() - tui.PanelStyle.GetHorizontalMargins()).Height(1)
	helpView := helpStyle.Render(m.help.View(m.keys))

	// Add viewport styles
	focusedViewport := tui.PanelStyle.BorderForeground(tui.Focused)
	focusedViewportHeaderColor := tui.Focused
	focusedDialogBox := tui.DialogBoxStyle.BorderForeground(tui.Focused)

	selectedResultDialogBox := tui.DialogBoxStyle
	compareResultDialogBox := tui.DialogBoxStyle
	findingsViewport := tui.PanelStyle
	findingsViewportHeader := tui.Highlight
	summaryViewport := tui.PanelStyle
	summaryViewportHeader := tui.Highlight
	observationsViewport := tui.PanelStyle
	observationsViewportHeader := tui.Highlight

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

	selectedResultLabel := tui.LabelStyle.Render("Selected Result")
	selectedResultText := tui.TruncateText(getResultText(m.selectedResult), dialogFixedWidth)
	selectedResultContent := selectedResultDialogBox.Width(dialogFixedWidth).Render(selectedResultText)
	selectedResult := lipgloss.JoinHorizontal(lipgloss.Top, selectedResultLabel, selectedResultContent)

	compareResultLabel := tui.LabelStyle.Render("Compare Result")
	compareResultText := tui.TruncateText(getResultText(m.compareResult), dialogFixedWidth)
	compareResultContent := compareResultDialogBox.Width(dialogFixedWidth).Render(compareResultText)
	compareResult := lipgloss.JoinHorizontal(lipgloss.Top, compareResultLabel, compareResultContent)

	resultSelectionContent := lipgloss.JoinHorizontal(lipgloss.Top, selectedResult, compareResult)

	// Add Controls panel + Results Table
	topSectionHeight := helpStyle.GetHeight() + tui.DialogBoxStyle.GetHeight()
	bottomHeight := totalHeight - topSectionHeight
	m.findings.SetShowTitle(false)
	m.findings.SetHeight(totalHeight - topSectionHeight - tui.PanelTitleStyle.GetHeight() - tui.PanelStyle.GetVerticalPadding())
	m.findings.SetWidth(leftWidth - tui.PanelStyle.GetHorizontalPadding())

	m.findingPicker.Style = findingsViewport
	m.findingPicker.SetContent(m.findings.View())
	m.findingPicker.Height = bottomHeight
	m.findingPicker.Width = leftWidth - tui.PanelStyle.GetHorizontalPadding()
	bottomLeftView := fmt.Sprintf("%s\n%s", tui.HeaderView("Findings List", m.findingPicker.Width-tui.PanelStyle.GetMarginRight(), findingsViewportHeader), m.findingPicker.View())

	// Add Summary and Observations panels
	bottomRightPanelHeight := (bottomHeight - 2*tui.PanelTitleStyle.GetHeight() - 2*tui.PanelTitleStyle.GetVerticalMargins()) / 2
	// summaryHeight := bottomHeight / 4
	m.findingSummary.Style = summaryViewport
	m.findingSummary.SetContent(m.renderSummary())
	m.findingSummary.Height = bottomRightPanelHeight
	m.findingSummary.Width = rightWidth
	summaryPanel := fmt.Sprintf("%s\n%s", tui.HeaderView("Summary", rightWidth-tui.PanelStyle.GetPaddingRight(), summaryViewportHeader), m.findingSummary.View()) // TODO: fix padding

	// observationsHeight := bottomHeight - summaryHeight - tui.PanelStyle.GetVerticalPadding() - tui.PanelStyle.GetPaddingBottom() - 2*tui.PanelTitleStyle.GetHeight()
	m.observationSummary.Style = observationsViewport
	m.observationSummary.SetContent(m.renderObservations())
	m.observationSummary.Height = bottomRightPanelHeight
	m.observationSummary.Width = rightWidth
	observationsPanel := fmt.Sprintf("%s\n%s", tui.HeaderView("Observations", rightWidth-tui.PanelStyle.GetPaddingRight(), observationsViewportHeader), m.observationSummary.View())

	bottomRightView := lipgloss.JoinVertical(lipgloss.Top, summaryPanel, observationsPanel)
	bottomContent := lipgloss.JoinHorizontal(lipgloss.Top, bottomLeftView, bottomRightView)

	return lipgloss.JoinVertical(lipgloss.Top, helpView, resultSelectionContent, bottomContent)
}

func (m Model) updateViewportContent(resultType string) string {
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

	return s.String()
}

func (m Model) renderSummary() string {
	return "Summary"
}

func (m Model) renderObservations() string {
	return "Observations"
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
