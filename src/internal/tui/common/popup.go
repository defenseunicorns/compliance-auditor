package common

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	popupWidth  = 40
	popupHeight = 10
)

type WarnPopupModel struct {
	Open  bool
	Saved bool
	Help  HelpModel
}

func NewWarnPopup() WarnPopupModel {
	help := NewHelpModel(true)
	help.ShortHelp = []key.Binding{
		CommonKeys.Confirm, CommonKeys.Save, CommonKeys.Cancel,
	}
	return WarnPopupModel{
		Open:  false,
		Saved: false,
		Help:  help,
	}
}

func (m WarnPopupModel) Update(_ tea.Msg) (WarnPopupModel, tea.Cmd) {
	// update Saved?
	return m, nil
}

func (m WarnPopupModel) View() string {
	popupStyle := OverlayWarnStyle.
		Width(popupWidth).
		Height(popupHeight)

	title := "Quit Console"
	text := "Are you sure you want to quit the Lula Console?"
	if !m.Saved {
		text += "⚠️ Changes not written ⚠️"
	}
	content := lipgloss.JoinVertical(lipgloss.Top, title, text, NewHelpModel(true).View())
	return popupStyle.Render(content)
}
