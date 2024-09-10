package assessmentresults

import (
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

type Model struct {
	open                bool
	help                common.HelpModel
	keys                keys
	focus               focus
	inResultOverlay     bool
	results             []result
	resultsPicker       viewport.Model
	selectedResult      result
	selectedResultIndex int
	compareResult       result
	compareResultIndex  int
	findings            blist.Model
	findingPicker       viewport.Model
	findingSummary      viewport.Model
	observationSummary  viewport.Model
	width               int
	height              int
}

type focus int

const (
	noFocus focus = iota
	focusResultSelection
	focusCompareSelection
	focusFindings
	focusSummary
	focusObservations
)

var maxFocus = focusObservations

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

func (m *Model) Close() {
	m.open = false
}

func (m *Model) Open(height, width int) {
	m.open = true
	m.UpdateSizing(height, width)
}

func (m *Model) UpdateSizing(height, width int) {
	m.height = height
	m.width = width

	totalHeight := m.height
	leftWidth := m.width / 4
	rightWidth := m.width - leftWidth - common.PanelStyle.GetHorizontalPadding() - common.PanelStyle.GetHorizontalMargins()

	topSectionHeight := common.HelpStyle(m.width).GetHeight() + common.DialogBoxStyle.GetHeight()
	bottomSectionHeight := totalHeight - topSectionHeight
	bottomRightPanelHeight := (bottomSectionHeight - 2*common.PanelTitleStyle.GetHeight() - 2*common.PanelTitleStyle.GetVerticalMargins()) / 2

	m.findings.SetHeight(totalHeight - topSectionHeight - common.PanelTitleStyle.GetHeight() - common.PanelStyle.GetVerticalPadding())
	m.findings.SetWidth(leftWidth - common.PanelStyle.GetHorizontalPadding())

	m.findingPicker.Height = bottomSectionHeight
	m.findingPicker.Width = leftWidth - common.PanelStyle.GetHorizontalPadding()

	m.findingSummary.Height = bottomRightPanelHeight
	m.findingSummary.Width = rightWidth

	m.observationSummary.Height = bottomRightPanelHeight
	m.observationSummary.Width = rightWidth
}

func (m *Model) GetDimensions() (height, width int) {
	return m.height, m.width
}
