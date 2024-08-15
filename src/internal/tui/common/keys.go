package common

import (
	"github.com/charmbracelet/bubbles/key"
)

type keys struct {
	Quit     key.Binding
	Confirm  key.Binding
	TabLeft  key.Binding
	TabRight key.Binding
	Help     key.Binding
}

var CommonHotkeys = keys{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	TabRight: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "tab right"),
	),
	TabLeft: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "tab left"),
	),
}

func ContainsKey(v string, a []string) string {
	for _, i := range a {
		if i == v {
			return v
		}
	}
	return ""
}

// func (m *model) handleKey(key string, cmd tea.Cmd) tea.Cmd {
// 	switch key {
// 	case containsKey(key, hotkeys.TabLeft):
// 		if m.cursor > 0 {
// 			m.cursor--
// 		}

// 	case containsKey(key, hotkeys.TabRight):
// 		if m.cursor < len(m.tabs)-1 {
// 			m.cursor++
// 		}

// 		case containsKey(key, hotkeys.Confirm):
// 			if m.componentModel.open {
// 				selectedItem := m.componentModel.controls.SelectedItem().(item)
// 				m.componentModel.content =
// 			}
// 	}

// 	return cmd
// }

type listKeys struct {
	Up      key.Binding
	Down    key.Binding
	Slash   key.Binding
	Confirm key.Binding
	Escape  key.Binding
	Help    key.Binding
}

var ListHotkeys = listKeys{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "move down"),
	),
	Slash: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func (k listKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Help}
}

func (k listKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down}, {k.Slash, k.Confirm}, {k.Escape, k.Help},
	}
}
