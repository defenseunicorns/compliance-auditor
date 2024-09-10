package assessmentresults

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

type keys struct {
	Validate      key.Binding
	Evaluate      key.Binding
	Confirm       key.Binding
	Cancel        key.Binding
	Navigation    key.Binding
	NavigateLeft  key.Binding
	NavigateRight key.Binding
	SwitchModels  key.Binding
	Up            key.Binding
	Down          key.Binding
	Help          key.Binding
	Quit          key.Binding
}

var assessmentHotkeys = keys{
	Quit: common.CommonKeys.Quit,
	Help: common.CommonKeys.Help,
	Validate: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "validate"),
	),
	Evaluate: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "evaluate"),
	),
	Confirm:       common.CommonKeys.Confirm,
	Cancel:        common.CommonKeys.Cancel,
	Navigation:    common.CommonKeys.Navigation,
	NavigateLeft:  common.CommonKeys.NavigateLeft,
	NavigateRight: common.CommonKeys.NavigateRight,
	SwitchModels:  common.CommonKeys.NavigateModels,
	Up:            common.PickerKeys.Up,
	Down:          common.PickerKeys.Down,
}

func (k keys) ShortHelp() []key.Binding {
	return []key.Binding{k.Validate, k.Evaluate, k.Help}
}

func (k keys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Evaluate}, {k.Confirm}, {k.Navigation}, {k.SwitchModels}, {k.Help}, {k.Quit},
	}
}

func (m *Model) updateKeyBindings() {
	m.findings.KeyMap = common.UnfocusedListKeyMap()
	m.findings.SetDelegate(common.NewUnfocusedDelegate())

	switch m.focus {
	case focusFindings:
		m.findings.KeyMap = common.FocusedListKeyMap()
		m.findings.SetDelegate(common.NewFocusedDelegate())
	}
}
