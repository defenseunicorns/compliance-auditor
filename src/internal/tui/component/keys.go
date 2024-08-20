package component

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

type keys struct {
	Edit     key.Binding
	Generate key.Binding
	Confirm  key.Binding
	Navigate key.Binding
	Help     key.Binding
	Quit     key.Binding
}

var componentKeys = keys{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Generate: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "generate"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Navigate: key.NewBinding(
		key.WithKeys("left", "right", "h", "l"),
		key.WithHelp("←/h/→/l", "navigate"),
	),
}

func (k keys) ShortHelp() []key.Binding {
	return []key.Binding{k.Generate, k.Help}
}

func (k keys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Generate}, {k.Confirm}, {k.Navigate}, {k.Help}, {k.Quit},
	}
}

func (m *Model) updateKeyBindings() {
	m.controls.KeyMap = common.UnfocusedListKeyMap()
	// m.controls.SetDelegate(common.NewUnfocusedDelegate())
	m.validations.KeyMap = common.UnfocusedListKeyMap()
	m.validations.SetDelegate(common.NewUnfocusedDelegate())

	m.remarks.KeyMap = common.UnfocusedPanelKeyMap()
	m.description.KeyMap = common.UnfocusedPanelKeyMap()

	switch m.focus {
	case focusComponentSelection:
	case focusValidations:
		m.validations.KeyMap = common.FocusedListKeyMap()
		m.validations.SetDelegate(common.NewFocusedDelegate())
	case focusControls:
		m.controls.KeyMap = common.FocusedListKeyMap()
		m.controls.SetDelegate(common.NewFocusedDelegate())
	case focusRemarks:
		m.remarks.KeyMap = common.FocusedPanelKeyMap()
	case focusDescription:
		m.description.KeyMap = common.FocusedPanelKeyMap()
	}
}
