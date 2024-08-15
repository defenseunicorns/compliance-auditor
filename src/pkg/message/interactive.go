package message

import (
	"fmt"

	"github.com/pterm/pterm"
)

func PromptForConfirmation(spinner *Spinner) bool {
	// Prompt the user to confirm the action
	if spinner != nil {
		spinnerText := spinner.Pause()
		defer spinner.Updatef(fmt.Sprintf("%s\n", spinnerText))
	}

	confirmation := pterm.DefaultInteractiveConfirm.WithDefaultText("Do you want to run executable validations?")
	result, _ := confirmation.Show()

	return result
}
