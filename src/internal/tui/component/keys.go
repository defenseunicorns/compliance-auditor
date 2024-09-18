package component

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

type keys struct {
	Edit          key.Binding
	Generate      key.Binding
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

var componentKeys = keys{
	Quit: common.CommonHotkeys.Quit,
	Help: common.CommonHotkeys.Help,
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
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
	return []key.Binding{k.Navigation, k.Help}
}

func (k keys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Confirm}, {k.Navigation}, {k.SwitchModels}, {k.Help}, {k.Quit},
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
