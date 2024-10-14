package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DecisionOwner string

type DecisionOpenFromMsg struct {
	Title, Content, Warning string
	HelpKeys                []key.Binding
	From                    DecisionOwner
}
type DecisionYesMsg struct {
	To DecisionOwner
}
type DecisionNoMsg struct {
	To DecisionOwner
}

type DecisionModel struct {
	Open    bool
	title   string
	content string
	warning string
	help    HelpModel
	owner   DecisionOwner
	height  int
	width   int
}

func NewDefaultDecisionModel() DecisionModel {
	help := NewHelpModel(true)
	help.ShortHelp = []key.Binding{CommonKeys.Confirm, CommonKeys.Cancel}
	return DecisionModel{
		help:   help,
		height: defaultPopupHeight,
		width:  defaultPopupWidth,
	}
}

func (m *DecisionModel) UpdateHelp(helpKeys []key.Binding) {
	m.help.ShortHelp = helpKeys
}

func (m *DecisionModel) UpdateText(title, content, warning string) {
	m.title = title
	m.content = content
	m.warning = warning
}

func (m *DecisionModel) SetDimensions(height, width int) {
	m.height = height
	m.width = width
}

func (m DecisionModel) Init() tea.Cmd {
	return nil
}

func (m DecisionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if m.Open {
			switch k {
			case ContainsKey(k, CommonKeys.Confirm.Keys()):
				m.Open = false
				return m, func() tea.Msg {
					return DecisionYesMsg{
						To: m.owner,
					}
				}
			case ContainsKey(k, CommonKeys.Cancel.Keys()):
				m.Open = false
				return m, func() tea.Msg {
					return DecisionNoMsg{
						To: m.owner,
					}
				}
			}
		}

	case DecisionOpenFromMsg:
		m.Open = true
		m.owner = msg.From
		// Update any field passed in the message, otherwise keep as they were
		if msg.Title != "" {
			m.title = msg.Title
		}
		if msg.Content != "" {
			m.content = msg.Content
		}
		if msg.Warning != "" {
			m.warning = msg.Warning
		}
		if len(msg.HelpKeys) > 0 {
			m.help.ShortHelp = msg.HelpKeys
		}
	}

	return m, nil
}

func (m DecisionModel) View() string {
	popupStyle := OverlayWarnStyle.
		Width(m.width).
		Height(m.height)

	content := strings.Builder{}
	content.WriteString(fmt.Sprintf("%s\n", m.title))
	if m.content != "" {
		content.WriteString(fmt.Sprintf("\n%s\n", m.content))
	}
	if m.warning != "" {
		content.WriteString(fmt.Sprintf("\n⚠️ %s ⚠️\n", m.warning))
	}

	popupContent := lipgloss.JoinVertical(lipgloss.Top, content.String(), m.help.View())
	return popupStyle.Render(popupContent)
}
