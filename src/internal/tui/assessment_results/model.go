package assessmentresults

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/internal/tui/component"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/evertras/bubble-table/table"
)

const (
	height           = 20
	width            = 12
	pickerHeight     = 20
	pickerWidth      = 80
	dialogFixedWidth = 40
)

const (
	resultPicker                 common.PickerKind = "result"
	comparedResultPicker         common.PickerKind = "compared result"
	ColumnKeyName                                  = "name"
	ColumnKeyStatus                                = "status"
	ColumnKeyDescription                           = "description"
	ColumnKeyStatusChange                          = "status_change"
	ColumnKeyFinding                               = "finding"
	ColumnKeyRelatedObs                            = "related_obs"
	ColumnKeyControlIds                            = "control_ids"
	ColumnKeyComparedFinding                       = "compared_finding"
	ColumnKeyObservation                           = "observation"
	ColumnKeyComparedObservation                   = "compared_observation"
	ColumnKeyValidationId                          = "validation_id"
)

type Model struct {
	open                  bool
	help                  common.HelpModel
	keys                  keys
	focus                 focus
	oscalAssessment       *oscalTypes.AssessmentResults
	results               []result
	resultsPicker         common.PickerModel
	selectedResult        result
	selectedResultIndex   int
	comparedResultsPicker common.PickerModel
	comparedResult        result
	findingsSummary       viewport.Model
	findingsTable         table.Model
	observationsSummary   viewport.Model
	observationsTable     table.Model
	currentObservations   []table.Row
	detailView            common.DetailModel
	width                 int
	height                int
}

type ModelOpenMsg struct {
	Height int
	Width  int
}
type ModelCloseMsg struct{}
type AssessmentUpdatedMsg struct{}

