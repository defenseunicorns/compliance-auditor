package component

import (
	"github.com/charmbracelet/bubbles/help"
	blist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/defenseunicorns/lula/src/types"
)

type Model struct {
	open                   bool
	help                   help.Model
	keys                   keys
	focus                  focus
	content                string
	inComponentOverlay     bool
	components             []component
	selectedComponent      component
	selectedComponentIndex int
	componentPicker        viewport.Model
	inFrameworkOverlay     bool
	frameworks             []framework
	selectedFramework      framework
	selectedFrameworkIndex int
	frameworkPicker        viewport.Model
	controlPicker          viewport.Model
	controls               blist.Model
	selectedControl        control
	remarks                viewport.Model
	description            viewport.Model
	validationPicker       viewport.Model
	validations            blist.Model
	selectedValidation     validationLink
	width                  int
	height                 int
}

type focus int

const (
	noComponentFocus focus = iota
	focusComponentSelection
	focusFrameworkSelection
	focusControls
	focusRemarks
	focusDescription
	focusValidations
)

type component struct {
	uuid, title, desc string
	frameworks        []framework
}

type framework struct {
	name     string
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
