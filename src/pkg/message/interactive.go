package message

import (
	"github.com/pterm/pterm"
)

func PromptForConfirmation() bool {
	// Prompt the user to confirm the action
	confirmation := pterm.DefaultInteractiveConfirm.WithDefaultText("Do you want to continue?")
	result, _ := confirmation.Show()
	return result
}
