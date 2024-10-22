package component

import (
	"fmt"

	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

const (
	componentPickerKind common.PickerKind = "component"
	frameworkPickerKind common.PickerKind = "framework"
	height                                = 20
	width                                 = 12
	dialogFixedWidth                      = 40
)

type Model struct {
	IsOpen             bool
	help               common.HelpModel
	keys               keys
	focus              focus
	oscalComponent     *oscalTypes_1_1_2.ComponentDefinition
	components         []component
	selectedComponent  component
	componentPicker    common.PickerModel
	frameworks         []framework
	selectedFramework  framework
	frameworkPicker    common.PickerModel
	controlPicker      viewport.Model
	controls           blist.Model
	selectedControl    control
	remarks            viewport.Model
	remarksEditor      textarea.Model
	description        viewport.Model
	descriptionEditor  textarea.Model
	validationPicker   viewport.Model
	validations        blist.Model
	selectedValidation validationLink
	detailView         common.DetailModel
	validateModel      ValidateModel
	width              int
	height             int
}

type ModelOpenMsg struct {
	Height int
	Width  int
}
type ModelCloseMsg struct{}

func NewComponentDefinitionModel(oscalComponent *oscalTypes_1_1_2.ComponentDefinition) Model {
	var selectedComponent component
	var selectedFramework framework
	viewedControls := make([]blist.Item, 0)
	viewedValidations := make([]blist.Item, 0)
	components := make([]component, 0)
	frameworks := make([]framework, 0)

	help := common.NewHelpModel(false)
	help.OneLine = true
	help.ShortHelp = shortHelpNoFocus

	// Initialize widgets
	componentPicker := common.NewPickerModel("Select a Component", componentPickerKind, []string{}, 0)
	frameworkPicker := common.NewPickerModel("Select a Framework", frameworkPickerKind, []string{}, 0)

	l := blist.New(viewedControls, common.NewUnfocusedDelegate(), width, height)
	l.SetShowHelp(false) // help to be at top right
	l.KeyMap = common.FocusedListKeyMap()

	v := blist.New(viewedValidations, common.NewUnfocusedDelegate(), width, height)
	v.SetShowHelp(false) // help to be at top right
	v.KeyMap = common.UnfocusedListKeyMap()

	controlPicker := viewport.New(width, height)
	controlPicker.Style = common.PanelStyle

	validationPicker := viewport.New(width, height)
	validationPicker.Style = common.PanelStyle

	remarks := viewport.New(width, height)
	remarks.Style = common.PanelStyle
	remarks.MouseWheelEnabled = false
	remarksEditor := textarea.New()
	remarksEditor.CharLimit = 0
	remarksEditor.KeyMap = common.UnfocusedTextAreaKeyMap()

	description := viewport.New(width, height)
	description.Style = common.PanelStyle
	description.MouseWheelEnabled = false
	descriptionEditor := textarea.New()
	descriptionEditor.CharLimit = 0
	descriptionEditor.KeyMap = common.UnfocusedTextAreaKeyMap()

	model := Model{
		keys:              componentKeys,
		help:              help,
		components:        components,
		selectedComponent: selectedComponent,
		componentPicker:   componentPicker,
		frameworks:        frameworks,
		selectedFramework: selectedFramework,
		frameworkPicker:   frameworkPicker,
		controlPicker:     controlPicker,
		controls:          l,
		remarks:           remarks,
		remarksEditor:     remarksEditor,
		description:       description,
		descriptionEditor: descriptionEditor,
		validationPicker:  validationPicker,
		validations:       v,
		detailView:        common.NewDetailModel(),
		validateModel:     NewValidateModel(oscalComponent),
	}

	model.UpdateWithComponentDefinition(oscalComponent)

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// up front so it doesn't capture the first key ('e')
	if m.remarksEditor.Focused() {
		m.remarksEditor, cmd = m.remarksEditor.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.descriptionEditor.Focused() {
		m.descriptionEditor, cmd = m.descriptionEditor.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case ModelOpenMsg:
		m.Open(msg.Height, msg.Width)

	case tea.WindowSizeMsg:
		m.updateSizing(msg.Height-common.TabOffset, msg.Width)

	case tea.KeyMsg:
		k := msg.String()
		switch k {
		case common.ContainsKey(k, m.keys.Help.Keys()):
			m.help.ShowAll = !m.help.ShowAll

		case common.ContainsKey(k, m.keys.NavigateLeft.Keys()):
			if !m.inOverlay() {
				if m.focus == 0 {
					m.focus = maxFocus
				} else {
					m.focus--
				}
				m.updateKeyBindings()
			}

		case common.ContainsKey(k, m.keys.NavigateRight.Keys()):
			if !m.inOverlay() {
				m.focus = (m.focus + 1) % (maxFocus + 1)
				m.updateKeyBindings()
			}

		case common.ContainsKey(k, m.keys.Confirm.Keys()):
			switch m.focus {
			case focusComponentSelection:
				if len(m.components) > 0 && !m.inOverlay() {
					return m, func() tea.Msg {
						return common.PickerOpenMsg{
							Kind: componentPickerKind,
						}
					}
				}

			case focusFrameworkSelection:
				if len(m.frameworks) > 0 && !m.inOverlay() {
					return m, func() tea.Msg {
						return common.PickerOpenMsg{
							Kind: frameworkPickerKind,
						}
					}
				}

			case focusControls:
				if selectedItem := m.controls.SelectedItem(); selectedItem != nil {
					m.selectedControl = m.controls.SelectedItem().(control)
					m.remarks.SetContent(m.selectedControl.OscalControl.Remarks)
					m.description.SetContent(m.selectedControl.OscalControl.Description)

					// update validations list for selected control
					validationItems := make([]blist.Item, len(m.selectedControl.Validations))
					for i, val := range m.selectedControl.Validations {
						validationItems[i] = val
					}
					m.validations.SetItems(validationItems)
				}

			case focusValidations:
				if selectedItem := m.validations.SelectedItem(); selectedItem != nil {
					m.selectedValidation = selectedItem.(validationLink)
				}

			case focusRemarks:
				if m.remarksEditor.Focused() {
					remarks := m.remarksEditor.Value()
					m.UpdateRemarks(remarks)
					m.remarksEditor.Blur()
					m.remarks.SetContent(remarks)
					m.updateKeyBindings()
				}

			case focusDescription:
				if m.descriptionEditor.Focused() {
					description := m.descriptionEditor.Value()
					m.UpdateDescription(description)
					m.descriptionEditor.Blur()
					m.description.SetContent(description)
					m.updateKeyBindings()
				}
			}

		case common.ContainsKey(k, m.keys.Edit.Keys()):
			if m.selectedControl.OscalControl != nil {
				switch m.focus {
				case focusRemarks:
					if !m.remarksEditor.Focused() {
						m.remarksEditor.SetValue(m.selectedControl.OscalControl.Remarks)
						m.remarks.SetContent(m.remarksEditor.View())
						_ = m.remarksEditor.Focus()
						m.updateKeyBindings()
					}
				case focusDescription:
					if !m.descriptionEditor.Focused() {
						m.descriptionEditor.SetValue(m.selectedControl.OscalControl.Description)
						m.description.SetContent(m.descriptionEditor.View())
						_ = m.descriptionEditor.Focus()
						m.updateKeyBindings()
					}
				}
			}

		case common.ContainsKey(k, m.keys.Detail.Keys()):
			switch m.focus {
			case focusValidations:
				// TODO: update the key locks
				if selectedItem := m.validations.SelectedItem(); selectedItem != nil {
					valLink := selectedItem.(validationLink)
					m.validations.KeyMap = common.UnfocusedListKeyMap()
					return m, func() tea.Msg {
						return common.DetailOpenMsg{
							Content:      getValidationText(valLink),
							WindowHeight: m.height,
							WindowWidth:  m.width,
						}
					}
				}
			}

		case common.ContainsKey(k, m.keys.Validate.Keys()):
			if !m.inOverlay() {
				return m, func() tea.Msg {
					return ValidateOpenMsg{
						Target: m.selectedFramework.Name,
						Height: m.height,
						Width:  m.width,
					}
				}
			}

		case common.ContainsKey(k, m.keys.Cancel.Keys()):
			if m.selectedControl.OscalControl != nil {
				switch m.focus {
				case focusRemarks:
					if m.remarksEditor.Focused() {
						m.remarksEditor.Blur()
						m.remarks.SetContent(m.selectedControl.OscalControl.Remarks)
					}

				case focusDescription:
					if m.descriptionEditor.Focused() {
						m.descriptionEditor.Blur()
						m.description.SetContent(m.selectedControl.OscalControl.Description)
					}
				}
			}
			m.updateKeyBindings()
		}

	case common.PickerItemSelected:
		// reset all the controls, contents - if component is selected, reset the framework list as well
		if msg.From == componentPickerKind {
			m.selectedComponent = m.components[msg.Selected]
			m.selectedFramework = framework{}

			// Update controls list
			if len(m.components[msg.Selected].Frameworks) > 0 {
				m.selectedFramework = m.components[msg.Selected].Frameworks[0]
			}
		} else if msg.From == frameworkPickerKind {
			m.selectedFramework = m.selectedComponent.Frameworks[msg.Selected]
		}

		m.resetWidgets()
	}

	mdl, cmd := m.componentPicker.Update(msg)
	m.componentPicker = mdl.(common.PickerModel)
	cmds = append(cmds, cmd)

	mdl, cmd = m.frameworkPicker.Update(msg)
	m.frameworkPicker = mdl.(common.PickerModel)
	cmds = append(cmds, cmd)

	mdl, cmd = m.detailView.Update(msg)
	m.detailView = mdl.(common.DetailModel)
	cmds = append(cmds, cmd)

	mdl, cmd = m.validateModel.Update(msg)
	m.validateModel = mdl.(ValidateModel)
	cmds = append(cmds, cmd)

	m.remarks, cmd = m.remarks.Update(msg)
	cmds = append(cmds, cmd)

	m.description, cmd = m.description.Update(msg)
	cmds = append(cmds, cmd)

	m.controls, cmd = m.controls.Update(msg)
	cmds = append(cmds, cmd)

	m.validations, cmd = m.validations.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.componentPicker.Open {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.componentPicker.View(), lipgloss.WithWhitespaceChars(" "))
	}
	if m.frameworkPicker.Open {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.frameworkPicker.View(), lipgloss.WithWhitespaceChars(" "))
	}
	if m.detailView.Open {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.detailView.View(), lipgloss.WithWhitespaceChars(" "))
	}
	if m.validateModel.IsOpen {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.validateModel.View(), lipgloss.WithWhitespaceChars(" "))
	}

	return m.mainView()
}

