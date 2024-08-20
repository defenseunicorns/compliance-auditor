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
	Quit: common.CommonHotkeys.Quit,
	Help: common.CommonHotkeys.Help,
	Validate: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "validate"),
	),
	Evaluate: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "evaluate"),
	),
	Confirm: common.PickerHotkeys.Confirm,
	Cancel:  common.PickerHotkeys.Cancel,
	Navigation: key.NewBinding(
		key.WithKeys("left", "h", "right", "l"),
		key.WithHelp("←/h, →/l", "navigation"),
	),
	NavigateLeft: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "navigate left"),
	),
	NavigateRight: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "navigate right"),
	),
	SwitchModels: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
		key.WithHelp("tab/shift+tab", "switch models"),
	),
	Up:   common.PickerHotkeys.Up,
	Down: common.PickerHotkeys.Down,
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
