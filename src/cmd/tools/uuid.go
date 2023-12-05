package tools

import (
	"fmt"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
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
		Use:     "uuidgen",
		Short:   "Generate a UUID",
		Long:    "Generate a UUID at random or deterministically with a provided input",
		Example: uuidHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Println(uuid.NewUUID())
				return nil
			} else if len(args) == 1 {
				fmt.Println(uuid.NewUUIDWithSource(args[0]))
				return nil
			} else {
				return fmt.Errorf("too many arguments")
			}
		},
	}

	toolsCmd.AddCommand(uuidCmd)
}