func (m Model) mainView() string {
	// Add viewport and focus styles
	focusedViewport := common.PanelStyle.BorderForeground(common.Focused)
	focusedViewportHeaderColor := common.Focused
	focusedDialogBox := common.DialogBoxStyle.BorderForeground(common.Focused)

	selectedComponentDialogBox := common.DialogBoxStyle
	selectedFrameworkDialogBox := common.DialogBoxStyle
	controlPickerViewport := common.PanelStyle
	controlHeaderColor := common.Highlight
	descViewport := common.PanelStyle
	descHeaderColor := common.Highlight
	remarksViewport := common.PanelStyle
	remarksHeaderColor := common.Highlight
	validationPickerViewport := common.PanelStyle
	validationHeaderColor := common.Highlight

	switch m.focus {
	case focusComponentSelection:
		selectedComponentDialogBox = focusedDialogBox
	case focusFrameworkSelection:
		selectedFrameworkDialogBox = focusedDialogBox
	case focusControls:
		controlPickerViewport = focusedViewport
		controlHeaderColor = focusedViewportHeaderColor
	case focusDescription:
		descViewport = focusedViewport
		descHeaderColor = focusedViewportHeaderColor
	case focusRemarks:
		remarksViewport = focusedViewport
		remarksHeaderColor = focusedViewportHeaderColor
	case focusValidations:
		validationPickerViewport = focusedViewport
		validationHeaderColor = focusedViewportHeaderColor
	}
	// Add help panel at the top right
	helpStyle := common.HelpStyle(m.width)
	helpView := helpStyle.Render(m.help.View())

	// Add widgets for dialogs
	selectedComponentLabel := common.LabelStyle.Render("Selected Component")
	selectedComponentText := common.TruncateText(getComponentText(m.selectedComponent), dialogFixedWidth)
	selectedComponentContent := selectedComponentDialogBox.Width(dialogFixedWidth).Render(selectedComponentText)
	selectedResult := lipgloss.JoinHorizontal(lipgloss.Top, selectedComponentLabel, selectedComponentContent)

	selectedFrameworkLabel := common.LabelStyle.Render("Selected Framework")
	selectedFrameworkText := common.TruncateText(getFrameworkText(m.selectedFramework), dialogFixedWidth)
	selectedFrameworkContent := selectedFrameworkDialogBox.Width(dialogFixedWidth).Render(selectedFrameworkText)
	selectedFramework := lipgloss.JoinHorizontal(lipgloss.Top, selectedFrameworkLabel, selectedFrameworkContent)

	componentSelectionContent := lipgloss.JoinHorizontal(lipgloss.Top, selectedResult, selectedFramework)

	m.controls.SetShowTitle(false)
	m.validations.SetShowTitle(false)

	m.controlPicker.Style = controlPickerViewport
	m.controlPicker.SetContent(m.controls.View())
	leftView := fmt.Sprintf("%s\n%s", common.HeaderView("Controls List", m.controlPicker.Width-common.PanelStyle.GetMarginRight(), controlHeaderColor), m.controlPicker.View())

	m.remarks.Style = remarksViewport
	m.description.Style = descViewport

	m.validationPicker.Style = validationPickerViewport
	m.validationPicker.SetContent(m.validations.View())

	// remarksView = m.remarks.View()
	if m.remarksEditor.Focused() {
		m.remarks.SetContent(lipgloss.JoinVertical(lipgloss.Top, m.remarksEditor.View()))
	} else if m.descriptionEditor.Focused() {
		m.description.SetContent(lipgloss.JoinVertical(lipgloss.Top, m.descriptionEditor.View()))
	}

	remarksPanel := fmt.Sprintf("%s\n%s", common.HeaderView("Remarks", m.remarks.Width-common.PanelStyle.GetPaddingRight(), remarksHeaderColor), m.remarks.View())
	descriptionPanel := fmt.Sprintf("%s\n%s", common.HeaderView("Description", m.description.Width-common.PanelStyle.GetPaddingRight(), descHeaderColor), m.description.View())
	validationsPanel := fmt.Sprintf("%s\n%s", common.HeaderView("Validations", m.validationPicker.Width-common.PanelStyle.GetPaddingRight(), validationHeaderColor), m.validationPicker.View())

	rightView := lipgloss.JoinVertical(lipgloss.Top, remarksPanel, descriptionPanel, validationsPanel)
	bottomContent := lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView)

	return lipgloss.JoinVertical(lipgloss.Top, helpView, componentSelectionContent, bottomContent)
}

