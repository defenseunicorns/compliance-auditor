package component

import (
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/types"
)

type Model struct {
	open               bool
	help               common.HelpModel
	keys               keys
	focus              focus
	componentModel     *oscalTypes_1_1_2.ComponentDefinition
	components         []component
	selectedComponent  component
	componentPicker    common.PickerModel
	frameworks         []framework
	selectedFramework  framework
	frameworkPicker    common.PickerModel
	controlPicker      viewport.Model
	controls           blist.Model
	selectedControl    control
	remarks            viewport.Model
	remarksEditor      textarea.Model
	description        viewport.Model
	descriptionEditor  textarea.Model
	validationPicker   viewport.Model
	validations        blist.Model
	selectedValidation validationLink
	width              int
	height             int
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
	oscalComponent    *oscalTypes_1_1_2.DefinedComponent
	uuid, title, desc string
	frameworks        []framework
}

type framework struct {
	oscalFramework *oscalTypes_1_1_2.ControlImplementationSet
	name, uuid     string
	controls       []control
}

type control struct {
	oscalControl *oscalTypes_1_1_2.ImplementedRequirementControlImplementation
	uuid, title  string
	validations  []validationLink
}

func (i control) Title() string       { return i.title }
func (i control) Description() string { return i.uuid }
func (i control) FilterValue() string { return i.title }

type validationLink struct {
	oscalLink  *oscalTypes_1_1_2.Link
	text       string
	validation *types.LulaValidation
}

func (i validationLink) Title() string       { return i.validation.Name }
func (i validationLink) Description() string { return i.text }
func (i validationLink) FilterValue() string { return i.validation.Name }

func (m *Model) Close() {
	m.open = false
}

func (m *Model) Open(height, width int) {
	m.open = true
	m.UpdateSizing(height, width)
}

// GetComponentDefinition returns the component definition model, used on save events
func (m *Model) GetComponentDefinition() *oscalTypes_1_1_2.ComponentDefinition {
	return m.componentModel
}

// TestSetSelectedControl is a test helper function to set the selected control
func (m *Model) TestSetSelectedControl(title string) {
	var idx int
	if len(m.selectedFramework.controls) > 0 {
		for i, c := range m.selectedFramework.controls {
			if c.title == title {
				m.selectedControl = c
				idx = i
				break
			}
		}
	}
	m.controls.Select(idx)
}

func (m *Model) UpdateRemarks(remarks string) {
	// TODO: handle any race conditions updating the control?
	if m.selectedControl.oscalControl != nil {
		m.selectedControl.oscalControl.Remarks = remarks
	}
}

func (m *Model) UpdateDescription(description string) {
	// TODO: handle any race conditions updating the control?
	if m.selectedControl.oscalControl != nil {
		m.selectedControl.oscalControl.Description = description
	}
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
	panelHeight := common.PanelTitleStyle.GetHeight() + common.PanelStyle.GetVerticalPadding() + common.PanelStyle.GetVerticalMargins() + 1 // 1 for border

	remarksOutsideHeight := bottomSectionHeight / 4
	remarksInsideHeight := remarksOutsideHeight - panelHeight

	descriptionOutsideHeight := bottomSectionHeight / 4
	descriptionInsideHeight := descriptionOutsideHeight - panelHeight
	validationsHeight := bottomSectionHeight - remarksOutsideHeight - descriptionOutsideHeight - 2*common.PanelTitleStyle.GetHeight()

	// Update widget sizing
	m.help.Width = m.width

	m.controls.SetHeight(m.height - common.PanelTitleStyle.GetHeight() - 1)
	m.controls.SetWidth(leftWidth - common.PanelStyle.GetHorizontalPadding())

	m.controlPicker.Height = bottomSectionHeight
	m.controlPicker.Width = leftWidth - common.PanelStyle.GetHorizontalPadding()

	m.remarks.Height = remarksInsideHeight
	m.remarks.Width = rightWidth
	m.remarks, _ = m.remarks.Update(tea.WindowSizeMsg{Width: rightWidth, Height: remarksInsideHeight}) // rebuild remarks for line wrapping?

	m.remarksEditor.SetHeight(m.remarks.Height - 1)
	m.remarksEditor.SetWidth(m.remarks.Width - 5) // probably need to fix this to be a func

	m.description.Height = descriptionInsideHeight
	m.description.Width = rightWidth
	m.description, _ = m.description.Update(tea.WindowSizeMsg{Width: rightWidth, Height: descriptionInsideHeight})

	m.descriptionEditor.SetHeight(m.description.Height - 1)
	m.descriptionEditor.SetWidth(m.description.Width - 5) // probably need to fix this to be a func

	m.validations.SetHeight(validationsHeight - common.PanelTitleStyle.GetHeight())
	m.validations.SetWidth(rightWidth - common.PanelStyle.GetHorizontalPadding())

	m.validationPicker.Height = validationsHeight
	m.validationPicker.Width = rightWidth
}

