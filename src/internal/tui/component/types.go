package component

import (
	"github.com/charmbracelet/bubbles/key"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/types"
)

type Model struct {
	open  bool
	help  common.HelpModel
	keys  keys
	focus focus
	// focusLock              bool
	componentModel *oscalTypes_1_1_2.ComponentDefinition
	// inComponentOverlay     bool
	components        []component
	selectedComponent component
	// selectedComponentIndex int
	componentPicker common.PickerModel
	// componentPicker        viewport.Model
	inFrameworkOverlay bool
	frameworks         []framework
	selectedFramework  framework
	// selectedFrameworkIndex int
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

// UpdateComponentDefinition updates and returns the component definition model, used on save events
// func (m *Model) UpdateComponentDefinition() *oscalTypes_1_1_2.ComponentDefinition {
// 	// Add all edits to the component definition
// 	oscalModel := &oscalTypes_1_1_2.OscalCompleteSchema{
// 		ComponentDefinition: m.componentModel,
// 	}
// 	var err error
// 	if !m.editor.IsEmpty() {
// 		// run inject for each edit...
// 		for path, edits := range m.editor.EditsByPath {
// 			// TODO: probably handle different edit types here?
// 			// Also, how to parse all edits? right now it's just updates so just taking the last but if I do add and delete ones...
// 			edit := edits[len(edits)-1]
// 			oscalModel, err = oscal.InjectIntoOSCALModel(oscalModel, edit.Value, path)
// 			if err != nil {
// 				common.PrintToLog("error injecting edits: %v", err)
// 				common.PrintToLog("non-injected edit at %s: %v", path, edit.Value)
// 			}
// 		}

// 		m.editor.ResetEditor()
// 	}
// 	m.componentModel = oscalModel.ComponentDefinition
// 	return m.componentModel
// }

// GetComponentDefinition returns the component definition model, used on save events
func (m *Model) GetComponentDefinition() *oscalTypes_1_1_2.ComponentDefinition {
	return m.componentModel
}

func (m *Model) TestSetControl() {
	controlItems := make([]blist.Item, len(m.selectedFramework.controls))
	if len(m.selectedFramework.controls) > 0 {
		for i, c := range m.selectedFramework.controls {
			controlItems[i] = c
		}
	}
	m.controls.SetItems(controlItems)

	if len(m.controls.Items()) != 0 {
		m.selectedControl = m.controls.Items()[0].(control)
	}
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

	remarksOutsideHeight := bottomSectionHeight / 4
	remarksInsideHeight := remarksOutsideHeight - common.PanelTitleStyle.GetHeight()

	descriptionOutsideHeight := bottomSectionHeight / 4
	descriptionInsideHeight := descriptionOutsideHeight - common.PanelTitleStyle.GetHeight()
	validationsHeight := bottomSectionHeight - remarksOutsideHeight - descriptionOutsideHeight - 2*common.PanelTitleStyle.GetHeight()

	// Update widget sizing
	m.help.Width = m.width

	m.controls.SetHeight(m.height - common.PanelTitleStyle.GetHeight() - 1)
	m.controls.SetWidth(leftWidth - common.PanelStyle.GetHorizontalPadding())

	m.controlPicker.Height = bottomSectionHeight
	m.controlPicker.Width = leftWidth - common.PanelStyle.GetHorizontalPadding()

	m.remarks.Height = remarksInsideHeight - 1
	m.remarks.Width = rightWidth
	m.remarks, _ = m.remarks.Update(tea.WindowSizeMsg{Width: rightWidth, Height: remarksInsideHeight - 1})

	m.remarksEditor.SetHeight(m.remarks.Height - 1)
	m.remarksEditor.SetWidth(m.remarks.Width - 5) // probably need to fix this to be a func

	m.description.Height = descriptionInsideHeight - 1
	m.description.Width = rightWidth
	m.description, _ = m.description.Update(tea.WindowSizeMsg{Width: rightWidth, Height: descriptionInsideHeight - 1})

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
	// m.controls.KeyMap = common.UnfocusedListKeyMap()
	// // m.controls.SetDelegate(common.NewUnfocusedDelegate())
	// m.validations.KeyMap = common.UnfocusedListKeyMap()
	// m.validations.SetDelegate(common.NewUnfocusedDelegate())

	// m.remarks.KeyMap = common.UnfocusedPanelKeyMap()
	// m.description.KeyMap = common.UnfocusedPanelKeyMap()
	m.outOfFocus()

	switch m.focus {
	case focusComponentSelection:
		m.setDialogBoxHelpKeys()

	case focusFrameworkSelection:
		m.setDialogBoxHelpKeys()

	case focusControls:
		m.setListHelpKeys()
		m.controls.KeyMap = common.FocusedListKeyMap()
		m.controls.SetDelegate(common.NewFocusedDelegate())

	case focusValidations:
		m.setListHelpKeys()
		m.validations.KeyMap = common.FocusedListKeyMap()
		m.validations.SetDelegate(common.NewFocusedDelegate())

	case focusRemarks:
		m.remarks.KeyMap = common.FocusedPanelKeyMap()
		if m.remarksEditor.Focused() {
			m.setEditingDialogBoxHelpKeys()
			m.remarksEditor.KeyMap = common.FocusedTextAreaKeyMap()
			m.keys = componentEditKeys
		} else {
			m.setEditableDialogBoxHelpKeys()
			m.remarksEditor.KeyMap = common.UnfocusedTextAreaKeyMap()
			m.keys = componentKeys
		}

	case focusDescription:
		m.description.KeyMap = common.FocusedPanelKeyMap()
		if m.descriptionEditor.Focused() {
			m.setEditingDialogBoxHelpKeys()
			m.descriptionEditor.KeyMap = common.FocusedTextAreaKeyMap()
			m.keys = componentEditKeys
		} else {
			m.setEditableDialogBoxHelpKeys()
			m.descriptionEditor.KeyMap = common.UnfocusedTextAreaKeyMap()
			m.keys = componentKeys
		}

	default:
		m.setNoFocusHelpKeys()
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

		case focusDescription:
			m.description.KeyMap = common.UnfocusedPanelKeyMap()
		}
	}
}