// UpdateWithComponentDefinition updates the model data given a component definition
// Useful when a new component definition is loaded (at init time or if component definition is generated)
func (m *Model) UpdateWithComponentDefinition(oscalComponent *oscalTypes_1_1_2.ComponentDefinition) {
	var selectedComponent component
	var selectedFramework framework
	components := make([]component, 0)
	frameworks := make([]framework, 0)

	m.components = components
	m.frameworks = frameworks
	m.selectedComponent = selectedComponent
	m.selectedFramework = selectedFramework

	// Update data if component definition is not nil
	if oscalComponent != nil {
		m.oscalComponent = oscalComponent
		components = GetComponents(oscalComponent)
	}

	if len(components) > 0 {
		m.components = components
		m.selectedComponent = components[0]
		if len(m.selectedComponent.Frameworks) > 0 {
			m.frameworks = m.selectedComponent.Frameworks
			if len(m.frameworks) > 0 {
				m.selectedFramework = m.frameworks[0]
			}
		}
	}

	componentItems := make([]string, len(m.components))
	for i, c := range m.components {
		componentItems[i] = getComponentText(c)
	}
	m.componentPicker.UpdateItems(componentItems)

	frameworkItems := make([]string, len(m.frameworks))
	for i, f := range m.frameworks {
		frameworkItems[i] = getFrameworkText(f)
	}
	m.frameworkPicker.UpdateItems(frameworkItems)

	m.resetWidgets()
}

