package common

import (
	"github.com/charmbracelet/bubbles/key"
)

type Keys struct {
	Quit       key.Binding
	Confirm    key.Binding
	ModelLeft  key.Binding
	ModelRight key.Binding
	Help       key.Binding
}

var CommonHotkeys = Keys{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	ModelRight: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "model right"),
	),
	ModelLeft: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "model left"),
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
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
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

type pickerKeys struct {
	Up      key.Binding
	Down    key.Binding
	Confirm key.Binding
	Cancel  key.Binding
}

var PickerHotkeys = pickerKeys{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↳", "select"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("esc/q", "cancel"),
	),
}

func (k pickerKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Confirm, k.Cancel}
}

func (k pickerKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up}, {k.Down}, {k.Confirm}, {k.Cancel},
	}
}
