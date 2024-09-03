package common

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TbdModal struct {
	open    bool
	content string
	width   int
	height  int
}

func NewTbdModal(modelName string) TbdModal {
	return TbdModal{
		open:    false,
		content: fmt.Sprintf("⚠️ %s Under Construction ⚠️", modelName),
	}
}

func (m TbdModal) Init() tea.Cmd {
	return nil
}

func (m TbdModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m TbdModal) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.content, lipgloss.WithWhitespaceChars(" "))
}

func (m *TbdModal) Close() {
	m.open = false
}

func (m *TbdModal) Open() {
	m.open = true
}