func NewAssessmentResultsModel(assessmentResults *oscalTypes.AssessmentResults) Model {
	help := common.NewHelpModel(false)
	help.OneLine = true
	help.ShortHelp = shortHelpNoFocus

	resultsPicker := common.NewPickerModel("Select a Result", resultPicker, []string{}, 0)
	comparedResultsPicker := common.NewPickerModel("Select a Result to Compare", comparedResultPicker, []string{}, 0)

	findingsSummary := viewport.New(width, height)
	findingsSummary.Style = common.PanelStyle
	observationsSummary := viewport.New(width, height)
	observationsSummary.Style = common.PanelStyle

	model := Model{
		keys:                  assessmentKeys,
		help:                  help,
		oscalAssessment:       assessmentResults,
		resultsPicker:         resultsPicker,
		comparedResultsPicker: comparedResultsPicker,
		findingsSummary:       findingsSummary,
		observationsSummary:   observationsSummary,
		detailView:            common.NewDetailModel(),
	}

	model.UpdateWithAssessmentResults(assessmentResults)

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ModelOpenMsg:
		m.Open(msg.Height, msg.Width)

	case tea.WindowSizeMsg:
		m.updateSizing(msg.Height-common.TabOffset, msg.Width)

	case tea.KeyMsg:

		common.DumpToLog(msg)
		k := msg.String()
		switch k {
		case common.ContainsKey(k, m.keys.Help.Keys()):
			m.help.ShowAll = !m.help.ShowAll

		case common.ContainsKey(k, m.keys.NavigateLeft.Keys()):
			if !m.inOverlay() {
				if m.focus == 0 {
					m.focus = maxFocus
				} else {
					m.focus--
				}
				m.updateKeyBindings()
			}

		case common.ContainsKey(k, m.keys.NavigateRight.Keys()):
			if !m.inOverlay() {
				m.focus = (m.focus + 1) % (maxFocus + 1)
				m.updateKeyBindings()
			}

		case common.ContainsKey(k, m.keys.Confirm.Keys()):
			m.keys = assessmentKeys
			switch m.focus {
			case focusResultSelection:
				if len(m.results) > 0 && !m.resultsPicker.Open {
					return m, func() tea.Msg {
						return common.PickerOpenMsg{
							Kind: resultPicker,
						}
					}
				}

			case focusCompareSelection:
				if len(m.results) > 0 && !m.comparedResultsPicker.Open {
					// TODO: get compared result items to send with picker open
					return m, func() tea.Msg {
						return common.PickerOpenMsg{
							Kind: comparedResultPicker,
						}
					}
				}

			case focusFindings:
				// Select the observations
				if !m.detailView.Open && m.findingsTable.HighlightedRow().Data != nil {
					m.observationsTable = m.observationsTable.WithRows(m.getObservationsByFinding(m.findingsTable.HighlightedRow().Data[ColumnKeyRelatedObs].([]string)))
				}
			}

		case common.ContainsKey(k, m.keys.Cancel.Keys()):
			m.keys = assessmentKeys
			switch m.focus {
			case focusFindings:
				m.observationsTable = m.observationsTable.WithRows(m.currentObservations)
			}
			m.updateKeyBindings()

		case common.ContainsKey(k, m.keys.Detail.Keys()):
			switch m.focus {
			case focusFindings:
				if m.findingsTable.HighlightedRow().Data != nil {
					m.findingsTable = m.findingsTable.WithKeyMap(common.UnfocusedTableKeyMap())
					return m, func() tea.Msg {
						return common.DetailOpenMsg{
							Content:      m.getFindingsDetail(),
							WindowHeight: m.height,
							WindowWidth:  m.width,
						}
					}
				}

			case focusObservations:
				if m.observationsTable.HighlightedRow().Data != nil {
					m.observationsTable = m.observationsTable.WithKeyMap(common.UnfocusedTableKeyMap())
					return m, func() tea.Msg {
						return common.DetailOpenMsg{
							Content:      m.getObsDetail(),
							WindowHeight: m.height,
							WindowWidth:  m.width,
						}
					}
				}
			}

		case common.ContainsKey(k, m.keys.Filter.Keys()):
			// Lock keys during table filter
			if m.focus == focusFindings && !m.detailView.Open {
				m.keys = assessmentKeysInFilter
			}
			if m.focus == focusObservations && !m.detailView.Open {
				m.keys = assessmentKeysInFilter
			}
		}

	case common.PickerItemSelected:
		if msg.From == resultPicker {
			m.selectedResultIndex = msg.Selected
			m.selectedResult = m.results[m.selectedResultIndex]
			m.findingsTable, m.observationsTable = m.getSingleResultTables(m.selectedResult.FindingsRows, m.selectedResult.ObservationsRows)
			m.currentObservations = m.selectedResult.ObservationsRows
			// Update comparison
			m.comparedResult = result{}
			m.comparedResultsPicker.UpdateItems(getComparedResults(m.results, m.selectedResult))
		} else if msg.From == comparedResultPicker {
			// First item will always be "None", so return single table if selected
			if msg.Selected == 0 {
				if m.comparedResult.OscalResult != nil {
					m.findingsTable, m.observationsTable = m.getSingleResultTables(m.selectedResult.FindingsRows, m.selectedResult.ObservationsRows)
					m.currentObservations = m.selectedResult.ObservationsRows
					m.comparedResult = result{}
				}
			} else {
				if m.selectedResultIndex < msg.Selected {
					m.comparedResult = m.results[msg.Selected]
				} else {
					m.comparedResult = m.results[msg.Selected-1]
				}
				m.findingsTable, m.observationsTable, m.currentObservations = m.getComparedResultTables(m.selectedResult, m.comparedResult)
			}
		}

	case component.ValidationDataMsg:
		err := m.MergeNewResults(msg.AssessmentResults)
		if err != nil {
			common.PrintToLog("error merging assessment results")
		} else {
			cmds = append(cmds, func() tea.Msg {
				return AssessmentUpdatedMsg{}
			})
		}
	}

	mdl, cmd := m.resultsPicker.Update(msg)
	m.resultsPicker = mdl.(common.PickerModel)
	cmds = append(cmds, cmd)

	mdl, cmd = m.comparedResultsPicker.Update(msg)
	m.comparedResultsPicker = mdl.(common.PickerModel)
	cmds = append(cmds, cmd)

	mdl, cmd = m.detailView.Update(msg)
	m.detailView = mdl.(common.DetailModel)
	cmds = append(cmds, cmd)

	m.findingsTable, cmd = m.findingsTable.Update(msg)
	cmds = append(cmds, cmd)

	m.observationsTable, cmd = m.observationsTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.resultsPicker.Open {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.resultsPicker.View(), lipgloss.WithWhitespaceChars(" "))
	}
	if m.comparedResultsPicker.Open {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.comparedResultsPicker.View(), lipgloss.WithWhitespaceChars(" "))
	}
	if m.detailView.Open {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.detailView.View(), lipgloss.WithWhitespaceChars(" "))
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
	comparedResultDialogBox := common.DialogBoxStyle
	findingsViewport := common.PanelStyle
	findingsViewportHeader := common.Highlight
	findingsTableStyle := common.TableStyleBase
	observationsViewport := common.PanelStyle
	observationsViewportHeader := common.Highlight
	observationsTableStyle := common.TableStyleBase

	switch m.focus {
	case focusResultSelection:
		selectedResultDialogBox = focusedDialogBox
	case focusCompareSelection:
		comparedResultDialogBox = focusedDialogBox
	case focusFindings:
		findingsViewport = focusedViewport
		findingsViewportHeader = focusedViewportHeaderColor
		findingsTableStyle = common.TableStyleActive
	case focusObservations:
		observationsViewport = focusedViewport
		observationsViewportHeader = focusedViewportHeaderColor
		observationsTableStyle = common.TableStyleActive
	}

	// add panels at the top for selecting a result, selecting a comparison result
	const dialogFixedWidth = 40

	selectedResultLabel := common.LabelStyle.Render("Selected Result")
	selectedResultText := common.TruncateText(getResultText(m.selectedResult), dialogFixedWidth)
	selectedResultContent := selectedResultDialogBox.Width(dialogFixedWidth).Render(selectedResultText)
	selectedResult := lipgloss.JoinHorizontal(lipgloss.Top, selectedResultLabel, selectedResultContent)

	comparedResultLabel := common.LabelStyle.Render("Compare Result")
	comparedResultText := common.TruncateText(getResultText(m.comparedResult), dialogFixedWidth)
	comparedResultContent := comparedResultDialogBox.Width(dialogFixedWidth).Render(comparedResultText)
	comparedResult := lipgloss.JoinHorizontal(lipgloss.Top, comparedResultLabel, comparedResultContent)

	resultSelectionContent := lipgloss.JoinHorizontal(lipgloss.Top, selectedResult, comparedResult)

	// Write summary
	findingsSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#3ad33c")).Render(fmt.Sprintf("%d", m.selectedResult.SummaryData.NumFindingsSatisfied))
	findingsNotSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#e36750")).Render(fmt.Sprintf("%d", m.selectedResult.SummaryData.NumFindings-m.selectedResult.SummaryData.NumFindingsSatisfied))
	observationsSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#3ad33c")).Render(fmt.Sprintf("%d", m.selectedResult.SummaryData.NumObservationsSatisfied))
	observationsNotSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#e36750")).Render(fmt.Sprintf("%d", m.selectedResult.SummaryData.NumObservations-m.selectedResult.SummaryData.NumObservationsSatisfied))
	summaryText := fmt.Sprintf("Summary: %d (%s/%s) Findings - %d (%s/%s) Observations",
		m.selectedResult.SummaryData.NumFindings, findingsSatisfied, findingsNotSatisfied,
		m.selectedResult.SummaryData.NumObservations, observationsSatisfied, observationsNotSatisfied,
	)

	// Write compared summary
	if m.comparedResult.OscalResult != nil {
		comparedFindingsSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#3ad33c")).Render(fmt.Sprintf("%d", m.comparedResult.SummaryData.NumFindingsSatisfied))
		comparedFindingsNotSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#e36750")).Render(fmt.Sprintf("%d", m.comparedResult.SummaryData.NumFindings-m.comparedResult.SummaryData.NumFindingsSatisfied))
		comparedObservationsSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#3ad33c")).Render(fmt.Sprintf("%d", m.comparedResult.SummaryData.NumObservationsSatisfied))
		comparedObservationsNotSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#e36750")).Render(fmt.Sprintf("%d", m.comparedResult.SummaryData.NumObservations-m.comparedResult.SummaryData.NumObservationsSatisfied))
		summaryText += fmt.Sprintf(" | Compared Summary: %d (%s/%s) Findings - %d (%s/%s) Observations",
			m.comparedResult.SummaryData.NumFindings, comparedFindingsSatisfied, comparedFindingsNotSatisfied,
			m.comparedResult.SummaryData.NumObservations, comparedObservationsSatisfied, comparedObservationsNotSatisfied,
		)
	}

	summary := lipgloss.JoinHorizontal(lipgloss.Top, common.SummaryTextStyle.Render(summaryText))

	// Add Tables
	m.findingsSummary.Style = findingsViewport
	m.findingsTable = m.findingsTable.WithBaseStyle(findingsTableStyle)
	m.findingsSummary.SetContent(m.findingsTable.View())
	findingsPanel := fmt.Sprintf("%s\n%s", common.HeaderView("Findings", m.findingsSummary.Width-common.PanelStyle.GetPaddingRight(), findingsViewportHeader), m.findingsSummary.View())

	m.observationsSummary.Style = observationsViewport
	m.observationsTable = m.observationsTable.WithBaseStyle(observationsTableStyle)
	m.observationsSummary.SetContent(m.observationsTable.View())
	observationsPanel := fmt.Sprintf("%s\n%s", common.HeaderView("Observations", m.observationsSummary.Width-common.PanelStyle.GetPaddingRight(), observationsViewportHeader), m.observationsSummary.View())

	bottomContent := lipgloss.JoinVertical(lipgloss.Top, summary, findingsPanel, observationsPanel)

	return lipgloss.JoinVertical(lipgloss.Top, helpView, resultSelectionContent, bottomContent)
}

