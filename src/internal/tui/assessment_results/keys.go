package assessmentresults

import (
	"github.com/charmbracelet/bubbles/key"
	tui "github.com/defenseunicorns/lula/src/internal/tui/common"
)

type keys struct {
	Quit     key.Binding
	Evaluate key.Binding
	Confirm  key.Binding
	Navigate key.Binding
	Help     key.Binding
}

var assessmentHotkeys = keys{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Evaluate: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "evaluate"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
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

func (a keys) ShortHelp() []key.Binding {
	return []key.Binding{a.Evaluate, a.Help}
}

func (a keys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{a.Evaluate}, {a.Confirm}, {a.Navigate}, {a.Help}, {a.Quit},
	}
}

func (m *Model) updateKeyBindings() {
	m.findings.KeyMap = tui.UnfocusedListKeyMap()
	m.findings.SetDelegate(tui.NewUnfocusedDelegate())

	switch m.focus {
	case focusFindings:
		m.findings.KeyMap = tui.FocusedListKeyMap()
		m.findings.SetDelegate(tui.NewFocusedDelegate())
	}
}
