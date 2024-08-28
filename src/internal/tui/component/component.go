package component

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	blist "github.com/charmbracelet/bubbles/list"
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
	pickerHeight     = 20
	pickerWidth      = 80
	dialogFixedWidth = 40
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
		componentFrameworks := oscal.NewComponentFrameworks(oscalComponent)

		validationStore := validationstore.NewValidationStore()
		if oscalComponent.BackMatter != nil {
			validationStore = validationstore.NewValidationStoreFromBackMatter(*oscalComponent.BackMatter)
		}

		for uuid, c := range componentFrameworks {
			frameworks := make([]framework, 0)
			for k, f := range c.Frameworks {
				controls := make([]control, 0)

				for _, controlImpl := range f {
					for _, implementedRequirement := range controlImpl.ImplementedRequirements {
						// get validations from implementedRequirement.Links
						validationLinks := make([]validationLink, 0)
						if implementedRequirement.Links != nil {
							for _, link := range *implementedRequirement.Links {
								if pkgcommon.IsLulaLink(link) {
									validation, err := validationStore.GetLulaValidation(link.Href)
									if err == nil {
										// add the lula validation to the validations array
										validationLinks = append(validationLinks, validationLink{
											text:       link.Text,
											validation: validation,
										})
									}
								}
							}
						}

						controls = append(controls, control{
							title:       implementedRequirement.ControlId,
							uuid:        implementedRequirement.UUID,
							desc:        implementedRequirement.Description,
							remarks:     implementedRequirement.Remarks,
							validations: validationLinks,
						})
					}
				}
				// sort controls by title
				sort.Slice(controls, func(i, j int) bool {
					// custom sort function to sort controls by title
					return oscal.CompareControls(controls[i].title, controls[j].title)
				})

				frameworks = append(frameworks, framework{
					name:     k,
					controls: controls,
				})

			}
			// sort frameworks by name
			sort.Slice(frameworks, func(i, j int) bool {
				return frameworks[i].name < frameworks[j].name
			})

			components = append(components, component{
				uuid:       uuid,
				title:      c.Component.Title,
				desc:       c.Component.Description,
				frameworks: frameworks,
			})
		}
	}

	if len(components) > 0 {
		// sort components by title
		sort.Slice(components, func(i, j int) bool {
			return components[i].title < components[j].title
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

	componentPicker := viewport.New(pickerWidth, pickerHeight)
	componentPicker.Style = common.OverlayStyle

	frameworkPicker := viewport.New(pickerWidth, pickerHeight)
	frameworkPicker.Style = common.OverlayStyle

	l := blist.New(viewedControls, common.NewUnfocusedDelegate(), width, height)
	l.KeyMap = common.FocusedListKeyMap()

	v := blist.New(viewedValidations, common.NewUnfocusedDelegate(), width, height)
	v.KeyMap = common.UnfocusedListKeyMap()

	controlPicker := viewport.New(width, height)
	controlPicker.Style = common.PanelStyle

	remarks := viewport.New(width, height)
	remarks.Style = common.PanelStyle
	description := viewport.New(width, height)
	description.Style = common.PanelStyle
	validationPicker := viewport.New(width, height)
	validationPicker.Style = common.PanelStyle

	return Model{
		keys:              componentKeys,
		help:              help.New(),
		components:        components,
		selectedComponent: selectedComponent,
		componentPicker:   componentPicker,
		frameworks:        frameworks,
		selectedFramework: selectedFramework,
		frameworkPicker:   frameworkPicker,
		controlPicker:     controlPicker,
		controls:          l,
		remarks:           remarks,
		description:       description,
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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateSizing(msg.Height-common.TabOffset, msg.Width)

	case tea.KeyMsg:
		if m.open {
			k := msg.String()
			switch k {

			case common.ContainsKey(k, m.keys.Quit.Keys()):
				return m, tea.Quit

			case common.ContainsKey(k, m.keys.Help.Keys()):
				m.help.ShowAll = !m.help.ShowAll

			case common.ContainsKey(k, m.keys.NavigateLeft.Keys()):
				if m.focus == 0 {
					m.focus = maxFocus
				} else {
					m.focus--
				}
				m.updateKeyBindings()

			case common.ContainsKey(k, m.keys.NavigateRight.Keys()):
				m.focus = (m.focus + 1) % (maxFocus + 1)
				m.updateKeyBindings()

			case common.ContainsKey(k, m.keys.Up.Keys()):
				if m.inComponentOverlay && m.selectedComponentIndex > 0 {
					m.selectedComponentIndex--
					m.componentPicker.SetContent(m.updateComponentPickerContent())
				} else if m.inFrameworkOverlay && m.selectedFrameworkIndex > 0 {
					m.selectedFrameworkIndex--
					m.frameworkPicker.SetContent(m.updateFrameworkPickerContent())
				}

			case common.ContainsKey(k, m.keys.Down.Keys()):
				if m.inComponentOverlay && m.selectedComponentIndex < len(m.components)-1 {
					m.selectedComponentIndex++
					m.componentPicker.SetContent(m.updateComponentPickerContent())
				} else if m.inFrameworkOverlay && m.selectedFrameworkIndex < len(m.selectedComponent.frameworks)-1 {
					m.selectedFrameworkIndex++
					m.frameworkPicker.SetContent(m.updateFrameworkPickerContent())
				}

			case common.ContainsKey(k, m.keys.Confirm.Keys()):
				switch m.focus {
				case focusComponentSelection:
					if m.inComponentOverlay {
						if len(m.components) > 1 {
							m.selectedComponent = m.components[m.selectedComponentIndex]
							m.selectedFrameworkIndex = 0

							// Update controls list
							if len(m.components[m.selectedComponentIndex].frameworks) > 0 {
								m.selectedFramework = m.components[m.selectedComponentIndex].frameworks[m.selectedFrameworkIndex]
							} else {
								m.selectedFramework = framework{}
							}
							controlItems := make([]blist.Item, len(m.selectedFramework.controls))
							if len(m.selectedFramework.controls) > 0 {
								for i, c := range m.selectedFramework.controls {
									controlItems[i] = c
								}
							}
							m.controls.SetItems(controlItems)
							m.controls.SetDelegate(common.NewUnfocusedDelegate())

							// Update remarks, description, and validations
							m.remarks.SetContent("")
							m.description.SetContent("")
							m.validations.SetItems(make([]blist.Item, 0))
						}

						m.inComponentOverlay = false
					} else {
						m.inComponentOverlay = true
						m.componentPicker.SetContent(m.updateComponentPickerContent())
					}
				case focusFrameworkSelection:
					if m.inFrameworkOverlay {
						if len(m.components) != 0 && len(m.components[m.selectedComponentIndex].frameworks) > 1 {
							m.selectedFramework = m.components[m.selectedComponentIndex].frameworks[m.selectedFrameworkIndex]

							// Update controls list
							controlItems := make([]blist.Item, len(m.selectedFramework.controls))
							if len(m.selectedFramework.controls) > 0 {
								for i, c := range m.selectedFramework.controls {
									controlItems[i] = c
								}
							}
							m.controls.SetItems(controlItems)
							m.controls.SetDelegate(common.NewUnfocusedDelegate())

							// Update remarks, description, and validations
							m.remarks.SetContent("")
							m.description.SetContent("")
							m.validations.SetItems(make([]blist.Item, 0))
						}

						m.inFrameworkOverlay = false
					} else {
						m.inFrameworkOverlay = true
						m.frameworkPicker.SetContent(m.updateFrameworkPickerContent())
					}

				case focusControls:
					if selectedItem := m.controls.SelectedItem(); selectedItem != nil {
						m.selectedControl = m.controls.SelectedItem().(control)
						m.remarks.SetContent(m.selectedControl.remarks)
						m.description.SetContent(m.selectedControl.desc)

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
				}

			case common.ContainsKey(k, m.keys.Cancel.Keys()):
				if m.inComponentOverlay {
					m.inComponentOverlay = false
				} else if m.inFrameworkOverlay {
					m.inFrameworkOverlay = false
				}
			}
		}
	}

	m.componentPicker, cmd = m.componentPicker.Update(msg)
	cmds = append(cmds, cmd)

	m.frameworkPicker, cmd = m.frameworkPicker.Update(msg)
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
	if m.inComponentOverlay {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.componentPicker.View(), lipgloss.WithWhitespaceChars(" "))
	}
	if m.inFrameworkOverlay {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.frameworkPicker.View(), lipgloss.WithWhitespaceChars(" "))
	}
	return m.mainView()
}

func (m Model) mainView() string {
	// Add help panel at the top left
	helpStyle := common.HelpStyle(m.width)
	helpView := helpStyle.Render(m.help.View(m.keys))

	// Add viewport styles
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

	remarksPanel := fmt.Sprintf("%s\n%s", common.HeaderView("Remarks", m.remarks.Width-common.PanelStyle.GetPaddingRight(), remarksHeaderColor), m.remarks.View())
	descriptionPanel := fmt.Sprintf("%s\n%s", common.HeaderView("Description", m.description.Width-common.PanelStyle.GetPaddingRight(), descHeaderColor), m.description.View())
	validationsPanel := fmt.Sprintf("%s\n%s", common.HeaderView("Validations", m.validationPicker.Width-common.PanelStyle.GetPaddingRight(), validationHeaderColor), m.validationPicker.View())

	rightView := lipgloss.JoinVertical(lipgloss.Top, remarksPanel, descriptionPanel, validationsPanel)
	bottomContent := lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView)

	return lipgloss.JoinVertical(lipgloss.Top, helpView, componentSelectionContent, bottomContent)
}

