package component

import (
	"github.com/charmbracelet/bubbles/help"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/types"
)

type Model struct {
	open                   bool
	help                   help.Model
	keys                   keys
	focus                  focus
	inComponentOverlay     bool
	components             []component
	selectedComponent      component
	selectedComponentIndex int
	componentPicker        viewport.Model
	inFrameworkOverlay     bool
	frameworks             []framework
	selectedFramework      framework
	selectedFrameworkIndex int
	frameworkPicker        viewport.Model
	controlPicker          viewport.Model
	controls               blist.Model
	selectedControl        control
	remarks                viewport.Model
	description            viewport.Model
	validationPicker       viewport.Model
	validations            blist.Model
	selectedValidation     validationLink
	width                  int
	height                 int
}

type focus int

const (
	noComponentFocus focus = iota
	focusComponentSelection
	focusFrameworkSelection
	focusControls
	focusRemarks
	focusDescription
	focusValidations
)

var maxFocus = focusValidations

type component struct {
	uuid, title, desc string
	frameworks        []framework
}

type framework struct {
	name     string
	controls []control
}

type validationLink struct {
	text       string
	validation *types.LulaValidation
}

func (i validationLink) Title() string       { return i.validation.Name }
func (i validationLink) Description() string { return i.text }
func (i validationLink) FilterValue() string { return i.validation.Name }

type control struct {
	uuid, remarks, title, desc string
	validations                []validationLink
}

func (i control) Title() string       { return i.title }
func (i control) Description() string { return i.uuid }
func (i control) FilterValue() string { return i.title }

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

	// Set internal sizing properties
	totalHeight := m.height
	leftWidth := m.width / 4
	rightWidth := m.width - leftWidth - common.PanelStyle.GetHorizontalPadding() - common.PanelStyle.GetHorizontalMargins()

	topSectionHeight := common.HelpStyle(m.width).GetHeight() + common.DialogBoxStyle.GetHeight()
	bottomSectionHeight := totalHeight - topSectionHeight

	remarksOutsideHeight := bottomSectionHeight / 4
	remarksInsideHeight := remarksOutsideHeight - common.PanelTitleStyle.GetHeight()

	descriptionOutsideHeight := bottomSectionHeight / 4
	descriptionInsideHeight := descriptionOutsideHeight - common.PanelTitleStyle.GetHeight()
	validationsHeight := bottomSectionHeight - remarksOutsideHeight - descriptionOutsideHeight - 2*common.PanelTitleStyle.GetHeight()

	// Update widget sizing
	m.controls.SetHeight(m.height - common.PanelTitleStyle.GetHeight() - 1)
	m.controls.SetWidth(leftWidth - common.PanelStyle.GetHorizontalPadding())

	m.controlPicker.Height = bottomSectionHeight
	m.controlPicker.Width = leftWidth - common.PanelStyle.GetHorizontalPadding()

	m.remarks.Height = remarksInsideHeight - 1
	m.remarks.Width = rightWidth
	m.remarks, _ = m.remarks.Update(tea.WindowSizeMsg{Width: rightWidth, Height: remarksInsideHeight - 1})

	m.description.Height = descriptionInsideHeight - 1
	m.description.Width = rightWidth
	m.description, _ = m.description.Update(tea.WindowSizeMsg{Width: rightWidth, Height: descriptionInsideHeight - 1})

	m.validations.SetHeight(validationsHeight - common.PanelTitleStyle.GetHeight())
	m.validations.SetWidth(rightWidth - common.PanelStyle.GetHorizontalPadding())

	m.validationPicker.Height = validationsHeight
	m.validationPicker.Width = rightWidth
}

func (m *Model) GetDimensions() (height, width int) {
	return m.height, m.width
}
