package tui

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	assessmentresults "github.com/defenseunicorns/lula/src/internal/tui/assessment_results"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	"github.com/defenseunicorns/lula/src/internal/tui/component"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

type SwitchTabMsg struct {
	ToTab int
}

type model struct {
	keys                      common.Keys
	tabs                      []string
	activeTab                 int
	componentFilePath         string
	assessmentResultsFilePath string
	writtenComponentModel     *oscalTypes.ComponentDefinition
	componentModel            component.Model
	assessmentResultsModel    assessmentresults.Model
	catalogModel              common.TbdModal
	planOfActionAndMilestones common.TbdModal
	assessmentPlanModel       common.TbdModal
	systemSecurityPlanModel   common.TbdModal
	profileModel              common.TbdModal
	closeModel                common.PopupModel
	saveModel                 common.SaveModel
	width                     int
	height                    int
}

func NewOSCALModel(modelMap map[string]*oscalTypes.OscalCompleteSchema, filePathMap map[string]string, dumpFile *os.File) model {
	tabs := []string{
		"ComponentDefinition",
		"AssessmentResults",
		"SystemSecurityPlan",
		"AssessmentPlan",
		"PlanOfActionAndMilestones",
		"Catalog",
		"Profile",
	}

	if dumpFile != nil {
		common.DumpFile = dumpFile
	}

	// get the right model types assigned to their respective tea models
	writtenComponentModel := new(oscalTypes.ComponentDefinition)
	componentModel := component.NewComponentDefinitionModel(writtenComponentModel)
	componentFilePath := "component.yaml"
	assessmentResultsModel := assessmentresults.NewAssessmentResultsModel(nil)
	assessmentResultsFilePath := "assessment-results.yaml"

	for k, v := range filePathMap {
		switch k {
		case "component":
			componentFilePath = v
		case "assessment-results":
			assessmentResultsFilePath = v
		}
	}

	for k, v := range modelMap {
		// TODO: update these with the UpdateModel functions for the respective models
		switch k {
		case "component":
			componentModel = component.NewComponentDefinitionModel(v.ComponentDefinition)
			err := common.DeepCopy(v.ComponentDefinition, writtenComponentModel)
			if err != nil {
				common.PrintToLog("error creating deep copy of component model: %v", err)
			}
		case "assessment-results":
			assessmentResultsModel = assessmentresults.NewAssessmentResultsModel(v.AssessmentResults)
		}
	}

	saveModel := common.NewSaveModel(componentFilePath)
	closeModel := common.NewPopupModel("Quit Console", "Are you sure you want to quit the Lula Console?", []key.Binding{common.CommonKeys.Confirm, common.CommonKeys.Cancel})

	return model{
		keys:                      common.CommonKeys,
		tabs:                      tabs,
		componentFilePath:         componentFilePath,
		assessmentResultsFilePath: assessmentResultsFilePath,
		writtenComponentModel:     writtenComponentModel,
		closeModel:                closeModel,
		saveModel:                 saveModel,
		componentModel:            componentModel,
		assessmentResultsModel:    assessmentResultsModel,
		systemSecurityPlanModel:   common.NewTbdModal("System Security Plan"),
		catalogModel:              common.NewTbdModal("Catalog"),
		profileModel:              common.NewTbdModal("Profile"),
		assessmentPlanModel:       common.NewTbdModal("Assessment Plan"),
		planOfActionAndMilestones: common.NewTbdModal("Plan of Action and Milestones"),
		width:                     common.DefaultWidth,
		height:                    common.DefaultHeight,
	}
}

func (m *model) isModelSaved() bool {
	return reflect.DeepEqual(m.writtenComponentModel, m.componentModel.GetComponentDefinition())
}