func getComponentText(component component) string {
	if component.uuid == "" {
		return "No Component Selected"
	}
	return fmt.Sprintf("%s - %s", component.title, component.uuid)
}

func getFrameworkText(framework framework) string {
	return framework.name
}

func (m Model) updateComponentPickerContent() string {
	helpStyle := common.HelpStyle(pickerWidth)
	helpView := helpStyle.Render(help.New().View(common.PickerHotkeys))

	s := strings.Builder{}
	s.WriteString("Select a Component:\n\n")

	for i, component := range m.components {
		if m.selectedComponentIndex == i {
			s.WriteString("(•) ") //[✔] Todo: many components?
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(getComponentText(component))
		s.WriteString("\n")
	}

	return lipgloss.JoinVertical(lipgloss.Top, helpView, s.String())
}

func (m Model) updateFrameworkPickerContent() string {
	helpStyle := common.HelpStyle(pickerWidth)
	helpView := helpStyle.Render(help.New().View(common.PickerHotkeys))

	s := strings.Builder{}
	s.WriteString("Select a Framework:\n\n")

	for i, fw := range m.selectedComponent.frameworks {
		if m.selectedFrameworkIndex == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(getFrameworkText(fw))
		s.WriteString("\n")
	}

	return lipgloss.JoinVertical(lipgloss.Top, helpView, s.String())
}
