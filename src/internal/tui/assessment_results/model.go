package assessmentresults

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
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
	columnKeyName                                  = "name"
	columnKeyStatus                                = "status"
	columnKeyDescription                           = "description"
	columnKeyStatusChange                          = "status_change"
	columnKeyFinding                               = "finding"
	columnKeyRelatedObs                            = "related_obs"
	columnKeyComparedFinding                       = "compared_finding"
	columnKeyObservation                           = "observation"
	columnKeyComparedObservation                   = "compared_observation"
	columnKeyValidationId                          = "validation_id"
)

type Model struct {
	open                  bool
	help                  common.HelpModel
	keys                  keys
	focus                 focus
	results               []result
	resultsPicker         common.PickerModel
	selectedResult        result
	comparedResultsPicker common.PickerModel
	comparedResult        result
	findingsSummary       viewport.Model
	findingsTable         table.Model
	observationsSummary   viewport.Model
	observationsTable     table.Model
	detailView            common.DetailModel
	width                 int
	height                int
}

func NewAssessmentResultsModel(assessmentResults *oscalTypes_1_1_2.AssessmentResults) Model {
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
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateSizing(msg.Height-common.TabOffset, msg.Width)

	case tea.KeyMsg:
		if m.open {
			common.DumpToLog(msg)
			k := msg.String()
			switch k {
			case common.ContainsKey(k, m.keys.Help.Keys()):
				m.help.ShowAll = !m.help.ShowAll

			case common.ContainsKey(k, m.keys.NavigateLeft.Keys()):
				if !m.resultsPicker.Open && !m.comparedResultsPicker.Open && !m.detailView.Open {
					if m.focus == 0 {
						m.focus = maxFocus
					} else {
						m.focus--
					}
					m.updateKeyBindings()
				}

			case common.ContainsKey(k, m.keys.NavigateRight.Keys()):
				if !m.resultsPicker.Open && !m.comparedResultsPicker.Open && !m.detailView.Open {
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
					if len(m.results) > 1 && !m.comparedResultsPicker.Open {
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
						m.observationsTable = m.observationsTable.WithRows(m.getObservationsByFinding(m.findingsTable.HighlightedRow().Data[columnKeyRelatedObs].([]string)))
					}
				}

			case common.ContainsKey(k, m.keys.Cancel.Keys()):
				m.keys = assessmentKeys
				switch m.focus {
				case focusFindings:
					m.observationsTable = m.observationsTable.WithRows(m.selectedResult.observationsRows)
				}

			case common.ContainsKey(k, m.keys.Detail.Keys()):
				switch m.focus {
				case focusFindings:
					if m.findingsTable.HighlightedRow().Data != nil {
						selected := m.findingsTable.HighlightedRow().Data[columnKeyFinding].(string)
						return m, func() tea.Msg {
							return common.DetailOpenMsg{
								Content:      selected,
								WindowHeight: (m.height + common.TabOffset),
								WindowWidth:  m.width,
							}
						}
					}

				case focusObservations:
					if m.observationsTable.HighlightedRow().Data != nil {
						selected := m.observationsTable.HighlightedRow().Data[columnKeyObservation].(string)
						return m, func() tea.Msg {
							return common.DetailOpenMsg{
								Content:      selected,
								WindowHeight: (m.height + common.TabOffset),
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
		}

	case common.PickerItemSelected:
		if m.open {
			if msg.From == resultPicker {
				m.selectedResult = m.results[msg.Selected]
				m.observationsTable = m.observationsTable.WithRows(m.selectedResult.observationsRows)
				m.findingsTable = m.findingsTable.WithRows(m.selectedResult.findingsRows)
			}
			// TODO: add logic for compared result picker
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
	findingsSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#3ad33c")).Render(fmt.Sprintf("%d", m.selectedResult.summaryData.numFindingsSatisfied))
	findingsNotSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#e36750")).Render(fmt.Sprintf("%d", m.selectedResult.summaryData.numFindings-m.selectedResult.summaryData.numFindingsSatisfied))
	observationsSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#3ad33c")).Render(fmt.Sprintf("%d", m.selectedResult.summaryData.numObservationsSatisfied))
	observationsNotSatisfied := lipgloss.NewStyle().Foreground(lipgloss.Color("#e36750")).Render(fmt.Sprintf("%d", m.selectedResult.summaryData.numObservations-m.selectedResult.summaryData.numObservationsSatisfied))
	summaryText := fmt.Sprintf("Summary: %d (%s/%s) Findings - %d (%s/%s) Observations",
		m.selectedResult.summaryData.numFindings, findingsSatisfied, findingsNotSatisfied,
		m.selectedResult.summaryData.numObservations, observationsSatisfied, observationsNotSatisfied,
	)

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

func (m *Model) UpdateWithAssessmentResults(assessmentResults *oscalTypes_1_1_2.AssessmentResults) {
	var selectedResult result

	results := GetResults(assessmentResults)

	if len(results) != 0 {
		selectedResult = results[0]
	}

	// Set up tables
	findingsTableColumns := []table.Column{
		table.NewFlexColumn(columnKeyName, "Control", 1).WithFiltered(true),
		table.NewFlexColumn(columnKeyStatus, "Status", 1),
		table.NewFlexColumn(columnKeyDescription, "Description", 4),
	}

	observationsTableColumns := []table.Column{
		table.NewFlexColumn(columnKeyName, "Observation", 1).WithFiltered(true),
		table.NewFlexColumn(columnKeyStatus, "Status", 1),
		table.NewFlexColumn(columnKeyDescription, "Remarks", 4),
	}

	findingsTable := table.New(findingsTableColumns).
		WithRows(selectedResult.findingsRows).
		WithBaseStyle(common.TableStyleBase).
		Filtered(true).
		SortByAsc(columnKeyName)

	observationsTable := table.New(observationsTableColumns).
		WithRows(selectedResult.observationsRows).
		WithBaseStyle(common.TableStyleBase).
		Filtered(true).
		SortByAsc(columnKeyName)

	// Update model parameters
	resultItems := make([]string, len(results))
	for i, c := range results {
		resultItems[i] = getResultText(c)
	}

	m.results = results
	m.selectedResult = selectedResult
	m.resultsPicker.UpdateItems(resultItems)
	comparedResultItems := getComparedResults(results, selectedResult)
	m.comparedResultsPicker.UpdateItems(comparedResultItems)

	m.observationsTable = observationsTable
	m.findingsTable = findingsTable
}

// func (m *Model) UpdateWithComparedResults(result, comparedResult *oscalTypes_1_1_2.Result) {
// 	resultComparisonMap := pkgResult.NewResultComparisonMap(*result, *comparedResult)
// }

func (m *Model) updateSizing(height, width int) {
	m.height = height
	m.width = width
	totalHeight := m.height

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

func (m *Model) GetDimensions() (height, width int) {
	return m.height, m.width
}

func (m *Model) getObservationsByFinding(relatedObs []string) []table.Row {
	obsRows := make([]table.Row, 0)
	for _, uuid := range relatedObs {
		if obsRow, ok := m.selectedResult.observationsMap[uuid]; ok {
			obsRows = append(obsRows, obsRow)
		}
	}

	return obsRows
}
