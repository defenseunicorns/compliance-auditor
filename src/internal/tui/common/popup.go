package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WarnPopupModal struct {
	Open  bool
	Saved bool
	Keys  KeyMap
}

func New() WarnPopupModal {
	return WarnPopupModal{
		Open:  false,
		Saved: false,
		Keys:  CommonKeys,
	}
}

func (m WarnPopupModal) Update(_ tea.Msg) (WarnPopupModal, tea.Cmd) {
	// update Saved?
	return m, nil
}

func (m WarnPopupModal) View() string {
	title := "Quit Console"
	content := "Are you sure you want to quit the Lula Console?"
	if !m.Saved {
		content += "⚠️ Changes not written ⚠️"
	}
	return lipgloss.JoinVertical(lipgloss.Top, title, content, NewHelpModel(true).View(m.Keys))
}
