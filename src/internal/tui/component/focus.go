package component

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

var maxFocus = focusValidations
