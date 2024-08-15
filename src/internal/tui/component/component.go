package component

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	tui "github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/pkg/common"
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
								if common.IsLulaLink(link) {
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
				frameworks = append(frameworks, framework{
					name:     k,
					controls: controls,
				})

			}
			components = append(components, component{
				uuid:       uuid,
				title:      c.Component.Title,
				desc:       c.Component.Description,
				frameworks: frameworks,
			})
		}
	}

	if len(components) > 0 {
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
	componentPicker.Style = tui.OverlayStyle

	frameworkPicker := viewport.New(pickerWidth, pickerHeight)
	frameworkPicker.Style = tui.OverlayStyle

	l := blist.New(viewedControls, tui.NewUnfocusedDelegate(), width, height)
	l.KeyMap = tui.FocusedListKeyMap()

	v := blist.New(viewedValidations, tui.NewUnfocusedDelegate(), width, height)
	v.KeyMap = tui.UnfocusedListKeyMap()

	controlPicker := viewport.New(width, height)
	controlPicker.Style = tui.PanelStyle

	remarks := viewport.New(width, height)
	remarks.Style = tui.PanelStyle
	description := viewport.New(width, height)
	description.Style = tui.PanelStyle
	validationPicker := viewport.New(width, height)
	validationPicker.Style = tui.PanelStyle

	return Model{
		content:           "Component Definition Content",
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
		keys:              componentKeys,
		help:              help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Printf("key: %s", msg.String())
		log.Printf("focus: %d", m.focus)
		if m.open {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "?":
				m.help.ShowAll = !m.help.ShowAll

			case "left", "h":
				if m.focus != 0 {
					m.focus--
					m.updateKeyBindings()
				}
			case "right", "l":
				if m.focus <= focusValidations {
					m.focus++
					m.updateKeyBindings()
				}
			case "up", "k":
				if m.inComponentOverlay && m.selectedComponentIndex > 0 {
					m.selectedComponentIndex--
					m.componentPicker.SetContent(m.updateComponentPickerContent())
				} else if m.inFrameworkOverlay && m.selectedFrameworkIndex > 0 {
					m.selectedFrameworkIndex--
					m.frameworkPicker.SetContent(m.updateFrameworkPickerContent())
				} else if m.focus == focusRemarks {
					m.remarks, _ = m.remarks.Update(msg)

				}
			case "down", "j":
				if m.inComponentOverlay && m.selectedComponentIndex < len(m.components)-1 {
					m.selectedComponentIndex++
					m.componentPicker.SetContent(m.updateComponentPickerContent())
				} else if m.inFrameworkOverlay && m.selectedFrameworkIndex < len(m.selectedComponent.frameworks)-1 {
					m.selectedFrameworkIndex++
					m.frameworkPicker.SetContent(m.updateFrameworkPickerContent())
				}
			case "enter":
				switch m.focus {
				case focusComponentSelection:
					if m.inComponentOverlay {
						m.selectedComponent = m.components[m.selectedComponentIndex]
						m.inComponentOverlay = false
					} else {
						m.inComponentOverlay = true
						m.componentPicker.SetContent(m.updateComponentPickerContent())
					}

				case focusFrameworkSelection:
					if m.inFrameworkOverlay {
						m.selectedFramework = m.components[m.selectedComponentIndex].frameworks[m.selectedFrameworkIndex]
						m.inFrameworkOverlay = false
					} else {
						m.inFrameworkOverlay = true
						m.frameworkPicker.SetContent(m.updateFrameworkPickerContent())
					}

				case focusControls:
					m.selectedControl = m.controls.SelectedItem().(control)
					m.remarks.SetContent(m.selectedControl.remarks)
					m.description.SetContent(m.selectedControl.desc)

					// update validations list for selected control
					validationItems := make([]blist.Item, len(m.selectedControl.validations))
					for i, val := range m.selectedControl.validations {
						validationItems[i] = val
					}
					m.validations.SetItems(validationItems)

				case focusValidations:
					if selectedItem := m.validations.SelectedItem(); selectedItem != nil {
						m.selectedValidation = selectedItem.(validationLink)
					}
				}
			case "esc", "q":
				if m.inComponentOverlay {
					m.inComponentOverlay = false
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
	totalHeight := m.height
	leftWidth := m.width / 4
	rightWidth := m.width - leftWidth - tui.PanelStyle.GetHorizontalPadding() - tui.PanelStyle.GetHorizontalMargins()

	// Add help panel at the top left
	helpStyle := lipgloss.NewStyle().Align(lipgloss.Right).Width(m.width - tui.PanelStyle.GetHorizontalPadding() - tui.PanelStyle.GetHorizontalMargins()).Height(1)
	helpView := helpStyle.Render(m.help.View(m.keys))

	// Add viewport styles
	focusedViewport := tui.PanelStyle.BorderForeground(tui.Focused)
	focusedViewportHeaderColor := tui.Focused
	focusedDialogBox := tui.DialogBoxStyle.BorderForeground(tui.Focused)

	selectedComponentDialogBox := tui.DialogBoxStyle
	selectedFrameworkDialogBox := tui.DialogBoxStyle
	controlPickerViewport := tui.PanelStyle
	controlHeaderColor := tui.Highlight
	descViewport := tui.PanelStyle
	descHeaderColor := tui.Highlight
	remarksViewport := tui.PanelStyle
	remarksHeaderColor := tui.Highlight
	validationPickerViewport := tui.PanelStyle
	validationHeaderColor := tui.Highlight

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
	selectedComponentLabel := tui.LabelStyle.Render("Selected Component")
	selectedComponentText := tui.TruncateText(getComponentText(m.selectedComponent), dialogFixedWidth)
	selectedComponentContent := selectedComponentDialogBox.Width(dialogFixedWidth).Render(selectedComponentText)
	selectedResult := lipgloss.JoinHorizontal(lipgloss.Top, selectedComponentLabel, selectedComponentContent)

	selectedFrameworkLabel := tui.LabelStyle.Render("Selected Framework")
	selectedFrameworkText := tui.TruncateText(getFrameworkText(m.selectedFramework), dialogFixedWidth)
	selectedFrameworkContent := selectedFrameworkDialogBox.Width(dialogFixedWidth).Render(selectedFrameworkText)
	selectedFramework := lipgloss.JoinHorizontal(lipgloss.Top, selectedFrameworkLabel, selectedFrameworkContent)

	componentSelectionContent := lipgloss.JoinHorizontal(lipgloss.Top, selectedResult, selectedFramework)

	// Build left panel with control picker
	topSectionHeight := helpStyle.GetHeight() + tui.DialogBoxStyle.GetHeight()
	bottomHeight := totalHeight - topSectionHeight
	m.controls.SetShowTitle(false)
	m.controls.SetHeight(m.height - tui.PanelTitleStyle.GetHeight() - tui.PanelStyle.GetVerticalPadding())
	m.controls.SetWidth(leftWidth - tui.PanelStyle.GetHorizontalPadding())

	m.controlPicker.Style = controlPickerViewport
	m.controlPicker.SetContent(m.controls.View())
	m.controlPicker.Height = bottomHeight
	m.controlPicker.Width = leftWidth - tui.PanelStyle.GetHorizontalPadding()
	leftView := fmt.Sprintf("%s\n%s", tui.HeaderView("Controls List", m.controlPicker.Width-tui.PanelStyle.GetMarginRight(), controlHeaderColor), m.controlPicker.View())

	// Add right panels with control details
	remarksHeight := bottomHeight / 5
	descriptionHeight := bottomHeight / 5
	validationsHeight := bottomHeight - remarksHeight - descriptionHeight - tui.PanelStyle.GetVerticalPadding() - tui.PanelStyle.GetPaddingBottom() - 3*tui.PanelTitleStyle.GetHeight() - 3

	m.remarks.Style = remarksViewport
	m.remarks.Height = remarksHeight // this needs to get passed back to update I think? because it's not changing the underlying model
	log.Printf("in view: remarks height: %d", m.remarks.Height)
	log.Printf("remarks: %s", m.remarks.View())
	m.remarks.Width = rightWidth
	m.remarks, _ = m.remarks.Update(tea.WindowSizeMsg{Width: rightWidth, Height: remarksHeight})

	m.description.Style = descViewport
	m.description.Height = descriptionHeight
	m.description.Width = rightWidth
	m.description, _ = m.description.Update(tea.WindowSizeMsg{Width: rightWidth, Height: descriptionHeight})

	m.validationPicker.Height = validationsHeight
	m.validationPicker.Width = rightWidth

	m.validations.SetShowTitle(false)
	m.validations.SetHeight(validationsHeight - tui.PanelTitleStyle.GetHeight() - tui.PanelStyle.GetVerticalPadding())
	m.validations.SetWidth(rightWidth - tui.PanelStyle.GetHorizontalPadding())

	m.validationPicker.Style = validationPickerViewport
	m.validationPicker.SetContent(m.validations.View())

	remarksPanel := fmt.Sprintf("%s\n%s", tui.HeaderView("Remarks", rightWidth-tui.PanelStyle.GetPaddingRight(), remarksHeaderColor), m.remarks.View())
	descriptionPanel := fmt.Sprintf("%s\n%s", tui.HeaderView("Description", rightWidth-tui.PanelStyle.GetPaddingRight(), descHeaderColor), m.description.View())
	validationsPanel := fmt.Sprintf("%s\n%s", tui.HeaderView("Validations", rightWidth-tui.PanelStyle.GetPaddingRight(), validationHeaderColor), m.validationPicker.View())

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
	// implement multi-select component picker...
	s := strings.Builder{}
	s.WriteString("Select one or many components:\n\n")

	for i, component := range m.components {
		if m.selectedComponentIndex == i {
			s.WriteString("[✔] ")
		} else {
			s.WriteString("[ ] ")
		}
		s.WriteString(getComponentText(component))
		s.WriteString("\n")
	}

	return s.String()
}

func (m Model) updateFrameworkPickerContent() string {
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

	return s.String()
}
