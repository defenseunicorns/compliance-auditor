package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	validationstore "github.com/defenseunicorns/lula/src/pkg/common/validation-store"
	"github.com/defenseunicorns/lula/src/types"
)

type componentModel struct {
	focus                  cFocus
	content                string
	keys                   componentKeys
	inComponentOverlay     bool
	components             []component
	selectedComponent      component
	selectedComponentIndex int
	selectedThreshold      threshold
	selectedThresholdIndex int
	componentPicker        viewport.Model
	controlPicker          viewport.Model
	controls               blist.Model
	selectedControl        control
	remarks                viewport.Model
	description            viewport.Model
	validationPicker       viewport.Model
	validations            blist.Model
	selectedValidation     validationLink
	open                   bool
	help                   help.Model
	width                  int
	height                 int
}

type cFocus int

const (
	focusControls cFocus = iota
	focusValidations
)

type component struct {
	uuid, remarks, title, desc string
	thresholds                 map[string]threshold
}

type threshold struct {
	controls []control
}

type validationLink struct {
	text       string
	validation *types.LulaValidation
}

func (i validationLink) Title() string       { return i.validation.Name }
func (i validationLink) Description() string { return i.text }
func (i validationLink) FilterValue() string { return i.validation.Name }

type control struct {
	uuid, remarks, title, desc string
	validations                []validationLink
}

func (i control) Title() string       { return i.title }
func (i control) Description() string { return i.uuid }
func (i control) FilterValue() string { return i.title }

func newDelegate() blist.DefaultDelegate {
	d := blist.NewDefaultDelegate()

	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{listHotkeys.Confirm, listHotkeys.Help}
	}

	return d
}

