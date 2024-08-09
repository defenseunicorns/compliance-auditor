package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/mattn/go-runewidth"
)

type assessmentResultsModel struct {
	focus               arFocus
	content             string
	keys                assessmentKeys
	results             []result
	findings            blist.Model
	findingPicker       viewport.Model
	findingSummary      viewport.Model
	observationSummary  viewport.Model
	inResultOverlay     bool
	selectedResult      result
	selectedResultIndex int
	compareResult       result
	compareResultIndex  int
	resultsPicker       viewport.Model
	open                bool
	help                help.Model
	width               int
	height              int
	pickerWidth         int
	pickerHeight        int
}

type arFocus int

const (
	noFocus arFocus = iota
	focusResultSelection
	focusCompareSelection
	focusFindings
	focusSummary
	focusObservations
)

type result struct {
	uuid, title  string
	findings     *[]oscalTypes_1_1_2.Finding
	observations *[]oscalTypes_1_1_2.Observation
}

type finding struct {
	title, uuid, controlId, state string
	observations                  []observation
}

func (i finding) Title() string       { return i.controlId }
func (i finding) Description() string { return i.state }
func (i finding) FilterValue() string { return i.title }

type observation struct {
	uuid, description, remarks, state, validationId string
}

func newUnfocusedDelegate() blist.DefaultDelegate {
	d := blist.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.NormalTitle
	d.Styles.SelectedDesc = d.Styles.NormalDesc

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{listHotkeys.Confirm, listHotkeys.Help}
	}

	return d
}

func newFocusedDelegate() blist.DefaultDelegate {
	d := blist.NewDefaultDelegate()

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{listHotkeys.Confirm, listHotkeys.Help}
	}

	return d
}

func NewAssessmentResultsModel(assessmentResults *oscalTypes_1_1_2.AssessmentResults) assessmentResultsModel {
	height := 60
	width := 12
	pickerHeight := 20
	pickerWidth := 80

	results := make([]result, 0)
	findings := make([]blist.Item, 0)

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
	var selectedResult result
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
	resultsPicker.Style = overlayStyle

	f := blist.New(findings, newUnfocusedDelegate(), width, height)
	findingPicker := viewport.New(width, height)
	findingPicker.Style = panelStyle

	findingSummary := viewport.New(width, height)
	findingSummary.Style = panelStyle
	observationSummary := viewport.New(width, height)
	observationSummary.Style = panelStyle

	return assessmentResultsModel{
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

func (m assessmentResultsModel) Init() tea.Cmd {
	return nil
}

func (m assessmentResultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			case "q":
				if m.inResultOverlay {
					m.inResultOverlay = false
				} else {
					return m, tea.Quit
				}
			case "?":
				m.help.ShowAll = !m.help.ShowAll
			}
		}
	}
	m.findings, cmd = m.findings.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m assessmentResultsModel) View() string {
	if m.inResultOverlay {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.resultsPicker.View(), lipgloss.WithWhitespaceChars(" "))
	}
	return m.mainView()
}

