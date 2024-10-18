package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	popupWidth  = 40
	popupHeight = 10
)

type PopupMsg struct {
	Title   string
	Content string
	Warning string
}
type PopupClose struct{}

type PopupFailMsg struct {
	Err error
}

type PopupModel struct {
	Open    bool
	Title   string
	Content string
	Warning string
	Help    HelpModel
}

func NewPopupModel(title, content string, helpKeys []key.Binding) PopupModel {
	help := NewHelpModel(true)
	help.ShortHelp = helpKeys
	return PopupModel{
		Help:    help,
		Title:   title,
		Content: content,
	}
}

func (m *PopupModel) UpdateText(title, content, warning string) {
	m.Title = title
	m.Content = content
	m.Warning = warning
}

func (m PopupModel) Init() tea.Cmd {
	return nil
}

func (m PopupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	DumpToLog(msg)
	switch msg := msg.(type) {
	case PopupMsg:
		m.UpdateText(msg.Title, msg.Content, msg.Warning)

	case PopupClose:
		m.Open = false
	}
	return m, nil
}

func (m PopupModel) View() string {
	popupStyle := OverlayWarnStyle.
		Width(popupWidth).
		Height(popupHeight)

	content := strings.Builder{}
	content.WriteString(fmt.Sprintf("%s\n", m.Title))
	if m.Content != "" {
		content.WriteString(fmt.Sprintf("\n%s\n", m.Content))
	}
	if m.Warning != "" {
		content.WriteString(fmt.Sprintf("\n⚠️ %s ⚠️\n", m.Warning))
	}

	popupContent := lipgloss.JoinVertical(lipgloss.Top, content.String(), m.Help.View())
	return popupStyle.Render(popupContent)
}
