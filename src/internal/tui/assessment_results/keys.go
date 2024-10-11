package assessmentresults

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
)

type keys struct {
	Validate      key.Binding
	Evaluate      key.Binding
	Confirm       key.Binding
	Cancel        key.Binding
	Navigation    key.Binding
	NavigateLeft  key.Binding
	NavigateRight key.Binding
	SwitchModels  key.Binding
	Detail        key.Binding
	Filter        key.Binding
	Help          key.Binding
	Quit          key.Binding
}

var assessmentKeys = keys{
	Quit:          common.CommonKeys.Quit,
	Help:          common.CommonKeys.Help,
	Confirm:       common.CommonKeys.Confirm,
	Cancel:        common.CommonKeys.Cancel,
	Navigation:    common.CommonKeys.Navigation,
	NavigateLeft:  common.CommonKeys.NavigateLeft,
	NavigateRight: common.CommonKeys.NavigateRight,
	SwitchModels:  common.CommonKeys.NavigateModels,
	Detail:        common.CommonKeys.Detail,
	Filter:        common.TableKeys.Filter,
}

var assessmentKeysInFilter = keys{
	Confirm: common.CommonKeys.Confirm,
	Cancel:  common.CommonKeys.Cancel,
}

var (
	// No focus
	shortHelpNoFocus = []key.Binding{
		assessmentKeys.Navigation, assessmentKeys.SwitchModels, assessmentKeys.Help,
	}
	fullHelpNoFocusOneLine = []key.Binding{
		assessmentKeys.Navigation, assessmentKeys.SwitchModels, assessmentKeys.Help,
	}
	fullHelpNoFocus = [][]key.Binding{
		{assessmentKeys.Navigation}, {assessmentKeys.SwitchModels}, {assessmentKeys.Help},
	}
)