// WriteOscalModel runs on save cmds
func (m *model) writeOscalModel() tea.Msg {
	common.PrintToLog("componentFilePath: %s", m.componentFilePath)

	saveStart := time.Now()
	err := oscal.OverwriteOscalModel(m.componentFilePath, &oscalTypes.OscalCompleteSchema{ComponentDefinition: m.componentModel.GetComponentDefinition()})
	saveDuration := time.Since(saveStart)
	// just adding a minimum of 2 seconds to the "saving" popup
	if saveDuration < time.Second*2 {
		time.Sleep(time.Second*2 - saveDuration)
	}
	if err != nil {
		common.PrintToLog("error writing component model: %v", err)
		return common.SaveFailMsg{Err: err}
	}
	common.PrintToLog("model saved")

	_ = common.DeepCopy(m.componentModel.GetComponentDefinition(), m.writtenComponentModel)
	return common.SaveSuccessMsg{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	common.DumpToLog(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		k := msg.String()

		switch k {
		case common.ContainsKey(k, m.keys.ModelRight.Keys()):
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			cmds = append(cmds, func() tea.Msg {
				return SwitchTabMsg{
					ToTab: m.activeTab,
				}
			})

		case common.ContainsKey(k, m.keys.ModelLeft.Keys()):
			if m.activeTab == 0 {
				m.activeTab = len(m.tabs) - 1
			} else {
				m.activeTab = m.activeTab - 1
			}
			cmds = append(cmds, func() tea.Msg {
				return SwitchTabMsg{
					ToTab: m.activeTab,
				}
			})

		case common.ContainsKey(k, m.keys.Confirm.Keys()):
			if m.closeModel.Open {
				return m, tea.Quit
			}

		case common.ContainsKey(k, m.keys.Save.Keys()):
			m.saveModel.RenderedDuringQuit = false
			if m.closeModel.Open {
				m.saveModel.RenderedDuringQuit = true
				if m.isModelSaved() {
					return m, nil
				} else {
					m.closeModel.Open = false
				}
			}

			m.saveModel.Open = true
			m.saveModel.Save = true
			if m.isModelSaved() {
				m.saveModel.Save = false
				m.saveModel.Content = "No changes to save"
				return m, nil
			}
			m.saveModel.Content = fmt.Sprintf("Save changes to %s?", m.saveModel.FilePath)
			// warning if file exists
			if _, err := os.Stat(m.componentFilePath); err == nil {
				m.saveModel.Warning = fmt.Sprintf("%s will be overwritten", m.componentFilePath)
			}

			return m, nil

		case common.ContainsKey(k, m.keys.Cancel.Keys()):
			if m.closeModel.Open {
				m.closeModel.Open = false
			} else if m.saveModel.Open {
				m.saveModel.Open = false
			}

		case common.ContainsKey(k, m.keys.Quit.Keys()):
			// add quit warn pop-up
			if !m.isModelSaved() {
				m.closeModel.Warning = "Changes not written"
				m.closeModel.Help.ShortHelp = []key.Binding{common.CommonKeys.Confirm, common.CommonKeys.Save, common.CommonKeys.Cancel}
			}
			if m.closeModel.Open {
				return m, tea.Quit
			} else {
				m.closeModel.Open = true
			}
		}

	case common.SaveModelMsg:
		saveResultMsg := m.writeOscalModel()
		common.DumpToLog(saveResultMsg)

		cmds = append(cmds, func() tea.Msg {
			return saveResultMsg
		}, func() tea.Msg {
			time.Sleep(time.Second * 2)
			return common.SaveCloseAndResetMsg{}
		})

		if (saveResultMsg == common.SaveSuccessMsg{}) && msg.InQuitWorkflow {
			cmds = append(cmds, tea.Quit)
		}
		return m, tea.Sequence(cmds...)

	case component.ValidationDataMsg:
		// Pass through to update assessment results model
		mdl, cmd := m.assessmentResultsModel.Update(msg)
		m.assessmentResultsModel = mdl.(assessmentresults.Model)
		cmds = append(cmds, cmd)

	case assessmentresults.AssessmentUpdatedMsg:
		// Save assessment results data
		// TODO: add to save workflow
		assessmentResults := m.assessmentResultsModel.GetAssessmentResults()
		if assessmentResults != nil {
			err := oscal.OverwriteOscalModel(m.assessmentResultsFilePath, &oscalTypes.OscalCompleteSchema{AssessmentResults: assessmentResults})
			if err != nil {
				common.PrintToLog("error writing assessment results model: %v", err)
			}

			m.activeTab = 1 // assessment results tab
			cmds = append(cmds, func() tea.Msg {
				return SwitchTabMsg{
					ToTab: m.activeTab,
				}
			})
		}

	case SwitchTabMsg:
		return m, m.openTab(msg.ToTab)
	}

	mdl, cmd := m.saveModel.Update(msg)
	m.saveModel = mdl.(common.SaveModel)
	cmds = append(cmds, cmd)

	// Only run update methods on active tab
	switch m.tabs[m.activeTab] {
	case "ComponentDefinition":
		mdl, cmd = m.componentModel.Update(msg)
		m.componentModel = mdl.(component.Model)
		cmds = append(cmds, cmd)
	case "AssessmentResults":
		mdl, cmd = m.assessmentResultsModel.Update(msg)
		m.assessmentResultsModel = mdl.(assessmentresults.Model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.closeModel.Open {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.closeModel.View(), lipgloss.WithWhitespaceChars(" "))
	}
	if m.saveModel.Open {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.saveModel.View(), lipgloss.WithWhitespaceChars(" "))
	}
	return m.mainView()
}

func (m model) mainView() string {
	var tabs []string
	for i, t := range m.tabs {
		if i == m.activeTab {
			tabs = append(tabs, common.ActiveTab.Render(t))
		} else {
			tabs = append(tabs, common.Tab.Render(t))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	gap := common.TabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	content := ""
	switch m.tabs[m.activeTab] {
	case "ComponentDefinition":
		content = m.componentModel.View()
	case "AssessmentResults":
		content = m.assessmentResultsModel.View()
	case "SystemSecurityPlan":
		content = m.systemSecurityPlanModel.View()
	case "AssessmentPlan":
		content = m.assessmentPlanModel.View()
	case "PlanOfActionAndMilestones":
		content = m.planOfActionAndMilestones.View()
	case "Catalog":
		content = m.catalogModel.View()
	case "Profile":
		content = m.profileModel.View()
	}

	body := lipgloss.NewStyle().PaddingTop(0).PaddingLeft(2).Render(content)
	return fmt.Sprintf("%s\n%s", row, body)
}

func (m *model) openTab(tab int) func() tea.Msg {
	switch m.tabs[tab] {
	case "ComponentDefinition":
		return func() tea.Msg {
			return component.ModelOpenMsg{
				Height: m.height - common.TabOffset,
				Width:  m.width,
			}
		}
	case "AssessmentResults":
		return func() tea.Msg {
			return assessmentresults.ModelOpenMsg{
				Height: m.height - common.TabOffset,
				Width:  m.width,
			}
		}
	}
	return func() tea.Msg {
		return nil
	}
}