func (m *Model) GetDimensions() (height, width int) {
	return m.height, m.width
}

func (m *Model) updateKeyBindings() {
	m.outOfFocus()
	m.updateFocusHelpKeys()

	switch m.focus {

	case focusControls:
		m.controls.KeyMap = common.FocusedListKeyMap()
		m.controls.SetDelegate(common.NewFocusedDelegate())

	case focusValidations:
		m.validations.KeyMap = common.FocusedListKeyMap()
		m.validations.SetDelegate(common.NewFocusedDelegate())

	case focusRemarks:
		m.remarks.KeyMap = common.FocusedPanelKeyMap()
		m.remarks.MouseWheelEnabled = true
		if m.remarksEditor.Focused() {
			m.remarksEditor.KeyMap = common.FocusedTextAreaKeyMap()
			m.keys = componentEditKeys
		} else {
			m.remarksEditor.KeyMap = common.UnfocusedTextAreaKeyMap()
			m.keys = componentKeys
		}

	case focusDescription:
		m.description.KeyMap = common.FocusedPanelKeyMap()
		m.description.MouseWheelEnabled = true
		if m.descriptionEditor.Focused() {
			m.descriptionEditor.KeyMap = common.FocusedTextAreaKeyMap()
			m.keys = componentEditKeys
		} else {
			m.descriptionEditor.KeyMap = common.UnfocusedTextAreaKeyMap()
			m.keys = componentKeys
		}

	}
}

// func for outOfFocus to run just when focus switches between items
func (m *Model) outOfFocus() {
	focusMinusOne := m.focus - 1
	focusPlusOne := m.focus + 1

	if m.focus == 0 {
		focusMinusOne = maxFocus
	}
	if m.focus == maxFocus {
		focusPlusOne = 0
	}

	for _, f := range []focus{focusMinusOne, focusPlusOne} {
		// Turn off keys for out of focus items
		switch f {
		case focusControls:
			m.controls.KeyMap = common.UnfocusedListKeyMap()

		case focusValidations:
			m.validations.KeyMap = common.UnfocusedListKeyMap()
			m.validations.SetDelegate(common.NewUnfocusedDelegate())

		case focusRemarks:
			m.remarks.KeyMap = common.UnfocusedPanelKeyMap()
			m.remarks.MouseWheelEnabled = false

		case focusDescription:
			m.description.KeyMap = common.UnfocusedPanelKeyMap()
			m.description.MouseWheelEnabled = false
		}
	}
}

func (m *Model) updateFocusHelpKeys() {
	switch m.focus {
	case focusComponentSelection:
		m.help.ShortHelp = shortHelpDialogBox
		m.help.FullHelpOneLine = fullHelpDialogBoxOneLine
		m.help.FullHelp = fullHelpDialogBox
	case focusFrameworkSelection:
		m.help.ShortHelp = shortHelpDialogBox
		m.help.FullHelpOneLine = fullHelpDialogBoxOneLine
		m.help.FullHelp = fullHelpDialogBox
	case focusControls:
		m.help.ShortHelp = common.ShortHelpList
		m.help.FullHelpOneLine = common.FullHelpListOneLine
		m.help.FullHelp = common.FullHelpList
	case focusRemarks:
		if m.remarksEditor.Focused() {
			m.help.ShortHelp = common.ShortHelpEditing
			m.help.FullHelpOneLine = common.FullHelpEditingOneLine
			m.help.FullHelp = common.FullHelpEditing
		} else {
			m.help.ShortHelp = shortHelpEditableDialogBox
			m.help.FullHelpOneLine = fullHelpEditableDialogBoxOneLine
			m.help.FullHelp = fullHelpEditableDialogBox
		}
	case focusDescription:
		if m.descriptionEditor.Focused() {
			m.help.ShortHelp = common.ShortHelpEditing
			m.help.FullHelpOneLine = common.FullHelpEditingOneLine
			m.help.FullHelp = common.FullHelpEditing
		} else {
			m.help.ShortHelp = shortHelpEditableDialogBox
			m.help.FullHelpOneLine = fullHelpEditableDialogBoxOneLine
			m.help.FullHelp = fullHelpEditableDialogBox
		}
	case focusValidations:
		m.help.ShortHelp = common.ShortHelpList
		m.help.FullHelpOneLine = common.FullHelpListOneLine
		m.help.FullHelp = common.FullHelpList
	default:
		m.help.ShortHelp = shortHelpNoFocus
		m.help.FullHelpOneLine = fullHelpNoFocusOneLine
		m.help.FullHelp = fullHelpNoFocus
	}
}

