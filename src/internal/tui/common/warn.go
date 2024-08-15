package common

type warnType int

type WarnModal struct {
	open     bool
	warnType warnType
	title    string
	content  string
}

// func (m model) warnModalRender() string {
// 	title := m.warnModel.title
// 	content := m.warnModel.content
// 	confirm := modalConfirm.Render(" (" + hotkeys.Confirm[0] + ") Confirm ")
// 	cancel := modalCancel.Render(" (" + hotkeys.Quit[0] + ") Cancel ")
// 	tip := confirm + lipgloss.NewStyle().Background(background).Render("           ") + cancel
// 	return modalBorderStyle(modalHeight, modalWidth).Render(title + "\n\n" + content + "\n\n" + tip)
// }