func (m assessmentResultsModel) mainView() string {
	totalHeight := m.height
	leftWidth := m.width / 4
	rightWidth := m.width - leftWidth - panelStyle.GetHorizontalPadding() - panelStyle.GetHorizontalMargins()

	// Add help panel at the top left
	helpStyle := lipgloss.NewStyle().Align(lipgloss.Right).Width(m.width - panelStyle.GetHorizontalPadding() - panelStyle.GetHorizontalMargins()).Height(1)
	helpView := helpStyle.Render(m.help.View(m.keys))

	// Add viewport styles
	focusedViewport := panelStyle.BorderForeground(focused)
	focusedViewportHeaderColor := focused
	focusedDialogBox := dialogBoxStyle.BorderForeground(focused)

	selectedResultDialogBox := dialogBoxStyle
	compareResultDialogBox := dialogBoxStyle
	findingsViewport := panelStyle
	findingsViewportHeader := highlight
	summaryViewport := panelStyle
	summaryViewportHeader := highlight
	observationsViewport := panelStyle
	observationsViewportHeader := highlight

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

	selectedResultLabel := labelStyle.Render("Selected Result")
	selectedResultText := truncateText(getResultText(m.selectedResult), dialogFixedWidth)
	selectedResultContent := selectedResultDialogBox.Width(dialogFixedWidth).Render(selectedResultText)
	selectedResult := lipgloss.JoinHorizontal(lipgloss.Top, selectedResultLabel, selectedResultContent)

	compareResultLabel := labelStyle.Render("Compare Result")
	compareResultText := truncateText(getResultText(m.compareResult), dialogFixedWidth)
	compareResultContent := compareResultDialogBox.Width(dialogFixedWidth).Render(compareResultText)
	compareResult := lipgloss.JoinHorizontal(lipgloss.Top, compareResultLabel, compareResultContent)

	resultSelectionContent := lipgloss.JoinHorizontal(lipgloss.Top, selectedResult, compareResult)

	// Add Controls panel + Results Table
	topSectionHeight := helpStyle.GetHeight() + dialogBoxStyle.GetHeight()
	bottomHeight := totalHeight - topSectionHeight
	m.findings.SetShowTitle(false)
	m.findings.SetHeight(totalHeight - topSectionHeight - panelTitleStyle.GetHeight() - panelStyle.GetVerticalPadding())
	m.findings.SetWidth(leftWidth - panelStyle.GetHorizontalPadding())

	m.findingPicker.Style = findingsViewport
	m.findingPicker.SetContent(m.findings.View())
	m.findingPicker.Height = bottomHeight
	m.findingPicker.Width = leftWidth - panelStyle.GetHorizontalPadding()
	bottomLeftView := fmt.Sprintf("%s\n%s", headerView("Findings List", m.findingPicker.Width-panelStyle.GetMarginRight(), findingsViewportHeader), m.findingPicker.View())

	// Add Summary and Observations panels
	bottomRightPanelHeight := (bottomHeight - 2*panelTitleStyle.GetHeight() - 2*panelTitleStyle.GetVerticalMargins()) / 2
	// summaryHeight := bottomHeight / 4
	m.findingSummary.Style = summaryViewport
	m.findingSummary.SetContent(m.renderSummary())
	m.findingSummary.Height = bottomRightPanelHeight
	m.findingSummary.Width = rightWidth
	summaryPanel := fmt.Sprintf("%s\n%s", headerView("Summary", rightWidth-panelStyle.GetPaddingRight(), summaryViewportHeader), m.findingSummary.View()) // TODO: fix padding

	// observationsHeight := bottomHeight - summaryHeight - panelStyle.GetVerticalPadding() - panelStyle.GetPaddingBottom() - 2*panelTitleStyle.GetHeight()
	m.observationSummary.Style = observationsViewport
	m.observationSummary.SetContent(m.renderObservations())
	m.observationSummary.Height = bottomRightPanelHeight
	m.observationSummary.Width = rightWidth
	observationsPanel := fmt.Sprintf("%s\n%s", headerView("Observations", rightWidth-panelStyle.GetPaddingRight(), observationsViewportHeader), m.observationSummary.View())

	bottomRightView := lipgloss.JoinVertical(lipgloss.Top, summaryPanel, observationsPanel)
	bottomContent := lipgloss.JoinHorizontal(lipgloss.Top, bottomLeftView, bottomRightView)

	return lipgloss.JoinVertical(lipgloss.Top, helpView, resultSelectionContent, bottomContent)
}

func (m assessmentResultsModel) updateViewportContent(resultType string) string {
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

func (m assessmentResultsModel) renderSummary() string {
	return "Summary"
}

func (m assessmentResultsModel) renderObservations() string {
	return "Observations"
}

func (m *assessmentResultsModel) updateKeyBindings() {
	m.findings.KeyMap = unfocusedListKeyMap()
	m.findings.SetDelegate(newUnfocusedDelegate())

	switch m.focus {
	case focusFindings:
		m.findings.KeyMap = focusedListKeyMap()
		m.findings.SetDelegate(newFocusedDelegate())
	}
}

func getResultText(result result) string {
	if result.uuid == "" {
		return "No Result Selected"
	}
	return fmt.Sprintf("%s - %s", result.title, result.uuid)
}

func truncateText(text string, width int) string {
	if runewidth.StringWidth(text) <= width {
		return text
	}

	ellipsis := "…"
	trimmedWidth := width - runewidth.StringWidth(ellipsis)
	trimmedText := runewidth.Truncate(text, trimmedWidth, "")

	return trimmedText + ellipsis
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
