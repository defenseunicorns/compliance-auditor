package component

import (
	"github.com/charmbracelet/bubbles/key"
	tui "github.com/defenseunicorns/lula/src/internal/tui/common"
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
		{k.Generate}, {k.Edit}, {k.Confirm}, {k.Navigate}, {k.Help}, {k.Quit},
	}
}

func (m *Model) updateKeyBindings() {
	m.controls.KeyMap = tui.UnfocusedListKeyMap()
	m.validations.KeyMap = tui.UnfocusedListKeyMap()
	m.remarks.KeyMap = tui.UnfocusedPanelKeyMap()
	m.controls.SetDelegate(tui.NewUnfocusedDelegate())
	m.validations.SetDelegate(tui.NewUnfocusedDelegate())

	switch m.focus {
	case focusValidations:
		m.validations.KeyMap = tui.FocusedListKeyMap()
		m.validations.SetDelegate(tui.NewFocusedDelegate())
	case focusControls:
		m.controls.KeyMap = tui.FocusedListKeyMap()
		m.controls.SetDelegate(tui.NewFocusedDelegate())
	case focusRemarks:
		m.remarks.KeyMap = tui.FocusedPanelKeyMap()
	}
}