func (m *Model) Close() {
	m.open = false
}

func (m *Model) Open(height, width int) {
	m.open = true
	m.updateSizing(height, width)
}

func (m *Model) GetAssessmentResults() *oscalTypes.AssessmentResults {
	return m.oscalAssessment
}

func (m *Model) UpdateWithAssessmentResults(assessmentResults *oscalTypes.AssessmentResults) {
	var selectedResult result
	m.oscalAssessment = assessmentResults

	results := GetResults(assessmentResults)

	if len(results) != 0 {
		selectedResult = results[0]
	}

	// Update model parameters
	resultItems := make([]string, len(results))
	for i, c := range results {
		resultItems[i] = getResultText(c)
	}

	m.results = results
	m.selectedResult = selectedResult
	m.resultsPicker.UpdateItems(resultItems)
	m.comparedResultsPicker.UpdateItems(getComparedResults(results, selectedResult))

	m.findingsTable, m.observationsTable = m.getSingleResultTables(selectedResult.FindingsRows, selectedResult.ObservationsRows)
	m.currentObservations = selectedResult.ObservationsRows
}

func (m *Model) MergeNewResults(newResults *oscalTypes.AssessmentResults) error {
	if newResults != nil {
		if m.oscalAssessment == nil {
			m.UpdateWithAssessmentResults(newResults)
			return nil
		}
		combinedAssessmentResults, err := oscal.MergeAssessmentResults(m.oscalAssessment, newResults)
		if err != nil {
			return fmt.Errorf("error merging assessment results: %v", err)
		}
		m.UpdateWithAssessmentResults(combinedAssessmentResults)
	}
	return nil
}

