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
	Use: "version",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.SkipLogFile = true
	},
	Short:   "Shows the current version of the Lula binary",
	Long:    "Shows the current version of the Lula binary",
	Example: versionHelp,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.CLIVersion)
	},
}

// Include adds the tools command to the root command.
func Include(rootCmd *cobra.Command) {
	rootCmd.AddCommand(versionCmd)
}