func (m *Model) Close() {
	m.IsOpen = false
}

func (m *Model) Open(height, width int) {
	m.IsOpen = true
	m.updateSizing(height, width)
}

func (m *Model) GetDimensions() (height, width int) {
	return m.height, m.width
}

// GetComponentDefinition returns the component definition model, used on save events
func (m *Model) GetComponentDefinition() *oscalTypes_1_1_2.ComponentDefinition {
	return m.oscalComponent
}

// GetSelectedFramework returns the selected framework
func (m *Model) GetSelectedFramework() framework {
	return m.selectedFramework
}

// TestSetSelectedControl is a test helper function to set the selected control
func (m *Model) TestSetSelectedControl(title string) {
	var idx int
	if len(m.selectedFramework.Controls) > 0 {
		for i, c := range m.selectedFramework.Controls {
			if c.Name == title {
				m.selectedControl = c
				idx = i
				break
			}
		}
	}
	m.controls.Select(idx)
}

func (m *Model) UpdateRemarks(remarks string) {
	// TODO: handle any race conditions updating the control?
	if m.selectedControl.OscalControl != nil {
		m.selectedControl.OscalControl.Remarks = remarks
	}
}

func (m *Model) UpdateDescription(description string) {
	// TODO: handle any race conditions updating the control?
	if m.selectedControl.OscalControl != nil {
		m.selectedControl.OscalControl.Description = description
	}
}

