package component

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

type keys struct {
	Edit          key.Binding
	Generate      key.Binding
	Confirm       key.Binding
	Select        key.Binding
	Detail        key.Binding
	Validate      key.Binding
	Save          key.Binding
	Cancel        key.Binding
	Navigation    key.Binding
	NavigateLeft  key.Binding
	NavigateRight key.Binding
	SwitchModels  key.Binding
	Help          key.Binding
	Quit          key.Binding
}

var componentKeys = keys{
	Quit:          common.CommonKeys.Quit,
	Help:          common.CommonKeys.Help,
	Edit:          common.CommonKeys.Edit,
	Detail:        common.CommonKeys.Detail,
	Save:          common.CommonKeys.Save,
	Validate:      common.CommonKeys.Validate,
	Select:        common.CommonKeys.Select,
	Confirm:       common.CommonKeys.Confirm,
	Cancel:        common.CommonKeys.Cancel,
	Navigation:    common.CommonKeys.Navigation,
	NavigateLeft:  common.CommonKeys.NavigateLeft,
	NavigateRight: common.CommonKeys.NavigateRight,
	SwitchModels:  common.CommonKeys.NavigateModels,
}

var componentEditKeys = keys{
	Confirm: common.PickerKeys.Select,
	Cancel:  common.PickerKeys.Cancel,
}

// Focus key help
var (
	// No focus
	shortHelpNoFocus = []key.Binding{
		componentKeys.Navigation, componentKeys.Validate, componentKeys.SwitchModels, componentKeys.Help,
	}
	fullHelpNoFocusOneLine = []key.Binding{
		componentKeys.Navigation, componentKeys.Validate, componentKeys.SwitchModels, componentKeys.Help,
	}
	fullHelpNoFocus = [][]key.Binding{
		{componentKeys.Navigation}, {componentKeys.Validate}, {componentKeys.SwitchModels}, {componentKeys.Help},
	}

	// focus dialog box
	shortHelpDialogBox = []key.Binding{
		componentKeys.Select, componentKeys.Navigation, componentKeys.Validate, componentKeys.Help,
	}
	fullHelpDialogBoxOneLine = []key.Binding{
		componentKeys.Select, componentKeys.Save, componentKeys.Navigation, componentKeys.Validate, componentKeys.SwitchModels, componentKeys.Help,
	}
	fullHelpDialogBox = [][]key.Binding{
		{componentKeys.Select}, {componentKeys.Save}, {componentKeys.Navigation}, {componentKeys.Validate}, {componentKeys.SwitchModels}, {componentKeys.Help},
	}

	// focus editable dialog box
	shortHelpEditableDialogBox = []key.Binding{
		componentKeys.Edit, componentKeys.Save, componentKeys.Navigation, componentKeys.Help,
	}
	fullHelpEditableDialogBoxOneLine = []key.Binding{
		componentKeys.Edit, componentKeys.Save, componentKeys.Navigation, componentKeys.Validate, componentKeys.SwitchModels, componentKeys.Help,
	}
	fullHelpEditableDialogBox = [][]key.Binding{
		{componentKeys.Edit}, {componentKeys.Save}, {componentKeys.Navigation}, {componentKeys.SwitchModels}, {componentKeys.Help},
	}

	// focus validations
	shortHelpValidations = []key.Binding{
		componentKeys.Detail, common.ListKeys.Up, common.ListKeys.Down, common.ListKeys.Filter, componentKeys.Help,
	}

	fullHelpValidationsOneLine = []key.Binding{
		componentKeys.Detail, common.ListKeys.Up, common.ListKeys.Down, common.ListKeys.Filter, componentKeys.Navigation, componentKeys.Validate, componentKeys.SwitchModels, componentKeys.Help,
	}

	fullHelpValidations = [][]key.Binding{
		{componentKeys.Detail}, {componentKeys.Navigation}, {componentKeys.SwitchModels}, {componentKeys.Help},
	}
)
