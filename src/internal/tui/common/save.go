package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SaveStartMsg struct{}
type SaveSuccessMsg struct{}
type SaveNoChangeMsg struct{}
type SaveFailMsg struct {
	Err error
}
type SaveCloseAndResetMsg struct{}
type SaveModelMsg struct {
	InQuitWorkflow bool
}

type SaveModel struct {
	Open               bool
	Save               bool
	FilePath           string
	Title              string
	Content            string
	Warning            string
	RenderedDuringQuit bool
	Help               HelpModel
}

func NewSaveModel(filepath string) SaveModel {
	help := NewHelpModel(true)
	help.ShortHelp = []key.Binding{CommonKeys.Confirm, CommonKeys.Cancel}
	return SaveModel{
		Help:     help,
		Title:    "Save OSCAL Model",
		Content:  fmt.Sprintf("Save changes to %s?", filepath),
		FilePath: filepath,
	}
}

func (m SaveModel) Init() tea.Cmd {
	return nil
}

func (m SaveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if m.Open {
			switch k {
			case ContainsKey(k, CommonKeys.Confirm.Keys()):
				if m.Save {
					var cmds []tea.Cmd
					cmds = append(cmds, func() tea.Msg {
						return SaveStartMsg{}
					})
					cmds = append(cmds, func() tea.Msg {
						return SaveModelMsg{
							InQuitWorkflow: m.RenderedDuringQuit,
						}
					})

					return m, tea.Sequence(cmds...)
				} else {
					m.Open = false
					return m, nil
				}

			case ContainsKey(k, CommonKeys.Cancel.Keys()):
				m.Open = false
				m.Save = false
			}
		}

	case SaveStartMsg:
		m.Title = "Saving..."
		m.Content = fmt.Sprintf("Saving to: %s", m.FilePath)
		m.Warning = ""
		m.Help.ShortHelp = []key.Binding{}

	case SaveFailMsg:
		m.Title = "Error saving model"
		m.Content = msg.Err.Error()
		m.Warning = "Changes not saved"
		m.Save = false

	case SaveSuccessMsg:
		m.Title = "Model saved!"
		m.Content = fmt.Sprintf("Model saved to %s", m.FilePath)
		m.Warning = ""
		m.Save = false

	case SaveCloseAndResetMsg:
		m.Title = "Save OSCAL Model"
		m.Content = fmt.Sprintf("Save changes to %s?", m.FilePath)
		m.Open = false
		m.Save = false
		m.Help.ShortHelp = []key.Binding{CommonKeys.Confirm, CommonKeys.Cancel}
	}
	return m, nil
}

func (m SaveModel) View() string {
	PrintToLog("in popup view")
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