func NewComponentDefinitionModel(oscalComponent *oscalTypes_1_1_2.ComponentDefinition) componentModel {
	height := 60
	width := 12
	controls := make([]blist.Item, 0)
	validations := make([]blist.Item, 0)

	if oscalComponent != nil {
		// Build content of the component model
		// Need to get components x targets x control implementations (2 filters for component and target to sort the list of controls -> like a pop-up list would be cool)
		componentTargets := oscal.NewComponentTargets(oscalComponent)

		validationStore := validationstore.NewValidationStore()
		if oscalComponent.BackMatter != nil {
			validationStore = validationstore.NewValidationStoreFromBackMatter(*oscalComponent.BackMatter)
		}

		for _, component := range componentTargets {
			for _, target := range component.Targets {
				for _, controlImpl := range target {
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
			}
		}
	}

	l := blist.New(controls, newDelegate(), width, height)
	l.KeyMap = focusedListKeyMap()

	v := blist.New(validations, newDelegate(), width, height)
	v.KeyMap = unfocusedListKeyMap()

	controlPicker := viewport.New(width, height)
	controlPicker.Style = panelStyle

	remarks := viewport.New(width, height)
	remarks.Style = panelStyle
	description := viewport.New(width, height)
	description.Style = panelStyle
	validationPicker := viewport.New(width, height)
	validationPicker.Style = panelStyle

	return componentModel{
		content:          "Component Definition Content",
		controlPicker:    controlPicker,
		controls:         l,
		remarks:          remarks,
		description:      description,
		validationPicker: validationPicker,
		focus:            focusControls,
		validations:      v,
		keys:             componentHotkeys,
		help:             help.New(),
	}
}

func (m componentModel) Init() tea.Cmd {
	return nil
}

func (m componentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.open {
			switch msg.String() {
			case "esc":
				return m, tea.Quit

			case "left", "h", "right", "l":
				if m.focus == focusControls {
					m.focus = focusValidations
					m.controls.KeyMap = unfocusedListKeyMap()
					m.validations.KeyMap = focusedListKeyMap()
				} else if m.focus == focusValidations {
					m.focus = focusControls
					m.validations.KeyMap = unfocusedListKeyMap()
					m.controls.KeyMap = focusedListKeyMap()
				}

			case "enter":
				if m.focus == focusControls {
					m.selectedControl = m.controls.SelectedItem().(control)
					m.remarks.SetContent(m.selectedControl.remarks)
					m.description.SetContent(m.selectedControl.desc)

					// update validations list for selected control
					validationItems := make([]blist.Item, len(m.selectedControl.validations))
					for i, val := range m.selectedControl.validations {
						validationItems[i] = val
					}
					m.validations.SetItems(validationItems)

					m.focus = focusValidations
					m.controls.KeyMap = unfocusedListKeyMap()
					m.validations.KeyMap = focusedListKeyMap()

				} else if m.focus == focusValidations {
					if selectedItem := m.validations.SelectedItem(); selectedItem != nil {
						m.selectedValidation = selectedItem.(validationLink)
					}
				}
			}
		}
	}
	m.controls, cmd = m.controls.Update(msg)
	cmds = append(cmds, cmd)

	m.validations, cmd = m.validations.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m componentModel) View() string {
	totalHeight := m.height
	leftWidth := m.width / 4
	rightWidth := m.width - leftWidth - panelStyle.GetHorizontalPadding() - panelStyle.GetHorizontalMargins()

	// Build left panel with control picker
	m.controls.SetShowTitle(false)
	m.controls.SetHeight(m.height - panelTitleStyle.GetHeight() - panelStyle.GetVerticalPadding())
	m.controls.SetWidth(leftWidth - panelStyle.GetHorizontalPadding())

	var controlPickerViewport lipgloss.Style
	var controlHeaderColor lipgloss.AdaptiveColor
	if m.focus == focusControls {
		controlPickerViewport = panelStyle.BorderForeground(focused)
		controlHeaderColor = focused
	} else {
		controlPickerViewport = panelStyle
		controlHeaderColor = highlight
	}

	m.controlPicker.Style = controlPickerViewport
	m.controlPicker.SetContent(m.controls.View())
	m.controlPicker.Height = totalHeight
	m.controlPicker.Width = leftWidth - panelStyle.GetHorizontalPadding()
	leftView := fmt.Sprintf("%s\n%s", headerView("Controls List", m.controlPicker.Width-panelStyle.GetMarginRight(), controlHeaderColor), m.controlPicker.View())

	// Add right panels with control details
	remarksHeight := totalHeight / 5
	descriptionHeight := totalHeight / 5
	validationsHeight := totalHeight - remarksHeight - descriptionHeight - panelStyle.GetVerticalPadding() - panelStyle.GetPaddingBottom() - 3*panelTitleStyle.GetHeight()

	m.remarks.Height = remarksHeight
	m.remarks.Width = rightWidth
	m.description.Height = descriptionHeight
	m.description.Width = rightWidth
	m.validationPicker.Height = validationsHeight
	m.validationPicker.Width = rightWidth

	m.validations.SetShowTitle(false)
	m.validations.SetHeight(validationsHeight - panelTitleStyle.GetHeight() - panelStyle.GetVerticalPadding())
	m.validations.SetWidth(rightWidth - panelStyle.GetHorizontalPadding())

	var validationPickerViewport lipgloss.Style
	var validationHeaderColor lipgloss.AdaptiveColor
	if m.focus == focusValidations {
		validationPickerViewport = panelStyle.BorderForeground(focused) // Highlight the focused area
		validationHeaderColor = focused
	} else {
		validationPickerViewport = panelStyle
		validationHeaderColor = highlight
	}
	m.validationPicker.Style = validationPickerViewport
	m.validationPicker.SetContent(m.validations.View())

	remarksPanel := fmt.Sprintf("%s\n%s", headerView("Remarks", rightWidth-panelStyle.GetPaddingRight(), highlight), m.remarks.View())
	descriptionPanel := fmt.Sprintf("%s\n%s", headerView("Description", rightWidth-panelStyle.GetPaddingRight(), highlight), m.description.View())
	validationsPanel := fmt.Sprintf("%s\n%s", headerView("Validations", rightWidth-panelStyle.GetPaddingRight(), validationHeaderColor), m.validationPicker.View())

	rightView := lipgloss.JoinVertical(lipgloss.Top, remarksPanel, descriptionPanel, validationsPanel)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView)
}

// func (m componentModel) renderValidations() string {

// 	return "Validations"
// }

// func componentModelRender(oscalComponent *oscalTypes_1_1_2.ComponentDefinition, width, height int) string {

// }
