package tools

import (
	"fmt"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

var uuidHelp = `
To create a new random UUID:
	lula tools uuidgen

To create a deterministic UUID given some source:
	lula tools uuidgen <source>
`

func init() {
	// Kubectl stub command.
	uuidCmd := &cobra.Command{
		Use:   "uuidgen",
		Short: "Generate a UUID",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.SkipLogFile = true
		},
		Long:    "Generate a UUID at random or deterministically with a provided string",
		Example: uuidHelp,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println(uuid.NewUUID())
			} else if len(args) == 1 {
				fmt.Println(uuid.NewUUIDWithSource(args[0]))
			} else {
				message.Fatal(fmt.Errorf("Too many arguments"), "Too many arguments")
			}
		},
	}

	toolsCmd.AddCommand(uuidCmd)
}
