package assessmentresults

import (
	"github.com/charmbracelet/bubbles/help"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
)

type Model struct {
	open                bool
	help                help.Model
	keys                keys
	focus               focus
	content             string
	inResultOverlay     bool
	results             []result
	resultsPicker       viewport.Model
	selectedResult      result
	selectedResultIndex int
	compareResult       result
	compareResultIndex  int
	findings            blist.Model
	findingPicker       viewport.Model
	findingSummary      viewport.Model
	observationSummary  viewport.Model
	width               int
	height              int
}

type focus int

const (
	noFocus focus = iota
	focusResultSelection
	focusCompareSelection
	focusFindings
	focusSummary
	focusObservations
)

type result struct {
	uuid, title  string
	findings     *[]oscalTypes_1_1_2.Finding
	observations *[]oscalTypes_1_1_2.Observation
}

type finding struct {
	title, uuid, controlId, state string
	observations                  []observation
}

func (i finding) Title() string       { return i.controlId }
func (i finding) Description() string { return i.state }
func (i finding) FilterValue() string { return i.title }

type observation struct {
	uuid, description, remarks, state, validationId string
}

func (m *Model) Close() {
	m.open = false
}

func (m *Model) Open() {
	m.open = true
}

func (m *Model) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}
