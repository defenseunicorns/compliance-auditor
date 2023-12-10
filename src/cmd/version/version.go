package version

import (
	"fmt"
	"github.com/defenseunicorns/lula/src/config"

	"github.com/spf13/cobra"
)

var versionHelp = `
Get the current Lula version:
	lula version
`

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Shows the current version of the Lula binary",
	Long:    "Shows the current version of the Lula binary",
	Example: versionHelp,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(config.CLIVersion)
		return nil
	},
}

// Include adds the tools command to the root command.
func Include(rootCmd *cobra.Command) {
	rootCmd.AddCommand(versionCmd)
}
