package component

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

type keys struct {
	Edit          key.Binding
	Generate      key.Binding
	Confirm       key.Binding
	Select        key.Binding
	Save          key.Binding
	Cancel        key.Binding
	Newline       key.Binding
	Navigation    key.Binding
	NavigateLeft  key.Binding
	NavigateRight key.Binding
	SwitchModels  key.Binding
	Up            key.Binding
	Down          key.Binding
	Help          key.Binding
	Quit          key.Binding
}

func (k keys) ShortHelp() []key.Binding {
	return []key.Binding{k.Navigation, k.Help}
}

func (k keys) SingleLineFullHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Navigation, k.SwitchModels, k.Help, k.Quit}
}

func (k keys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Confirm}, {k.Navigation}, {k.SwitchModels}, {k.Help}, {k.Quit},
	}
}

var componentKeys = keys{
	Quit:          common.CommonKeys.Quit,
	Help:          common.CommonKeys.Help,
	Edit:          common.CommonKeys.Edit,
	Save:          common.CommonKeys.Save,
	Select:        common.CommonKeys.Select,
	Confirm:       common.CommonKeys.Confirm,
	Cancel:        common.CommonKeys.Cancel,
	Newline:       common.CommonKeys.Newline,
	Navigation:    common.CommonKeys.Navigation,
	NavigateLeft:  common.CommonKeys.NavigateLeft,
	NavigateRight: common.CommonKeys.NavigateRight,
	SwitchModels:  common.CommonKeys.NavigateModels,
	Up:            common.CommonKeys.Up,
	Down:          common.CommonKeys.Down,
}

var componentEditKeys = keys{
	Save:    common.EditHotkeys.Save,
	Confirm: common.PickerHotkeys.Confirm,
	Cancel:  common.PickerHotkeys.Cancel,
}
