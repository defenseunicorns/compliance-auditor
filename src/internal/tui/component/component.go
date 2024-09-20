package component

import (
	"fmt"
	"slices"
	"strings"

	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	pkgcommon "github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
)

const (
	height           = 20
	width            = 12
	dialogFixedWidth = 40
)

const (
	componentPicker common.PickerKind = "component"
	frameworkPicker common.PickerKind = "framework"
)

// NewComponentDefinitionModel create new model for component definition view
func NewComponentDefinitionModel(oscalComponent *oscalTypes_1_1_2.ComponentDefinition) Model {
	var selectedComponent component
	var selectedFramework framework
	viewedControls := make([]blist.Item, 0)
	viewedValidations := make([]blist.Item, 0)
	components := make([]component, 0)
	frameworks := make([]framework, 0)

	if oscalComponent != nil {
		validationStore := validationstore.NewValidationStore()
		if oscalComponent.BackMatter != nil {
			validationStore = validationstore.NewValidationStoreFromBackMatter(*oscalComponent.BackMatter)
		}

		if oscalComponent.Components != nil {
			for cIdx, c := range *oscalComponent.Components {
				// for each component, add the control implementation to the framework
				componentFrameworks := make([]framework, 0)
				if c.ControlImplementations != nil {
					for ctrlImpIdx, controlImpl := range *c.ControlImplementations {
						// get the controls for each framework
						controls := make([]control, 0, len(controlImpl.ImplementedRequirements))
						for reqIdx, implementedRequirement := range controlImpl.ImplementedRequirements {
							// get validations from implementedRequirement.Links
							validationLinks := make([]validationLink, 0)
							if implementedRequirement.Links != nil {
								for _, link := range *implementedRequirement.Links {
									if pkgcommon.IsLulaLink(link) {
										validation, err := validationStore.GetLulaValidation(link.Href)
										if err == nil {
											// add the lula validation to the validations array
											validationLinks = append(validationLinks, validationLink{
												oscalLink:  &link, //&(*(*c.ControlImplementations)[ctrlImpIdx].ImplementedRequirements[reqIdx].Links)[linkIdx],
												text:       link.Text,
												validation: validation,
											})
										}
									}
								}
							}

							controls = append(controls, control{
								oscalControl: &(*c.ControlImplementations)[ctrlImpIdx].ImplementedRequirements[reqIdx], //&implementedRequirement,
								title:        implementedRequirement.ControlId,
								uuid:         implementedRequirement.UUID,
								validations:  validationLinks,
							})
						}

						// sort controls by title
						slices.SortStableFunc(controls, func(a, b control) int {
							return oscal.CompareControlsInt(a.title, b.title)
						})

						componentFrameworks = append(componentFrameworks, framework{
							oscalFramework: &(*c.ControlImplementations)[ctrlImpIdx], //&controlImpl,
							name:           controlImpl.Source,
							uuid:           controlImpl.UUID,
							controls:       controls,
						})

						// Add named framework if set
						status, value := oscal.GetProp("framework", oscal.LULA_NAMESPACE, controlImpl.Props)
						if status {
							componentFrameworks = append(componentFrameworks, framework{
								oscalFramework: &(*c.ControlImplementations)[ctrlImpIdx], //&controlImpl,
								name:           value,
								uuid:           controlImpl.UUID,
								controls:       controls,
							})
						}
					}
				}

				// sort componentFrameworks by name
				slices.SortStableFunc(componentFrameworks, func(a, b framework) int {
					return strings.Compare(a.name, b.name)
				})

				components = append(components, component{
					oscalComponent: &(*oscalComponent.Components)[cIdx], //&c,
					uuid:           c.UUID,
					title:          c.Title,
					desc:           c.Description,
					frameworks:     componentFrameworks,
				})
			}
		}
	}

	if len(components) > 0 {
		// sort components by title
		slices.SortStableFunc(components, func(a, b component) int {
			return strings.Compare(a.title, b.title)
		})

		selectedComponent = components[0]
		if len(selectedComponent.frameworks) > 0 {
			frameworks = selectedComponent.frameworks
			for _, fw := range selectedComponent.frameworks {
				selectedFramework = fw
				if len(selectedFramework.controls) > 0 {
					for _, c := range selectedFramework.controls {
						viewedControls = append(viewedControls, c)
					}
				}
				break
			}
		}
	}

	componentItems := make([]string, len(components))
	for i, c := range components {
		componentItems[i] = getComponentText(c)
	}
	componentPicker := common.NewPickerModel("Select a Component", componentPicker, componentItems, 0)

	frameworkItems := make([]string, len(frameworks))
	for i, f := range frameworks {
		frameworkItems[i] = getFrameworkText(f)
	}
	frameworkPicker := common.NewPickerModel("Select a Framework", frameworkPicker, frameworkItems, 0)

	l := blist.New(viewedControls, common.NewUnfocusedDelegate(), width, height)
	l.SetShowHelp(false) // help to be at top right
	l.KeyMap = common.FocusedListKeyMap()

	v := blist.New(viewedValidations, common.NewUnfocusedDelegate(), width, height)
	v.SetShowHelp(false) // help to be at top right
	v.KeyMap = common.UnfocusedListKeyMap()

	controlPicker := viewport.New(width, height)
	controlPicker.Style = common.PanelStyle

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

	validationPicker := viewport.New(width, height)
	validationPicker.Style = common.PanelStyle

	help := common.NewHelpModel(false)
	help.OneLine = true
	help.ShortHelp = shortHelpNoFocus

	return Model{
		keys:              componentKeys,
		help:              help,
		componentModel:    oscalComponent,
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
	}
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
	case tea.WindowSizeMsg:
		m.UpdateSizing(msg.Height-common.TabOffset, msg.Width)

	case tea.KeyMsg:
		if m.open {
			k := msg.String()
			switch k {
			case common.ContainsKey(k, m.keys.Help.Keys()):
				m.help.ShowAll = !m.help.ShowAll

			case common.ContainsKey(k, m.keys.NavigateLeft.Keys()):
				if !m.componentPicker.Open && !m.frameworkPicker.Open {
					if m.focus == 0 {
						m.focus = maxFocus
					} else {
						m.focus--
					}
					m.updateKeyBindings()
				}

			case common.ContainsKey(k, m.keys.NavigateRight.Keys()):
				if !m.componentPicker.Open && !m.frameworkPicker.Open {
					m.focus = (m.focus + 1) % (maxFocus + 1)
					m.updateKeyBindings()
				}

			case common.ContainsKey(k, m.keys.Confirm.Keys()):
				switch m.focus {
				case focusComponentSelection:
					if len(m.components) > 0 && !m.componentPicker.Open {
						return m, func() tea.Msg {
							return common.PickerOpenMsg{
								Kind: componentPicker,
							}
						}
					}

				case focusFrameworkSelection:
					if len(m.frameworks) > 0 && !m.frameworkPicker.Open {
						return m, func() tea.Msg {
							return common.PickerOpenMsg{
								Kind: frameworkPicker,
							}
						}
					}

				case focusControls:
					if selectedItem := m.controls.SelectedItem(); selectedItem != nil {
						m.selectedControl = m.controls.SelectedItem().(control)
						m.remarks.SetContent(m.selectedControl.oscalControl.Remarks)
						m.description.SetContent(m.selectedControl.oscalControl.Description)

						// update validations list for selected control
						validationItems := make([]blist.Item, len(m.selectedControl.validations))
						for i, val := range m.selectedControl.validations {
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
				if m.selectedControl.oscalControl != nil {
					switch m.focus {
					case focusRemarks:
						if !m.remarksEditor.Focused() {
							m.remarksEditor.SetValue(m.selectedControl.oscalControl.Remarks)
							m.remarks.SetContent(m.remarksEditor.View())
							_ = m.remarksEditor.Focus()
							m.updateKeyBindings()
						}
					case focusDescription:
						if !m.descriptionEditor.Focused() {
							m.descriptionEditor.SetValue(m.selectedControl.oscalControl.Description)
							m.description.SetContent(m.descriptionEditor.View())
							_ = m.descriptionEditor.Focus()
							m.updateKeyBindings()
						}
					}
				}

			case common.ContainsKey(k, m.keys.Cancel.Keys()):
				if m.selectedControl.oscalControl != nil {
					switch m.focus {
					case focusRemarks:
						if m.remarksEditor.Focused() {
							m.remarksEditor.Blur()
							m.remarks.SetContent(m.selectedControl.oscalControl.Remarks)
							m.updateKeyBindings()
						}

					case focusDescription:
						if m.descriptionEditor.Focused() {
							m.descriptionEditor.Blur()
							m.description.SetContent(m.selectedControl.oscalControl.Description)
							m.updateKeyBindings()
						}
					}
				}
			}
		}

	case common.PickerItemSelected:
		// reset all the controls, contents - if component is selected, reset the framework list as well
		if msg.From == componentPicker {
			m.selectedComponent = m.components[msg.Selected]
			m.selectedFramework = framework{}

			// Update controls list
			if len(m.components[msg.Selected].frameworks) > 0 {
				m.selectedFramework = m.components[msg.Selected].frameworks[0]
			}
		} else if msg.From == frameworkPicker {
			m.selectedFramework = m.selectedComponent.frameworks[msg.Selected]
		}
		if m.selectedFramework.oscalFramework != nil {
			controlItems := make([]blist.Item, len(m.selectedFramework.controls))
			if len(m.selectedFramework.controls) > 0 {
				for i, c := range m.selectedFramework.controls {
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

	mdl, cmd := m.componentPicker.Update(msg)
	m.componentPicker = mdl.(common.PickerModel)
	cmds = append(cmds, cmd)

	mdl, cmd = m.frameworkPicker.Update(msg)
	m.frameworkPicker = mdl.(common.PickerModel)
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

func getComponentText(component component) string {
	if component.uuid == "" {
		return "No Components"
	}
	return fmt.Sprintf("%s - %s", component.title, component.uuid)
}

func getFrameworkText(framework framework) string {
	if framework.name == "" {
		return "No Frameworks"
	}
	return framework.name
}
