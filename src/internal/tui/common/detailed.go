package common

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ContentType string

const (
	detailWidth  = 80
	detailHeight = 20
	widthScale   = 0.9
	heightScale  = 0.9
)

type DetailOpenMsg struct {
	Content      string
	WindowHeight int
	WindowWidth  int
}

type DetailModel struct {
	Open            bool
	help            HelpModel
	contentViewport viewport.Model
	width           int
	height          int
}

func NewDetailModel() DetailModel {
	help := NewHelpModel(true)
	help.ShortHelp = ShortHelpDetail

	return DetailModel{
		help:            help,
		contentViewport: viewport.New(detailWidth, detailHeight),
	}
}

func (m DetailModel) Init() tea.Cmd {
	return nil
}

func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateSizing(int(float64(msg.Height)*heightScale), int(float64(msg.Width)*widthScale))

	case tea.KeyMsg:
		k := msg.String()
		if m.Open {
			switch k {
			case ContainsKey(k, CommonKeys.Cancel.Keys()):
				m.Open = false
			}
		}

	case DetailOpenMsg:
		m.Open = true
		m.contentViewport.GotoTop()
		m.contentViewport.SetContent(msg.Content)
		m.updateSizing(int(float64(msg.WindowHeight)*heightScale), int(float64(msg.WindowWidth)*widthScale))
	}

	m.contentViewport, cmd = m.contentViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m DetailModel) View() string {
	overlayDetailStyle := OverlayStyle.
		Width(m.width).
		Height(m.height)

	detailContent := lipgloss.JoinVertical(lipgloss.Top, m.contentViewport.View(), "\n", m.help.View())

	return overlayDetailStyle.Render(detailContent)
}

func (m *DetailModel) updateSizing(height, width int) {
	m.height = height
	m.width = width

	m.contentViewport.Height = height - OverlayStyle.GetVerticalPadding() - 2
	m.contentViewport.Width = width - OverlayStyle.GetHorizontalPadding()
}