// func (m *Model) setNoFocusHelpKeys() {
// 	m.help.ShortHelp = []key.Binding{
// 		componentKeys.Navigation, componentKeys.SwitchModels, componentKeys.Help,
// 	}
// 	m.help.FullHelpOneLine = []key.Binding{
// 		componentKeys.Save, componentKeys.Navigation, componentKeys.SwitchModels, componentKeys.Help, componentKeys.Quit,
// 	}
// 	// This is currently unused - TODO: help overlay?
// 	m.help.FullHelp = [][]key.Binding{
// 		{componentKeys.SwitchModels}, {componentKeys.Navigation}, {componentKeys.Help}, {componentKeys.Quit},
// 	}
// }

// func (m *Model) setDialogBoxHelpKeys() {
// 	m.help.ShortHelp = []key.Binding{
// 		componentKeys.Select, componentKeys.Navigation, componentKeys.SwitchModels, componentKeys.Help,
// 	}
// 	m.help.FullHelpOneLine = []key.Binding{
// 		componentKeys.Select, componentKeys.Save, componentKeys.Navigation, componentKeys.SwitchModels, componentKeys.Help, componentKeys.Quit,
// 	}
// 	// This is currently unused - TODO: help overlay?
// 	m.help.FullHelp = [][]key.Binding{
// 		{componentKeys.SwitchModels}, {componentKeys.Navigation}, {componentKeys.Help}, {componentKeys.Quit},
// 	}
// }

// func (m *Model) setEditableDialogBoxHelpKeys() {
// 	m.help.ShortHelp = []key.Binding{
// 		componentKeys.Edit, componentKeys.Save, componentKeys.Navigation, componentKeys.SwitchModels, componentKeys.Help,
// 	}
// 	m.help.FullHelpOneLine = []key.Binding{
// 		componentKeys.Edit, componentKeys.Save, componentKeys.Navigation, componentKeys.SwitchModels, componentKeys.Help, componentKeys.Quit,
// 	}
// 	// This is currently unused - TODO: help overlay?
// 	m.help.FullHelp = [][]key.Binding{
// 		{componentKeys.Edit}, {componentKeys.Navigation}, {componentKeys.Help}, {componentKeys.Quit},
// 	}
// }

// func (m *Model) setEditingDialogBoxHelpKeys() {
// 	m.help.ShortHelp = []key.Binding{
// 		componentKeys.Confirm, componentKeys.Newline, componentKeys.Cancel,
// 	}
// 	m.help.FullHelpOneLine = []key.Binding{
// 		componentKeys.Confirm, componentKeys.Newline, componentKeys.Cancel, componentKeys.Save, componentKeys.Help, componentKeys.Quit,
// 	}
// 	// This is currently unused - TODO: help overlay?
// 	m.help.FullHelp = [][]key.Binding{
// 		{componentKeys.Confirm}, {componentKeys.Newline}, {componentKeys.Cancel}, {componentKeys.Quit},
// 	}
// }

// func (m *Model) setListHelpKeys() {
// 	m.help.ShortHelp = []key.Binding{
// 		componentKeys.Select, componentKeys.Up, componentKeys.Down, common.CommonKeys.Filter, componentKeys.Help,
// 	}
// 	m.help.FullHelpOneLine = []key.Binding{
// 		componentKeys.Select, componentKeys.Up, componentKeys.Down, common.CommonKeys.Filter, componentKeys.Cancel, componentKeys.Help, componentKeys.Quit,
// 	}
// 	// This is currently unused - TODO: help overlay?
// 	m.help.FullHelp = [][]key.Binding{
// 		{componentKeys.Edit}, {componentKeys.Navigation}, {componentKeys.Help}, {componentKeys.Quit},
// 	}
// }