func (m *Model) updateSizing(height, width int) {
	m.height = height
	m.width = width
	totalHeight := m.height - 4 // Subtract vertical margin/padding not accounted for below

	topSectionHeight := common.HelpStyle(m.width).GetHeight() + common.DialogBoxStyle.GetHeight()
	bottomSectionHeight := totalHeight - topSectionHeight - 2 // 2 for summary height
	bottomPanelHeight := (bottomSectionHeight - 2*common.PanelTitleStyle.GetHeight() - 2*common.PanelTitleStyle.GetVerticalMargins()) / 2
	panelWidth := width - 4
	panelInternalWidth := panelWidth - common.PanelStyle.GetHorizontalPadding() - common.PanelStyle.GetHorizontalMargins() - 2

	// Update widget dimensions
	m.findingsSummary.Height = bottomPanelHeight
	m.findingsSummary.Width = panelWidth
	findingsRowHeight := bottomPanelHeight - common.PanelTitleStyle.GetHeight() - common.PanelStyle.GetVerticalPadding() - 6
	m.findingsTable = m.findingsTable.WithTargetWidth(panelInternalWidth).WithPageSize(findingsRowHeight)
	m.observationsSummary.Height = bottomPanelHeight
	m.observationsSummary.Width = panelWidth
	observationsRowHeight := bottomPanelHeight - common.PanelTitleStyle.GetHeight() - common.PanelStyle.GetVerticalPadding() - 6
	m.observationsTable = m.observationsTable.WithTargetWidth(panelInternalWidth).WithPageSize(observationsRowHeight)
}

