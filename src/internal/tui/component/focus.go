package component

import (
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

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
			m.validations.ResetSelected()

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
		m.help.ShortHelp = shortHelpValidations
		m.help.FullHelpOneLine = fullHelpValidationsOneLine
		m.help.FullHelp = fullHelpValidations
	default:
		m.help.ShortHelp = shortHelpNoFocus
		m.help.FullHelpOneLine = fullHelpNoFocusOneLine
		m.help.FullHelp = fullHelpNoFocus
	}
}
