package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PickerKind string

const (
	pickerWidth  = 80
	pickerHeight = 20
)

type PickerOpenMsg struct {
	Kind PickerKind
}
type PickerItemSelected struct {
	Selected int
	From     PickerKind
}

type PickerModel struct {
	Open     bool
	items    []string
	selected int
	title    string
	kind     PickerKind
	help     HelpModel
	viewer   viewport.Model
}

func NewPickerModel(title string, kind PickerKind, items []string, initSelected int) PickerModel {
	help := NewHelpModel(true)
	help.ShortHelp = ShortHelpPicker
	return PickerModel{
		items:    items,
		selected: initSelected,
		title:    title,
		kind:     kind,
		help:     help,
		viewer:   viewport.New(pickerWidth, pickerHeight),
	}
}

func (m PickerModel) Init() tea.Cmd {
	return nil
}

func (m PickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if m.Open {
			switch k {
			case ContainsKey(k, PickerKeys.Up.Keys()):
				m.selected--
				if m.selected < 0 {
					m.selected = len(m.items) - 1
				}

			case ContainsKey(k, PickerKeys.Down.Keys()):
				m.selected++
				if m.selected >= len(m.items) {
					m.selected = 0
				}

			case ContainsKey(k, PickerKeys.Select.Keys()):
				m.Open = false
				return m, func() tea.Msg {
					return PickerItemSelected{
						Selected: m.selected,
						From:     m.kind,
					}
				}

			case ContainsKey(k, PickerKeys.Cancel.Keys()):
				m.Open = false
			}
		}
	case PickerOpenMsg:
		if msg.Kind == m.kind {
			m.Open = true
		}
	}
	return m, nil
}

func (m PickerModel) View() string {
	itemHeight := len(m.items) + 2
	if itemHeight > pickerHeight {
		itemHeight = pickerHeight
	}

	overlayPickerStyle := OverlayStyle.
		Width(pickerWidth).
		Height(itemHeight)

	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%s\n\n", m.title))

	for i, itm := range m.items {
		if m.selected == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(itm)
		s.WriteString("\n")
	}
	// m.viewer.SetContent(s.String())
	pickerContent := lipgloss.JoinVertical(lipgloss.Top, s.String(), m.help.View())
	return overlayPickerStyle.Render(pickerContent)
}
