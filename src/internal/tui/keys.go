package tui

import (
	"github.com/charmbracelet/bubbles/key"
	blist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type keys struct {
	Quit         []string
	Confirm      []string
	Left         []string
	Right        []string
	Up           []string
	Down         []string
	TabLeft      []string
	TabRight     []string
	CancelTyping []string
}

type assessmentKeys struct {
	Quit     key.Binding
	Evaluate key.Binding
	Confirm  key.Binding
	Navigate key.Binding
	Help     key.Binding
}

type componentKeys struct {
	Quit     key.Binding
	Generate key.Binding
	Confirm  key.Binding
	Navigate key.Binding
	Help     key.Binding
}

type listKeys struct {
	Up      key.Binding
	Down    key.Binding
	Slash   key.Binding
	Confirm key.Binding
	Escape  key.Binding
	Help    key.Binding
}

var hotkeys = keys{
	Quit:         []string{"q", "ctrl+c"},
	Confirm:      []string{"enter"},
	Left:         []string{"left", "h"},
	Right:        []string{"right", "l"},
	Up:           []string{"up", "k"},
	Down:         []string{"down", "j"},
	TabLeft:      []string{"["},
	TabRight:     []string{"]"},
	CancelTyping: []string{"esc"},
}

var assessmentHotkeys = assessmentKeys{
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

var componentHotkeys = componentKeys{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Generate: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "generate"),
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

func (a assessmentKeys) ShortHelp() []key.Binding {
	return []key.Binding{a.Evaluate, a.Help}
}

func (a assessmentKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{a.Evaluate}, {a.Confirm}, {a.Navigate}, {a.Help}, {a.Quit},
	}
}

var listHotkeys = listKeys{
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

func containsKey(v string, a []string) string {
	for _, i := range a {
		if i == v {
			return v
		}
	}
	return ""
}

func (m *model) handleKey(key string, cmd tea.Cmd) tea.Cmd {
	switch key {
	case containsKey(key, hotkeys.TabLeft):
		if m.cursor > 0 {
			m.cursor--
		}

	case containsKey(key, hotkeys.TabRight):
		if m.cursor < len(m.tabs)-1 {
			m.cursor++
		}

		// case containsKey(key, hotkeys.Confirm):
		// 	if m.componentModel.open {
		// 		selectedItem := m.componentModel.controls.SelectedItem().(item)
		// 		m.componentModel.content =
		// 	}
	}

	return cmd
}

// type keyMap struct {
// 	blist.DefaultKeyMap
// 	UnbindLeft  tea.Key
// 	UnbindRight tea.Key
// }

func focusedListKeyMap() blist.KeyMap {
	km := blist.DefaultKeyMap()
	km.NextPage.Unbind()
	km.PrevPage.Unbind()
	// do other bindings/unbindings...

	return km
}

func unfocusedListKeyMap() blist.KeyMap {
	km := blist.KeyMap{}

	return km
}

// func newListKeyMap() blist.KeyMap {
// 	km := blist.DefaultKeyMap()
// 	km.NextPage.Unbind()
// 	km.PrevPage.Unbind()

// 	return km
// }

// func (m *model) confirmToQuit(msg string) bool {
// 	switch msg {
// 	case containsKey(msg, hotkeys.Quit), containsKey(msg, hotkeys.CancelTyping):
// 		m.cancelWarnModal()
// 		m.confirmToQuit = false
// 		return false
// 	case containsKey(msg, hotkeys.Confirm):
// 		return true
// 	default:
// 		return false
// 	}
// }