func (m *Model) setNoFocusHelpKeys() {
	m.help.ShortHelp = []key.Binding{
		componentKeys.Navigation, componentKeys.Help,
	}
	m.help.FullHelpOneLine = []key.Binding{
		componentKeys.Navigation, componentKeys.Help, componentKeys.Quit,
	}
	m.help.FullHelp = [][]key.Binding{
		{componentKeys.Navigation}, {componentKeys.Help}, {componentKeys.Quit},
	}
}

func (m *Model) setDialogBoxHelpKeys() {
	m.help.ShortHelp = []key.Binding{
		componentKeys.Select, componentKeys.Help,
	}
	m.help.FullHelpOneLine = []key.Binding{
		componentKeys.Select, componentKeys.Navigation, componentKeys.Help, componentKeys.Quit,
	}
	m.help.FullHelp = [][]key.Binding{
		{componentKeys.Select}, {componentKeys.Navigation}, {componentKeys.Help}, {componentKeys.Quit},
	}
}

func (m *Model) setEditableDialogBoxHelpKeys() {
	m.help.ShortHelp = []key.Binding{
		componentKeys.Edit, componentKeys.Help,
	}
	m.help.FullHelpOneLine = []key.Binding{
		componentKeys.Edit, componentKeys.Navigation, componentKeys.Help, componentKeys.Quit,
	}
	m.help.FullHelp = [][]key.Binding{
		{componentKeys.Edit}, {componentKeys.Navigation}, {componentKeys.Help}, {componentKeys.Quit},
	}
}

func (m *Model) setEditingDialogBoxHelpKeys() {
	m.help.ShortHelp = []key.Binding{
		componentKeys.Confirm, componentKeys.Newline, componentKeys.Cancel, componentKeys.Help,
	}
	m.help.FullHelpOneLine = []key.Binding{
		componentKeys.Confirm, componentKeys.Newline, componentKeys.Cancel, componentKeys.Save, componentKeys.Help, componentKeys.Quit,
	}
	m.help.FullHelp = [][]key.Binding{
		{componentKeys.Confirm}, {componentKeys.Newline}, {componentKeys.Cancel}, {componentKeys.Quit},
	}
}

func (m *Model) setListHelpKeys() {
	m.help.ShortHelp = []key.Binding{
		componentKeys.Up, componentKeys.Down, common.CommonKeys.Filter, componentKeys.Help,
	}
	m.help.FullHelpOneLine = []key.Binding{
		componentKeys.Up, componentKeys.Down, common.CommonKeys.Filter, componentKeys.Cancel, componentKeys.Help, componentKeys.Quit,
	}
	m.help.FullHelp = [][]key.Binding{
		{componentKeys.Edit}, {componentKeys.Navigation}, {componentKeys.Help}, {componentKeys.Quit},
	}
}
