package tui

import (
	"fmt"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
)

type model struct {
	cursor                    int
	tabs                      []string
	activeTab                 int
	content                   string
	focusPanel                int
	oscalModel                oscalTypes_1_1_2.OscalCompleteSchema
	componentModel            componentModel
	catalogModel              componentModel
	assessmentResultsModel    assessmentResultsModel
	planOfActionAndMilestones componentModel
	assessmentPlanModel       componentModel
	systemSecurityPlanModel   componentModel
	profileModel              componentModel
	warnModel                 warnModal
	width                     int
	height                    int
}

func NewOSCALModel(oscalModel oscalTypes_1_1_2.OscalCompleteSchema) model {
	// tabs := checkNonNullFields(oscalModel)
	tabs := []string{
		"ComponentDefinition",
		"AssessmentResults",
		"SystemSecurityPlan",
		"AssessmentPlan",
		"PlanOfActionAndMilestones",
		"Catalog",
		"Profile",
	}

	return model{
		tabs:                   tabs,
		content:                "Welcome to the Lula TUI!",
		oscalModel:             oscalModel,
		componentModel:         NewComponentDefinitionModel(oscalModel.ComponentDefinition),
		assessmentResultsModel: NewAssessmentResultsModel(oscalModel.AssessmentResults),
	}
}

func checkNonNullFields(v interface{}) []string {
	var result []string
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsNil() {
			fieldName := typ.Field(i).Name
			result = append(result, fieldName)
		}
	}

	return result
}

func (m model) closeAllTabs() {
	m.catalogModel.open = false
	m.profileModel.open = false
	m.componentModel.open = false
	m.systemSecurityPlanModel.open = false
	m.assessmentPlanModel.open = false
	m.assessmentResultsModel.open = false
	m.planOfActionAndMilestones.open = false
}

func (m model) loadTabModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.closeAllTabs()
	switch m.tabs[m.activeTab] {
	case "Catalog":
		return nil, nil
	case "Profile":
		return nil, nil
	case "ComponentDefinition":
		m.componentModel.open = true
		return m.componentModel, nil
	case "SystemSecurityPlan":
		return nil, nil
	case "AssessmentPlan":
		return nil, nil
	case "AssessmentResults":
		m.assessmentResultsModel.open = true
		return m.assessmentResultsModel, nil
	case "PlanOfActionAndMilestones":
		return nil, nil
	}
	return nil, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		// m.handleKey(msg.String(), nil)
		switch msg.String() {
		case "tab":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
		case "shift+tab":
			if m.activeTab == 0 {
				m.activeTab = len(m.tabs) - 1
			} else {
				m.activeTab = m.activeTab - 1
			}
		case "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		contentHeight := m.height - 10
		contentWidth := m.width

		// Set the height and width of the models
		m.componentModel.height = contentHeight
		m.componentModel.width = contentWidth

		m.assessmentResultsModel.height = contentHeight
		m.assessmentResultsModel.width = contentWidth
	}

	// if m.tabs[m.activeTab] == "ComponentDefinition" {
	// 	var cmd tea.Cmd
	// 	m.componentModel, cmd = m.componentModel.Update(msg) // Ensure the component model is updated
	// 	cmds = append(cmds, cmd)
	// }

	var cmd tea.Cmd
	tabModel, cmd := m.loadTabModel(msg)
	if tabModel != nil {
		newTabModel, newCmd := tabModel.Update(msg)
		if newTabModel != nil {
			switch m.tabs[m.activeTab] {
			case "ComponentDefinition":
				m.componentModel = newTabModel.(componentModel)
			case "AssessmentResults":
				m.assessmentResultsModel = newTabModel.(assessmentResultsModel)
			}
		}
		cmds = append(cmds, newCmd)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var tabs []string
	for i, t := range m.tabs {
		if i == m.activeTab {
			tabs = append(tabs, activeTab.Render(t))
		} else {
			tabs = append(tabs, tab.Render(t))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	gap := tabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	tabModel, _ := m.loadTabModel(nil)
	if tabModel != nil {
		body := lipgloss.NewStyle().PaddingTop(0).PaddingLeft(2).Render(tabModel.View())
		return fmt.Sprintf("%s\n%s", row, body)
	}

	return row
}