func (m *Model) inOverlay() bool {
	return m.resultsPicker.Open || m.comparedResultsPicker.Open || m.detailView.Open
}

func (m *Model) getObservationsByFinding(relatedObs []string) []table.Row {
	obsRows := make([]table.Row, 0)
	for _, uuid := range relatedObs {
		if obsRow, ok := m.selectedResult.ObservationsMap[uuid]; ok {
			obsRows = append(obsRows, obsRow)
		}
	}

	return obsRows
}

func (m *Model) getSingleResultTables(findingsRows, observationsRows []table.Row) (findingsTable table.Model, observationsTable table.Model) {
	findingsTableColumns := []table.Column{
		table.NewFlexColumn(ColumnKeyName, "Control", 1).WithFiltered(true),
		table.NewFlexColumn(ColumnKeyStatus, "Status", 1),
		table.NewFlexColumn(ColumnKeyDescription, "Description", 4),
	}

	observationsTableColumns := []table.Column{
		table.NewFlexColumn(ColumnKeyName, "Observation", 1).WithFiltered(true),
		table.NewFlexColumn(ColumnKeyStatus, "Status", 1),
		table.NewFlexColumn(ColumnKeyControlIds, "Controls", 1),
		table.NewFlexColumn(ColumnKeyDescription, "Remarks", 4),
	}

	findingsHeight, findingsWidth := getTableDimensions(m.findingsSummary.Height, m.findingsSummary.Width)
	observationsHeight, observationsWidth := getTableDimensions(m.observationsSummary.Height, m.observationsSummary.Width)

	findingsTable = table.New(findingsTableColumns).
		WithRows(findingsRows).
		WithBaseStyle(common.TableStyleBase).
		Filtered(true).
		SortByAsc(ColumnKeyName).
		WithTargetWidth(findingsWidth).
		WithPageSize(findingsHeight)

	observationsTable = table.New(observationsTableColumns).
		WithRows(observationsRows).
		WithBaseStyle(common.TableStyleBase).
		Filtered(true).
		SortByAsc(ColumnKeyName).
		WithTargetWidth(observationsWidth).
		WithPageSize(observationsHeight)

	return findingsTable, observationsTable
}