func (m *Model) inOverlay() bool {
	if m.componentPicker.Open || m.frameworkPicker.Open || m.detailView.Open || m.validateModel.IsOpen {
		return true
	}
	return false
}

func (m *Model) resetWidgets() {
	if m.selectedFramework.OscalFramework != nil {
		controlItems := make([]blist.Item, len(m.selectedFramework.Controls))
		if len(m.selectedFramework.Controls) > 0 {
			for i, c := range m.selectedFramework.Controls {
				controlItems[i] = c
			}
		}
		m.controls.SetItems(controlItems)
		m.controls.ResetSelected()
	}

	// Update remarks, description, and validations
	m.controls.SetDelegate(common.NewUnfocusedDelegate())
	m.controls.ResetSelected()
	m.controls.ResetFilter()
	m.selectedControl = control{}
	m.remarks.SetContent("")
	m.remarksEditor.SetValue("")
	m.description.SetContent("")
	m.descriptionEditor.SetValue("")
	m.validations.SetItems(make([]blist.Item, 0))
	m.validations.ResetSelected()
	m.validations.ResetFilter()
	m.selectedValidation = validationLink{}
}

func (m *Model) updateSizing(height, width int) {
	m.height = height
	m.width = width

	// Set internal sizing properties
	totalHeight := m.height - 4 // Subtract vertical margin/padding not accounted for below
	leftWidth := m.width / 4
	rightWidth := m.width - leftWidth - common.PanelStyle.GetHorizontalPadding() - common.PanelStyle.GetHorizontalMargins()

	topSectionHeight := common.HelpStyle(m.width).GetHeight() + common.DialogBoxStyle.GetHeight()
	bottomSectionHeight := totalHeight - topSectionHeight
	panelHeight := common.PanelTitleStyle.GetHeight() + common.PanelStyle.GetVerticalPadding() + common.PanelStyle.GetVerticalMargins() + 1 // 1 for border

	remarksOutsideHeight := bottomSectionHeight / 4
	remarksInsideHeight := remarksOutsideHeight - panelHeight

	descriptionOutsideHeight := bottomSectionHeight / 4
	descriptionInsideHeight := descriptionOutsideHeight - panelHeight
	validationsHeight := bottomSectionHeight - remarksOutsideHeight - descriptionOutsideHeight - 2*common.PanelTitleStyle.GetHeight()

	// Update widget sizing
	m.help.Width = m.width

	m.controls.SetHeight(m.height - common.PanelTitleStyle.GetHeight() - 1)
	m.controls.SetWidth(leftWidth - common.PanelStyle.GetHorizontalPadding())

	m.controlPicker.Height = bottomSectionHeight
	m.controlPicker.Width = leftWidth - common.PanelStyle.GetHorizontalPadding()

	m.remarks.Height = remarksInsideHeight
	m.remarks.Width = rightWidth
	m.remarks, _ = m.remarks.Update(tea.WindowSizeMsg{Width: rightWidth, Height: remarksInsideHeight}) // rebuild remarks for line wrapping?

	m.remarksEditor.SetHeight(m.remarks.Height - 1)
	m.remarksEditor.SetWidth(m.remarks.Width - 5) // probably need to fix this to be a func

	m.description.Height = descriptionInsideHeight
	m.description.Width = rightWidth
	m.description, _ = m.description.Update(tea.WindowSizeMsg{Width: rightWidth, Height: descriptionInsideHeight})

	m.descriptionEditor.SetHeight(m.description.Height - 1)
	m.descriptionEditor.SetWidth(m.description.Width - 5) // probably need to fix this to be a func

	m.validations.SetHeight(validationsHeight - common.PanelTitleStyle.GetHeight())
	m.validations.SetWidth(rightWidth - common.PanelStyle.GetHorizontalPadding())

	m.validationPicker.Height = validationsHeight
	m.validationPicker.Width = rightWidth
}