func (m *Model) getComparedResultTables(selectedResult, comparedResult result) (findingsTable table.Model, observationsTable table.Model, currentObservations []table.Row) {
	findingsRows, observationsRows := GetResultComparison(selectedResult, comparedResult)

	// Set up tables
	findingsTableColumns := []table.Column{
		table.NewFlexColumn(ColumnKeyName, "Control", 1).WithFiltered(true),
		table.NewFlexColumn(ColumnKeyStatus, "Status", 1),
		table.NewFlexColumn(ColumnKeyStatusChange, "Status Change", 1).WithFiltered(true),
		table.NewFlexColumn(ColumnKeyDescription, "Description", 4),
	}

	observationsTableColumns := []table.Column{
		table.NewFlexColumn(ColumnKeyName, "Observation", 1).WithFiltered(true),
		table.NewFlexColumn(ColumnKeyStatus, "Status", 1),
		table.NewFlexColumn(ColumnKeyStatusChange, "Status Change", 1).WithFiltered(true),
		table.NewFlexColumn(ColumnKeyControlIds, "Controls", 1).WithFiltered(true),
		table.NewFlexColumn(ColumnKeyDescription, "Remarks", 4),
	}

	findingsHeight, findingsWidth := getTableDimensions(m.findingsSummary.Height, m.findingsSummary.Width)
	observationsHeight, observationsWidth := getTableDimensions(m.observationsSummary.Height, m.observationsSummary.Width)

	findingsTable = table.New(findingsTableColumns).
		WithRows(findingsRows).
		WithBaseStyle(common.TableStyleBase).
		Filtered(true).
		SortByAsc(ColumnKeyName).
		WithTargetWidth(findingsWidth).
		WithPageSize(findingsHeight)

	observationsTable = table.New(observationsTableColumns).
		WithRows(observationsRows).
		WithBaseStyle(common.TableStyleBase).
		Filtered(true).
		SortByAsc(ColumnKeyName).
		WithTargetWidth(observationsWidth).
		WithPageSize(observationsHeight)

	return findingsTable, observationsTable, observationsRows
}

func getTableDimensions(parentHeight, parentWidth int) (height int, width int) {
	height = parentHeight - common.PanelTitleStyle.GetHeight() - common.PanelStyle.GetVerticalPadding() - 6
	width = parentWidth - common.PanelStyle.GetHorizontalPadding() - common.PanelStyle.GetHorizontalMargins() - 2
	return height, width
}

func (m *Model) getFindingsDetail() string {
	var text strings.Builder
	important := lipgloss.NewStyle().Bold(true).
		Foreground(common.Special)

	text.WriteString(fmt.Sprintf("%s\n\n", important.Render("Finding: "+m.findingsTable.HighlightedRow().Data[ColumnKeyName].(string))))
	text.WriteString(m.findingsTable.HighlightedRow().Data[ColumnKeyFinding].(string))

	if m.comparedResult.OscalResult != nil {
		text.WriteString(fmt.Sprintf("\n\n%s\n\n", important.Render("Compared Finding: "+m.findingsTable.HighlightedRow().Data[ColumnKeyName].(string))))
		text.WriteString(m.findingsTable.HighlightedRow().Data[ColumnKeyComparedFinding].(string))
	}

	return text.String()
}

func (m *Model) getObsDetail() string {
	var text strings.Builder
	important := lipgloss.NewStyle().Bold(true).
		Foreground(common.Special)

	text.WriteString(fmt.Sprintf("Control IDs: %s\n\n", m.observationsTable.HighlightedRow().Data[ColumnKeyControlIds].(string)))
	text.WriteString(fmt.Sprintf("%s\n\n", important.Render("Observation: "+m.observationsTable.HighlightedRow().Data[ColumnKeyName].(string))))
	text.WriteString(m.observationsTable.HighlightedRow().Data[ColumnKeyObservation].(string))

	if m.comparedResult.OscalResult != nil {
		text.WriteString(fmt.Sprintf("\n\n%s\n\n", important.Render("Compared Observation: "+m.observationsTable.HighlightedRow().Data[ColumnKeyName].(string))))
		text.WriteString(m.observationsTable.HighlightedRow().Data[ColumnKeyComparedObservation].(string))
	}

	return text.String()
}
